package api

import (
	"encoding/base64"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/config"
	"net/http"
	"time"

	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/gorilla/securecookie"
)

const sessionTimeout = 24

type Session interface {
	Write(*echo.Context, string, interface{}) error
	Read(*echo.Context, string) (interface{}, error)
}

type secureCookieSession struct {
	*securecookie.SecureCookie
	isSecure bool
}

func (s *secureCookieSession) Write(c *echo.Context, key string, value interface{}) error {
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

func (s *secureCookieSession) Read(c *echo.Context, key string) (interface{}, error) {
	cookie, err := c.Request().Cookie(key)
	if err != nil {
		return nil, err
	}
	var dst interface{}
	err = s.Decode(key, cookie.Value, dst)
	return dst, err
}

func newSession(cfg *config.Config) Session {
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
