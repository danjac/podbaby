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

func handleBatch(batch []models.Channel, s store.Store, f feedparser.Feedparser) {

	var wg sync.WaitGroup
	wg.Add(len(batch))

	for _, channel := range batch {
		go func(channel models.Channel) {
			defer wg.Done()
			if err := fetchChannel(&channel, s, f); err != nil {
				log.Printf("Error fetching channel %s: %v", channel.Title, err)
			}
		}(channel)
	}

	wg.Wait()
}

func fetchChannel(channel *models.Channel, s store.Store, f feedparser.Feedparser) error {

	if channel.URL == "" {
		return nil
	}

	log.Printf("Channel: %s", channel.Title)

	channelStore := s.Channels()

	if err := f.Fetch(channel); err != nil {
		return err
	}

	tx, err := s.Conn().Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err = channelStore.CreateOrUpdate(tx, channel); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil

}

// Fetch retrieves latest podcasts
func Fetch(cfg *config.Config) {

	log.Println("Starting fetch")

	start := time.Now()

	s, err := store.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()

	conn := s.Conn()

	var channels []models.Channel
	if err := s.Channels().SelectAll(conn, &channels); err != nil {
		log.Fatalln(err)
	}
	numChannels := len(channels)

	log.Printf("%d channels to fetch", numChannels)

	if err != nil {
		log.Fatalln(err)
	}

	f := feedparser.New()

	batchSize := cfg.MaxDBConnections / 10
	batch := make([]models.Channel, batchSize)

	for i, channel := range channels {

		batch = append(batch, channel)

		if i > 0 && (i%batchSize == 0 || i == numChannels-1) {
			handleBatch(batch, s, f)
			batch = make([]models.Channel, batchSize)
		}

	}
	log.Printf("Fetch completed, %d channels fetched in %v", numChannels, time.Since(start))

}
