package configs

// LogConfig encapsulates the log related configs
//go:generate envparser generate -t LogConfig -f $GOFILE
type LogConfig struct {
	Level *string `env:"LOG_LEVEL"`
	Path  *string `env:"LOG_PATH"`
}

// GithubConfig stores github oauth creds
//go:generate envparser generate -t GithubConfig -f $GOFILE
type GithubConfig struct {
	ClientID     *string  `env:"GITHUB_CLIENT_ID"`
	ClientSecret *string  `env:"GITHUB_CLIENT_SECRET"`
	AuthScopes   []string `env:"GITHUB_AUTH_SCOPES"`
	AuthRedirect *string  `env:"GITHUB_AUTH_REDIRECT"`
}

// ProfileManager stores configs for the grpc conn
//go:generate envparser generate -t ProfileManager -f $GOFILE
type ProfileManager struct {
	URL                  *string `env:"PROFILE_MANAGER_URL"`
	InitConnections      *int    `env:"PROFILE_MANAGER_INIT_CONN"`
	ConnectionCapacity   *int    `env:"PROFILE_MANAGER_CONN_CAPACITY"`
	ConnectionTTLMinutes *int    `env:"PROFILE_MANAGER_CONN_TTL"`
}
