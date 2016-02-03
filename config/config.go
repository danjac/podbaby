package config

import (
	"errors"
	"github.com/danjac/podbaby/config/Godeps/_workspace/src/github.com/gorilla/securecookie"
	"log"
)

const (
	defaultPort            = 5000
	defaultDBConnections   = 99
	defaultStaticURL       = "/static/"
	defaultStaticDir       = "./static/"
	defaultTemplateDir     = "./templates"
	defaultMailTemplateDir = "./templates/email"
	devStaticURL           = "http://localhost:8080/static/"
)

var (
	ErrMissingDatabaseURL = errors.New("Database URL is missing")
	ErrMissingSecretKey   = errors.New("Secret key is missing")
)

func New() *Config {
	return &Config{
		Mail: &MailConfig{
			Host:        "localhost",
			ID:          "localhost",
			TemplateDir: defaultMailTemplateDir,
		},
		Env:               "prod",
		Port:              defaultPort,
		MaxDBConnections:  defaultDBConnections,
		StaticDir:         defaultStaticDir,
		StaticURL:         defaultStaticURL,
		TemplateDir:       defaultTemplateDir,
		DynamicContentURL: devStaticURL,
		SecretKey:         generateRandomKey(),
	}
}

func (cfg *Config) Validate() error {
	if cfg.DatabaseURL == "" {
		return ErrMissingDatabaseURL
	}
	if cfg.SecretKey == "" {
		return ErrMissingSecretKey
	}
	return nil
}

func (cfg *Config) MustValidate() {
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
}

func (cfg *Config) IsDev() bool {
	return cfg.Env == "dev"
}

func (cfg *Config) IsProd() bool {
	return cfg.Env == "prod"
}

// MailConfig contains SMTP settings

type MailConfig struct {
	TemplateDir,
	Addr,
	ID,
	User,
	Password,
	Host string
}

// Config is server configuration
type Config struct {
	Mail             *MailConfig
	Port             int
	MaxDBConnections int
	Env,
	DatabaseURL,
	StaticURL,
	DynamicContentURL,
	TemplateDir,
	StaticDir,
	GoogleAnalyticsID,
	SecureCookieKey,
	SecretKey string
}

func generateRandomKey() string {
	return string(securecookie.GenerateRandomKey(32))
}
