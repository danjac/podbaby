package commands

import (
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
	"log"
	"time"
)

func fetchChannel(channel *models.Channel, store store.Store, f feedparser.Feedparser) error {

	log.Printf("Channel: %s", channel.Title)

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

	log.Println("Starting fetch")

	start := time.Now()

	store, err := store.New(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Conn().Close()

	channels, err := store.Channels().SelectAll(store.Conn())
	numChannels := len(channels)
	log.Printf("%d channels to fetch", numChannels)

	if err != nil {
		panic(err)
	}

	f := feedparser.New()

	for _, channel := range channels {

		if err := fetchChannel(&channel, store, f); err != nil {
			continue
		}

	}
	log.Printf("Fetch completed, %d channels fetched in %v", numChannels, time.Since(start))

}
