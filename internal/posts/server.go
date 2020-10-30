package posts

import (
	"context"

	"github.com/firstcontributions/backend/internal/posts/models"
	"github.com/firstcontributions/backend/internal/profile/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Service keeps configs, connection pool objects etc.
type Service struct {
	*configs.Config
	MongoClient *mongo.Client
	PostManager *PostManager
}

// PostManager is wrapper for model store interface
type PostManager struct {
	models.Store
}

// NewPostManager returns an instance of post manager
func NewPostManager(store models.Store) *PostManager {
	return &PostManager{
		Store: store,
	}
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
	ctx := context.Background()
	if err := s.DecodeEnv(); err != nil {
		return err
	}
	if err := s.InitMongoClient(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) InitMongoClient(ctx context.Context) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(*s.MongoURL))
	if err != nil {
		return err
	}
	if err := client.Connect(ctx); err != nil {
		return err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	s.MongoClient = client
	return nil
}
