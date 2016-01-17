package server

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"gopkg.in/unrolled/render.v1"
)

const (
	cookieUserID  = "userid"
	userKey       = "user"
	cookieTimeout = 24
)

type authLevel int

const (
	authLevelIgnore authLevel = iota
	authLevelOptional
	authLevelRequired
)

var (
	errNotAuthenticated = errors.New("User not logged in")
	errBadRequest       = errors.New("Bad request")
)

type Server struct {
	DB         *database.DB
	Config     *config.Config
	Log        *logrus.Logger
	Render     *render.Render
	Cookie     *securecookie.SecureCookie
	Feedparser feedparser.Feedparser
	Mailer     mailer.Mailer
}

func New(db *database.DB,
	mailer mailer.Mailer,
	log *logrus.Logger,
	cfg *config.Config) *Server {

	secureCookieKey, _ := base64.StdEncoding.DecodeString(cfg.SecureCookieKey)

	cookie := securecookie.New(
		[]byte(cfg.SecretKey),
		secureCookieKey,
	)

	renderOptions := render.Options{
		IsDevelopment: cfg.IsDev(),
	}

	renderer := render.New(renderOptions)

	feedparser := feedparser.New()

	return &Server{
		DB:         db,
		Config:     cfg,
		Log:        log,
		Render:     renderer,
		Cookie:     cookie,
		Feedparser: feedparser,
		Mailer:     mailer,
	}
}

func (s *Server) Configure() http.Handler {
	return s.configureMiddleware(s.configureRoutes())
}

type handlerFunc func(*Server, http.ResponseWriter, *http.Request) error

type handler struct {
	*Server
	authLevel authLevel
	H         handlerFunc
}

func (s *Server) handler(authLevel authLevel, fn handlerFunc) http.Handler {
	return handler{s, authLevel, fn}
}

func (s *Server) authIgnoreHandler(fn handlerFunc) http.Handler {
	return s.handler(authLevelIgnore, fn)
}

func (s *Server) authOptionalHandler(fn handlerFunc) http.Handler {
	return s.handler(authLevelIgnore, fn)
}

func (s *Server) authRequiredHandler(fn handlerFunc) http.Handler {
	return s.handler(authLevelRequired, fn)
}

func (h handler) authorize(w http.ResponseWriter, r *http.Request) error {

	if h.authLevel == authLevelIgnore {
		return nil
	}

	var errNotAuthenticated error
	if h.authLevel == authLevelOptional {
		errNotAuthenticated = nil
	} else {
		errNotAuthenticated = httpError{
			errors.New("You must be logged in"),
			http.StatusUnauthorized,
		}
	}

	// get user from cookie
	cookie, err := r.Cookie(cookieUserID)
	if err != nil {
		return errNotAuthenticated
	}

	var userID int64

	if err := h.Cookie.Decode(cookieUserID, cookie.Value, &userID); err != nil {
		return errNotAuthenticated
	}

	if userID == 0 {
		return errNotAuthenticated
	}

	user, err := h.DB.Users.GetByID(userID)
	if err != nil {
		if isErrNoRows(err) {
			return errNotAuthenticated
		}
		return err
	}

	// all ok...
	context.Set(r, userKey, user)
	return nil
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = h.authorize(w, r); err == nil {
		err = h.H(h.Server, w, r)
	}
	if err != nil {
		h.abort(w, r, err)
	}
}

func (s *Server) setAuthCookie(w http.ResponseWriter, userID int64) {

	if encoded, err := s.Cookie.Encode(cookieUserID, userID); err == nil {
		cookie := &http.Cookie{
			Name:    cookieUserID,
			Value:   encoded,
			Expires: time.Now().Add(time.Hour * cookieTimeout),
			//Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}
}

func (s *Server) abort(w http.ResponseWriter, r *http.Request, err error) {
	switch err.(error).(type) {

	case decoders.Errors:
		s.Render.JSON(w, http.StatusBadRequest, err)

	case database.DBError:
		dbErr := err.(database.DBError)

		if dbErr.IsNoRows() {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		logger := s.Log.WithFields(logrus.Fields{
			"URL":    r.URL,
			"Method": r.Method,
			"Error":  dbErr,
			"SQL":    dbErr.Query(),
		})
		logger.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	default:
		if err == errNotAuthenticated {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		logger := s.Log.WithFields(logrus.Fields{
			"URL":    r.URL,
			"Method": r.Method,
			"Error":  err,
		})

		logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
