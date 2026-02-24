package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port        int             `envconfig:"PORT" default:"8080"`
	Environment string          `envconfig:"ENVIRONMENT" default:"development"`
	SiteURL     string          `envconfig:"SITE_URL" default:"http://localhost:3000"`
	Database    DatabaseConfig  `envconfig:"DATABASE"`
	JWT         JWTConfig       `envconfig:"JWT"`
	SMTP        SMTPConfig      `envconfig:"SMTP"`
	OAuth       OAuthConfig     `envconfig:"OAUTH"`
	Redis       RedisConfig     `envconfig:"REDIS"`
	Storage     StorageConfig   `envconfig:"STORAGE"`
	RateLimit   RateLimitConfig `envconfig:"RATE_LIMIT"`
	Seed        SeedConfig      `envconfig:"SEED"`
	LLM         LLMConfig       `envconfig:"LLM"`
}

type LLMConfig struct {
	Provider  string `envconfig:"PROVIDER" default:"openai"`
	APIKey    string `envconfig:"API_KEY"`
	Model     string `envconfig:"MODEL" default:"gpt-4o-mini"`
	MaxTokens int    `envconfig:"MAX_TOKENS" default:"1024"`
	Enabled   bool   `envconfig:"ENABLED" default:"false"`
}

type SeedConfig struct {
	AdminEmail     string `envconfig:"ADMIN_EMAIL" default:"admin@open-pm.dev"`
	AdminPassword  string `envconfig:"ADMIN_PASSWORD" default:"password123"`
	AdminFirstName string `envconfig:"ADMIN_FIRST_NAME" default:"Admin"`
	AdminLastName  string `envconfig:"ADMIN_LAST_NAME" default:"User"`
	WorkspaceName  string `envconfig:"WORKSPACE_NAME" default:"My Workspace"`
	WorkspaceSlug  string `envconfig:"WORKSPACE_SLUG" default:"my-workspace"`
}

type DatabaseConfig struct {
	URL             string `envconfig:"URL" required:"true"`
	MaxOpenConns    int    `envconfig:"MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int    `envconfig:"MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime int    `envconfig:"CONN_MAX_LIFETIME" default:"300"`
}

type JWTConfig struct {
	Secret     string `envconfig:"SECRET" required:"true"`
	Exp        int    `envconfig:"EXP" default:"3600"`
	RefreshExp int    `envconfig:"REFRESH_EXP" default:"604800"`
	Issuer     string `envconfig:"ISSUER" default:"open-pm"`
}

type SMTPConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"1025"`
	User     string `envconfig:"USER"`
	Pass     string `envconfig:"PASS"`
	FromName string `envconfig:"FROM_NAME" default:"Open-PM"`
	FromAddr string `envconfig:"FROM_ADDR" default:"noreply@open-pm.dev"`
}

type OAuthConfig struct {
	GitHub OAuthProviderConfig `envconfig:"GITHUB"`
	Google OAuthProviderConfig `envconfig:"GOOGLE"`
}

type OAuthProviderConfig struct {
	Enabled      bool   `envconfig:"ENABLED" default:"false"`
	ClientID     string `envconfig:"CLIENT_ID"`
	ClientSecret string `envconfig:"CLIENT_SECRET"`
	RedirectURL  string `envconfig:"REDIRECT_URL"`
}

type RedisConfig struct {
	URL string `envconfig:"URL" default:"redis://localhost:6379/0"`
}

type StorageConfig struct {
	Endpoint  string `envconfig:"ENDPOINT" default:"localhost:9000"`
	AccessKey string `envconfig:"ACCESS_KEY" default:"minioadmin"`
	SecretKey string `envconfig:"SECRET_KEY" default:"minioadmin"`
	Bucket    string `envconfig:"BUCKET" default:"open-pm-assets"`
	UseSSL    bool   `envconfig:"USE_SSL" default:"false"`
}

type RateLimitConfig struct {
	Enabled bool    `envconfig:"ENABLED" default:"true"`
	RPS     float64 `envconfig:"RPS" default:"10"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("OPENPM", &cfg)
	return &cfg, err
}
