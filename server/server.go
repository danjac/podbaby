package server

import (
	"database/sql"
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
	"github.com/danjac/podbaby/models"
	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/unrolled/render"
)

const (
	cookieUserID  = "userid"
	userKey       = "user"
	cookieTimeout = 24
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

	f := feedparser.New(db, log)

	return &Server{
		DB:         db,
		Config:     cfg,
		Log:        log,
		Render:     renderer,
		Cookie:     cookie,
		Feedparser: f,
		Mailer:     mailer,
	}
}

func (s *Server) Handler() http.Handler {
	return s.configureMiddleware(s.configureRoutes())
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

func (s *Server) requireAuth(fn http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if user already set elsewhere
		if _, ok := getUser(r); ok {
			fn(w, r)
			return
		}
		// get user from cookie
		user, err := s.getUserFromCookie(r)
		if err != nil {
			s.abort(w, r, err)
			return
		}
		// all ok...
		context.Set(r, userKey, user)
		fn(w, r)
	})

}

func (s *Server) getUserFromCookie(r *http.Request) (*models.User, error) {

	cookie, err := r.Cookie(cookieUserID)
	if err != nil {
		return nil, errNotAuthenticated
	}

	var userID int64

	if err := s.Cookie.Decode(cookieUserID, cookie.Value, &userID); err != nil {
		return nil, errNotAuthenticated
	}

	if userID == 0 {
		return nil, errNotAuthenticated
	}

	user, err := s.DB.Users.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errNotAuthenticated
		}
		return nil, err
	}
	return user, nil

}

func (s *Server) abort(w http.ResponseWriter, r *http.Request, err error) {
	if err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err == errNotAuthenticated {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	switch err.(error).(type) {

	case decoders.Errors:
		s.Render.JSON(w, http.StatusBadRequest, err)
	default:
		logger := s.Log.WithFields(logrus.Fields{
			"URL":    r.URL,
			"Method": r.Method,
			"Error":  err,
		})

		logger.Error(err)
		http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
	}
}
