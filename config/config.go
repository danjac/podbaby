package config

import (
	"encoding/base64"
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
	// ErrMissingDatabaseURL is returned if no database connection URL is set in config
	ErrMissingDatabaseURL = errors.New("Database URL is missing")
	// ErrMissingSecretKey is returned if no secret key  is set in config
	ErrMissingSecretKey = errors.New("Secret key is missing")
)

// Default returns a ready-made configuration with sensible settings
func Default() *Config {
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
		SecretKey:         RandomKey(),
	}
}

// Validate checks that all required configuration settings are present and correct
func (cfg *Config) Validate() error {
	if cfg.DatabaseURL == "" {
		return ErrMissingDatabaseURL
	}
	if cfg.SecretKey == "" {
		return ErrMissingSecretKey
	}
	return nil
}

// MustValidate checks that all required configuration settings are present and correct
func (cfg *Config) MustValidate() {
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
}

// IsDev checks if configuration is set to development mode
func (cfg *Config) IsDev() bool {
	return cfg.Env == "dev"
}

// IsProd checks if configuration is set to production mode
func (cfg *Config) IsProd() bool {
	return !cfg.IsDev()
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

// RandomKey generates a printable, random 32 byte string
func RandomKey() string {
	return base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
}
