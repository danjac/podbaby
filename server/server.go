package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/database"
	"github.com/gorilla/securecookie"
	"github.com/unrolled/render"
)

// Config is server configuration
type Config struct {
	StaticURL string
	StaticDir string
	SecretKey string
}

type Server struct {
	DB     *database.DB
	Config *Config
	Log    *logrus.Logger
	Render *render.Render
	Cookie *securecookie.SecureCookie
}

func New(db *database.DB, log *logrus.Logger, cfg *Config) *Server {

	cookie := securecookie.New(
		[]byte(cfg.SecretKey),
		securecookie.GenerateRandomKey(32))

	return &Server{
		DB:     db,
		Config: cfg,
		Log:    log,
		Render: render.New(),
		Cookie: cookie,
	}
}
