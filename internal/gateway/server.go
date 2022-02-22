package gateway

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/firstcontributions/backend/internal/gateway/configs"
	"github.com/firstcontributions/backend/internal/gateway/csrf"
	"github.com/firstcontributions/backend/internal/gateway/models/redis"
	graphqlschema "github.com/firstcontributions/backend/internal/graphql/schema"
	"github.com/firstcontributions/backend/internal/models/issuesstore/githubstore"
	"github.com/firstcontributions/backend/internal/models/usersstore/mongo"
	"github.com/firstcontributions/backend/internal/storemanager"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/rs/cors"
)

// Server encapsulates the connection objects and configs for
// the gateway server
type Server struct {
	*configs.Config
	SessionManager *session.Manager
	CookieManager  *securecookie.SecureCookie
	Router         *mux.Router
	CSRFManager    *csrf.Manager
	Store          *storemanager.Store
}

// NewServer returns an instance of server
func NewServer() *Server {
	return &Server{
		Config: &configs.Config{},
	}
}

// Init will initialise configs, connections etc.
func (s *Server) Init() error {
	var err error
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

	ctx := context.Background()
	userStore, err := mongo.NewUsersStore(ctx, *s.MongoURL)
	if err != nil {
		return err
	}
	s.Store = storemanager.NewStore(githubstore.NewGitHubStore(*s.GithubConfig), userStore)
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

	resolver := &graphqlschema.Resolver{}
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	return graphql.MustParseSchema(string(schema), resolver, opts...), nil

}
