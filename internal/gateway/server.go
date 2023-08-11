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
	"github.com/firstcontributions/backend/internal/gateway/session"
	graphqlschema "github.com/firstcontributions/backend/internal/graphql/schema"
	"github.com/firstcontributions/backend/internal/models/issuesstore/githubstore"
	storymongo "github.com/firstcontributions/backend/internal/models/storiesstore/mongo"
	usermongo "github.com/firstcontributions/backend/internal/models/usersstore/mongo"

	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	otelgraphql "github.com/graph-gophers/graphql-go/trace/otel"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
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
	userStore, err := usermongo.NewUsersStore(ctx, *s.Config.MongoURL)
	if err != nil {
		return err
	}
	storyStore, err := storymongo.NewStoriesStore(ctx, *s.Config.MongoURL)
	if err != nil {
		return err
	}
	s.Store = storemanager.NewStore(githubstore.NewGitHubStore(*s.GithubConfig), storyStore, userStore)
	if err := s.InitRoutes(); err != nil {
		return err
	}
	return nil
}

func (s *Server) InitRoutes() error {
	r := mux.NewRouter()
	r.Use(otelmux.Middleware("gateway"))
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
		AllowedOrigins:   []string{"http://explorer.opensource.forum", "http://app.opensource.forum"},
		AllowCredentials: true,
	})
	tp, err := tracerProvider()
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)
	return http.ListenAndServe(":"+*s.Port, c.Handler(s.Router))
}

func (s *Server) GetGraphqlSchema() (*graphql.Schema, error) {
	schema, err := ioutil.ReadFile("./assets/schema.graphql")
	if err != nil {
		return nil, err
	}

	resolver := &graphqlschema.Resolver{}
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.Tracer(otelgraphql.DefaultTracer())}
	return graphql.MustParseSchema(string(schema), resolver, opts...), nil

}

func tracerProvider() (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithAgentEndpoint())
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("gateway"),
			attribute.String("environment", "prod"),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}
