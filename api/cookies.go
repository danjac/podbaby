package api

import (
	"encoding/base64"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/config"
	"net/http"
	"time"

	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/gorilla/securecookie"
)

const cookieTimeout = 24

type CookieStore interface {
	Write(*echo.Context, string, interface{}) error
	Read(*echo.Context, string, interface{}) error
}

type secureCookieStore struct {
	*securecookie.SecureCookie
	isSecure bool
}

func (s *secureCookieStore) Write(c *echo.Context, key string, value interface{}) error {
	encoded, err := s.Encode(key, value)
	if err == nil {
		cookie := &http.Cookie{
			Name:     key,
			Value:    encoded,
			Expires:  time.Now().Add(time.Hour * cookieTimeout),
			Secure:   s.isSecure,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(c.Response(), cookie)
	}
	return err
}

func (s *secureCookieStore) Read(c *echo.Context, key string, dst interface{}) error {
	cookie, err := c.Request().Cookie(key)
	if err != nil {
		return err
	}
	return s.Decode(key, cookie.Value, dst)
}

func newCookieStore(cfg *config.Config) CookieStore {
	secureCookieKey, _ := base64.StdEncoding.DecodeString(cfg.SecureCookieKey)
	cookie := securecookie.New(
		[]byte(cfg.SecretKey),
		secureCookieKey,
	)
	return &secureCookieStore{
		SecureCookie: cookie,
		isSecure:     !cfg.IsDev(),
	}
}
