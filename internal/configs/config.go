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
