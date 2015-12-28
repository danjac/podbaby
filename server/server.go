package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/database"
	"github.com/gorilla/mux"
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

func getInt64(r *http.Request, name string) (int64, error) {
	badRequest := HTTPError{http.StatusBadRequest, errors.New("Invalid parameter for " + name)}
	value, ok := mux.Vars(r)[name]
	if !ok {
		return 0, badRequest
	}
	intval, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, badRequest
	}
	return intval, nil
}

func getPage(r *http.Request) int64 {
	value := r.FormValue("page")
	if value == "" {
		return 1
	}
	page, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		page = 1
	}
	return page
}
