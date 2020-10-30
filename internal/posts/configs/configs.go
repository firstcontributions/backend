package configs

import (
	commonconfigs "github.com/firstcontributions/backend/internal/configs"
)

// Config keeps the service level config data for luffy
//go:generate envparser generate -t Config -f $GOFILE
type Config struct {
	Log      *commonconfigs.LogConfig
	Port     *string `env:"POSTS_PORT"`
	MongoURL *string `env:"MONGO_URL"`
}
