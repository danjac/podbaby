package commands

import (
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
)

func fetchChannel(channel *models.Channel, store store.Store, f feedparser.Feedparser) error {

	channelStore := store.Channels()

	tx, err := store.Conn().Begin()
	if err != nil {
		return err
	}

	if err := f.Fetch(channel); err != nil {
		return err
	}

	if err := channelStore.AddCategories(tx, channel); err != nil {
		return err
	}

	if err := channelStore.AddPodcasts(tx, channel); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil

}

// Fetch retrieves latest podcasts
func Fetch(cfg *config.Config) {

	store, err := store.New(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Conn().Close()

	channels, err := store.Channels().SelectAll(store.Conn())

	if err != nil {
		panic(err)
	}

	f := feedparser.New()

	for _, channel := range channels {

		if err := fetchChannel(&channel, store, f); err != nil {
			continue
		}

	}

}
