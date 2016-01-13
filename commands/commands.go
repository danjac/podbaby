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
)

// Serve runs the webserver
func Serve(cfg *config.Config) {

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	log.Info("Starting web service...")

	db := database.MustConnect(cfg)
	defer db.Close()

	mailer, err := mailer.New(
		cfg.Mail.Addr,
		smtp.PlainAuth(
			cfg.Mail.ID,
			cfg.Mail.User,
			cfg.Mail.Password,
			cfg.Mail.Host,
		),
		"./templates/email",
	)

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
