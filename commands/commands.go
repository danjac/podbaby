package commands

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/models"
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

	channels, err := db.Channels.SelectAll()
	if err != nil {
		panic(err)
	}

	channelHandler := func(ch *models.Channel) error {
		log.Info("Fetching from channel: " + ch.Title)
		return db.Channels.Create(ch)
	}

	podcastHandler := func(p *models.Podcast) error {
		return db.Podcasts.Create(p)
	}

	f := feedparser.New(channelHandler, podcastHandler)

	if err := f.FetchAll(channels); err != nil {
		panic(err)
	}

}
