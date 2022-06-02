package configs

import (
	commonconfigs "github.com/firstcontributions/backend/internal/configs"
)

// Config keeps the service level config data for luffy
//go:generate envparser generate -t Config -f $GOFILE
type Config struct {
	Log                *commonconfigs.LogConfig
	RedisConfig        *RedisConfig
	Port               *string `env:"GATEWAY_PORT"`
	SessionTTLDays     *int    `env:"SESSION_TTL_DAYS"`
	CSRFTTLSeconds     *int    `env:"CSRF_TTL_SECONDS"`
	HashKey            *string `env:"HASH_KEY"`
	BlockKey           *string `env:"BLOCK_KEY"`
	GithubConfig       *commonconfigs.GithubConfig
	MongoURL           *string `env:"MONGO_URL"`
	UsersServiceConfig *UsersServiceConfig
}

// RedisConfig encapsulates the redis configs
//go:generate envparser generate -t RedisConfig -f $GOFILE
type RedisConfig struct {
	Port     *string `env:"REDIS_PORT"`
	Host     *string `env:"REDIS_HOST"`
	Password *string `env:"REDIS_PASSWORD"`
}

// UsersServiceConfig encapsulates the user service grpc
//go:generate envparser generate -t UsersServiceConfig -f $GOFILE
type UsersServiceConfig struct {
	URL                *string `env:"USER_SERVICE_URL"`
	InitConnections    *int    `env:"USER_SERVICE_INIT_CONN"`
	ConnectionCapacity *int    `env:"USER_SERVICE_CONN_CAPACITY"`
	TTL                *int    `env:"USER_SERVICE_CONN_TTL"`
}
