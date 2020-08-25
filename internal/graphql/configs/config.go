package config

import (
	commonconfigs "github.com/firstcontributions/firstcontributions/internal/configs"
)

// Config keeps the service level config data for luffy
//go:generate envparser generate -t Config -f $GOFILE
type Config struct {
	Log  *commonconfigs.LogConfig
	Port *string `env:"PROFILE_PORT"`
}
