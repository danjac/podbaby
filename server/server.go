package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/database"
	"github.com/unrolled/render"
)

type Config struct {
	StaticURL string
	StaticDir string
}

type Server struct {
	DB     *database.DB
	Config *Config
	Log    *logrus.Logger
	Render *render.Render
}

func New(db *database.DB, log *logrus.Logger, cfg *Config) *Server {
	return &Server{
		DB:     db,
		Config: cfg,
		Log:    log,
		Render: render.New(),
	}
}
