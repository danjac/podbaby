package commands

import (
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
	"log"
	"sync"
	"time"
)

func fetchChannel(channel *models.Channel, store store.Store, f feedparser.Feedparser) error {

	log.Printf("Channel: %s", channel.Title)

	channelStore := store.Channels()

	tx, err := store.Conn().Begin()

	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

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
		log.Fatalln(err)
	}
	defer store.Conn().Close()

	var channels []models.Channel
	if err := store.Channels().SelectAll(store.Conn(), &channels); err != nil {
		log.Fatalln(err)
	}
	numChannels := len(channels)
	log.Printf("%d channels to fetch", numChannels)

	if err != nil {
		log.Fatalln(err)
	}

	f := feedparser.New()

	var wg sync.WaitGroup

	wg.Add(numChannels)

	for _, channel := range channels {

		go func(channel models.Channel) {
			defer wg.Done()
			if err := fetchChannel(&channel, store, f); err != nil {
				log.Printf("Error fetching channel %s: %v", channel.Title, err)
			}
		}(channel)

	}
	wg.Wait()
	log.Printf("Fetch completed, %d channels fetched in %v", numChannels, time.Since(start))

}
