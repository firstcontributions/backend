package gateway

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/firstcontributions/backend/internal/gateway/configs"
	"github.com/firstcontributions/backend/internal/gateway/csrf"
	graphqlschema "github.com/firstcontributions/backend/internal/gateway/graphql"
	"github.com/firstcontributions/backend/internal/gateway/models/redis"
	"github.com/firstcontributions/backend/internal/gateway/rpcs"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	pool "github.com/processout/grpc-go-pool"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

// Server encapsulates the connection objects and configs for
// the gateway server
type Server struct {
	*configs.Config
	SessionManager *session.Manager
	CookieManager  *securecookie.SecureCookie
	Router         *mux.Router
	CSRFManager    *csrf.Manager
	ProfileManager *rpcs.ProfileManager
}

// NewServer returns an instance of server
func NewServer() *Server {
	return &Server{
		Config: &configs.Config{},
	}
}

// Init will initialise configs, connections etc.
func (s *Server) Init() error {
	if err := s.DecodeEnv(); err != nil {
		return fmt.Errorf("could not initialize config, Error: %w", err)
	}
	s.SessionManager = session.NewManager(
		redis.NewSessionStore(
			*s.RedisConfig.Host,
			*s.RedisConfig.Port,
			*s.RedisConfig.Password,
			time.Duration(*s.SessionTTLDays)*time.Hour*24,
		),
	)
	s.CSRFManager = csrf.NewManager(
		redis.NewCSRFStore(
			*s.RedisConfig.Host,
			*s.RedisConfig.Port,
			*s.RedisConfig.Password,
			time.Duration(*s.CSRFTTLSeconds)*time.Second,
		),
	)
	s.CookieManager = securecookie.New([]byte(*s.HashKey), []byte(*s.BlockKey))

	profileConn, err := pool.New(
		func() (*grpc.ClientConn, error) {
			return grpc.Dial(*s.Profile.URL, grpc.WithInsecure())
		},
		*s.Profile.InitConnections,
		*s.Profile.ConnectionCapacity,
		time.Duration(*s.Profile.ConnectionTTLMinutes)*time.Minute,
	)
	if err != nil {
		return fmt.Errorf("could not initialize connection to profile manager, Error: %w", err)
	}
	s.ProfileManager = rpcs.NewProfileManager(profileConn)
	if err := s.InitRoutes(); err != nil {
		return err
	}
	return nil
}

func (s *Server) InitRoutes() error {
	r := mux.NewRouter()
	r.HandleFunc("/v1/auth/redirect", s.AuthRedirect)
	r.HandleFunc("/v1/auth/callback", s.AuthCallback)

	schema, err := s.GetGraphqlSchema()
	if err != nil {
		return err
	}
	r.Handle("/v1/graphql", s.HandleSession(&relay.Handler{Schema: schema}))
	r.Handle("/v1/session", s.HandleSession(s.SessionHandler()))

	s.Router = r
	return nil
}

func (s *Server) ListenAndServe() error {
	log.Printf("listening at :%s", *s.Port)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://explorer.firstcontributions.com"},
		AllowCredentials: true,
	})
	return http.ListenAndServe(":"+*s.Port, c.Handler(s.Router))
}

func (s *Server) GetGraphqlSchema() (*graphql.Schema, error) {
	schema, err := ioutil.ReadFile("./schema.graphql")
	if err != nil {
		return nil, err
	}
	resolver := &graphqlschema.Resolver{
		ProfileManager: s.ProfileManager,
	}
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	return graphql.MustParseSchema(string(schema), resolver, opts...), nil

}
