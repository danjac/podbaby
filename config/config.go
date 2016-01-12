package config

import (
	"errors"
	"github.com/gorilla/securecookie"
	"log"
)

const (
	defaultPort      = 5000
	defaultStaticURL = "/static/"
	defaultStaticDir = "./static/"
	devStaticURL     = "http://localhost:8080/static/"
)

var (
	ErrMissingDatabaseURL = errors.New("Database URL is missing")
	ErrMissingSecretKey   = errors.New("Secret key is missing")
)

func New() *Config {
	return &Config{
		Mail:              &MailConfig{},
		Env:               "prod",
		Port:              defaultPort,
		StaticDir:         defaultStaticDir,
		StaticURL:         defaultStaticURL,
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
	Addr,
	ID,
	User,
	Password,
	Host string
}

// Config is server configuration
type Config struct {
	Mail *MailConfig
	Port int
	Env,
	DatabaseURL,
	StaticURL,
	DynamicContentURL,
	StaticDir,
	GoogleAnalyticsID,
	SecureCookieKey,
	SecretKey string
}

func generateRandomKey() string {
	return string(securecookie.GenerateRandomKey(32))
}
