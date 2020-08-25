package profile

import (
	"github.com/firstcontributions/firstcontributions/internal/profile/configs"
)

// Service keeps configs, connection pool objects etc.
type Service struct {
	*configs.Config
}

// NewService returns an instance of profile service
func NewService() *Service {
	return &Service{
		Config: &configs.Config{},
	}
}

// Init will initialize the configs, all connections to dbs
// and other infra layers
func (s *Service) Init() error {
	if err := s.DecodeEnv(); err != nil {
		return err
	}
	return nil
}
