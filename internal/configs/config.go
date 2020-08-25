package configs

// LogConfig encapsulates the log related configs
//go:generate envparser generate -t LogConfig -f $GOFILE
type LogConfig struct {
	Level *string `env:"LOG_LEVEL"`
	Path  *string `env:"LOG_PATH"`
}
