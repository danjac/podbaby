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

	handler := server.New(db, mailer, log, cfg).Configure()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler); err != nil {
		panic(err)
	}

}

// Fetch retrieves latest podcasts
func Fetch(cfg *config.Config) {

	db := database.MustConnect(cfg)
	defer db.Close()

	log := configureLogger()
	log.Info("fetching...")

	channels, err := db.Channels.SelectAll()

	if err != nil {
		panic(err)
	}

	f := feedparser.New()

	for _, channel := range channels {

		log.Info("Channel:" + channel.Title)

		if err := f.Fetch(&channel); err != nil {
			log.Error(err)
			continue
		}

		if err := db.Channels.Create(&channel); err != nil {
			log.Error(err)
			continue
		}

		for _, p := range channel.Podcasts {
			p.ChannelID = channel.ID
			if err := db.Podcasts.Create(p); err != nil {
				log.Error(err)
				continue
			}
		}

	}

}
