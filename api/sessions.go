package api

import (
	"encoding/base64"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/config"
	"net/http"
	"time"

	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/gorilla/securecookie"
)

const sessionTimeout = 24 * 30

type session interface {
	write(*echo.Context, string, interface{}) error
	read(*echo.Context, string, interface{}) (bool, error)
	readInt(*echo.Context, string) (int, bool, error)
}

type secureCookieSession struct {
	*securecookie.SecureCookie
	isSecure bool
}

func (s *secureCookieSession) write(c *echo.Context, key string, value interface{}) error {
	encoded, err := s.Encode(key, value)
	if err == nil {
		cookie := &http.Cookie{
			Name:     key,
			Value:    encoded,
			Expires:  time.Now().Add(time.Hour * sessionTimeout),
			Secure:   s.isSecure,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(c.Response(), cookie)
	}
	return err
}

func (s *secureCookieSession) read(c *echo.Context, key string, dst interface{}) (bool, error) {
	cookie, err := c.Request().Cookie(key)
	if err == http.ErrNoCookie {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	err = s.Decode(key, cookie.Value, dst)
	return err != nil, err
}

func (s *secureCookieSession) readInt(c *echo.Context, key string) (int, bool, error) {
	var rv int
	ok, err := s.read(c, key, &rv)
	return rv, ok, err
}

func newSession(cfg *config.Config) session {
	secureCookieKey, _ := base64.StdEncoding.DecodeString(cfg.SecureCookieKey)
	cookie := securecookie.New(
		[]byte(cfg.SecretKey),
		secureCookieKey,
	)
	return &secureCookieSession{
		SecureCookie: cookie,
		isSecure:     !cfg.IsDev(),
	}
}
