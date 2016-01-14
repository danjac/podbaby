package commands

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/server"
)

func configureLogger() *logrus.Logger {
	logger := logrus.New()

	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	return logger

}

// Serve runs the webserver
func Serve(cfg *config.Config) {

	log := configureLogger()

	db := database.MustConnect(cfg)
	defer db.Close()

	mailer, err := mailer.New(cfg)

	if err != nil {
		panic(err)
	}

	handler := server.New(db, mailer, log, cfg).Handler()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler); err != nil {
		panic(err)
	}

}

// Fetch retrieves latest podcasts
func Fetch(cfg *config.Config) {

	db := database.MustConnect(cfg)
	defer db.Close()

	log := configureLogger()

	f := feedparser.New(db, log)
	if err := f.FetchAll(); err != nil {
		panic(err)
	}

}
