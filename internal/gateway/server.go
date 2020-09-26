package gateway

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/firstcontributions/firstcontributions/internal/gateway/configs"
	"github.com/firstcontributions/firstcontributions/internal/gateway/csrf"
	"github.com/firstcontributions/firstcontributions/internal/gateway/models/redis"
	"github.com/firstcontributions/firstcontributions/internal/gateway/session"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

// Server encapsulates the connection objects and configs for
// the gateway server
type Server struct {
	*configs.Config
	SessionManager *session.Manager
	CookieManager  *securecookie.SecureCookie
	Router         *mux.Router
	CSRFManager    *csrf.Manager
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
			time.Duration(*s.SessionTTLDays)*time.Hour*24,
		),
	)
	s.CookieManager = securecookie.New([]byte(*s.HashKey), []byte(*s.BlockKey))

	s.InitRoutes()
	return nil
}

func (s *Server) InitRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/v1/auth/redirect", s.AuthRedirect)

	r.HandleFunc("/v1/auth/callback", s.AuthCallback)

	s.Router = r
}

func (s *Server) ListenAndServe() error {
	log.Printf("listening at :%s", *s.Port)
	return http.ListenAndServe(":"+*s.Port, s.Router)
}
