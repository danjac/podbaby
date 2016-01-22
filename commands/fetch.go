package commands

import (
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
)

func fetchChannel(channel *models.Channel, db *database.DB, f feedparser.Feedparser) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := f.Fetch(channel); err != nil {
		return err
	}

	if err := db.Channels.Create(channel); err != nil {
		return err
	}

	if err := db.Categories.Create(channel); err != nil {
		return err
	}

	for _, p := range channel.Podcasts {
		p.ChannelID = channel.ID
		if err := db.Podcasts.Create(p); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil

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
		if err := fetchChannel(&channel, db, f); err != nil {
			log.Error(err)
			continue
		}

	}

}
