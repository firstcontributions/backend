package configs

import (
	commonconfigs "github.com/firstcontributions/firstcontributions/internal/configs"
)

// Config keeps the service level config data for luffy
//go:generate envparser generate -t Config -f $GOFILE
type Config struct {
	Log            *commonconfigs.LogConfig
	RedisConfig    *RedisConfig
	Port           *string `env:"GATEWAY_PORT"`
	SessionTTLDays *int    `env:"SESSION_TTL_DAYS"`
	HashKey        *string `env:"HASH_KEY"`
	BlockKey       *string `env:"BLOCK_KEY"`
	GithubConfig   *commonconfigs.GithubConfig
	ProfileManager *commonconfigs.ProfileManager
}

// RedisConfig encapsulates the redis configs
//go:generate envparser generate -t RedisConfig -f $GOFILE
type RedisConfig struct {
	Port     *string `env:"REDIS_PORT"`
	Host     *string `env:"REDIS_HOST"`
	Password *string `env:"REDIS_PASSWORD"`
}
