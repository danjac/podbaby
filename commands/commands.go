package commands

import (
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/server"
	"github.com/jmoiron/sqlx"
)

func mustConnect(url string) *database.DB {
	return database.New(sqlx.MustConnect("postgres", url))
}

// Serve runs the webserver
func Serve(cfg *config.Config) {

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	log.Info("Starting web service...")

	db := mustConnect(cfg.DatabaseURL)
	defer db.Close()

	mailer := mailer.New(
		cfg.Mail.Addr,
		smtp.PlainAuth(
			cfg.Mail.ID,
			cfg.Mail.User,
			cfg.Mail.Password,
			cfg.Mail.Host,
		),
	)

	s := server.New(db, mailer, log, cfg)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), s.Handler()); err != nil {
		panic(err)
	}

}

// Fetch retrieves latest podcasts
func Fetch(cfg *config.Config) {

	db := mustConnect(cfg.DatabaseURL)
	defer db.Close()

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	f := feedparser.New(db, log)
	if err := f.FetchAll(); err != nil {
		panic(err)
	}

}
