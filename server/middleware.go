package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"net/http"
	"net/http/httptest"
	"time"
)

type timerMiddleware struct {
	log     *logrus.Logger
	handler http.Handler
}

func newTimerMiddleware(logger *logrus.Logger) alice.Constructor {
	return func(handler http.Handler) http.Handler {
		return &timerMiddleware{
			logger,
			handler,
		}
	}
}

func (m *timerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	// we need to record the original content
	// otherwise headers are sent before we can add our own
	rec := httptest.NewRecorder()

	// save the original content
	m.handler.ServeHTTP(rec, r)

	// copy response headers
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	timeTaken := time.Since(start)

	if m.log != nil {

		logger := m.log.WithFields(logrus.Fields{
			"URL":    r.URL.Path,
			"Method": r.Method,
			"Time":   timeTaken,
		})

		logger.Info()
	}

	w.Header().Set("X-Response-Time", timeTaken.String())

	// set the correct status code
	w.WriteHeader(rec.Code)

	// write the original content
	w.Write(rec.Body.Bytes())

}

func (s *Server) configureMiddleware(handler http.Handler) http.Handler {
	var middleware = []alice.Constructor{
		nosurf.NewPure,
	}

	// just log requests in development
	var log *logrus.Logger
	if s.Config.IsDev() {
		log = s.Log
	}
	middleware = append(middleware, newTimerMiddleware(log))

	return alice.New(middleware...).Then(handler)

}
