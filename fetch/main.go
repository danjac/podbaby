package main

import (
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"

	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var url = flag.String("url", "", "database connection url")

func main() {

	flag.Parse()

	db := database.New(sqlx.MustConnect("postgres", *url))

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	log.Info("Starting podcast fetching...")

	channels, err := db.Channels.GetAll()

	if err != nil {
		panic(err)
	}

	for _, channel := range channels {

		result, err := feedparser.Fetch(channel.URL)

		if err != nil {
			log.Error(err)
			continue
		}

		// update channel

		log.Info("Channel:" + channel.Title)

		channel.Title = result.Channel.Title
		channel.Image = result.Channel.Image.Url
		channel.Description = result.Channel.Description

		if err := db.Channels.Create(&channel); err != nil {
			log.Error(err)
			return
		}

		for _, item := range result.Items {
			podcast := &models.Podcast{
				ChannelID:   channel.ID,
				Title:       item.Title,
				Description: item.Description,
			}
			if len(item.Enclosures) == 0 {
				log.Debug("Item has no enclosures")
				continue
			}
			podcast.EnclosureURL = item.Enclosures[0].Url
			pubDate, _ := item.ParsedPubDate()
			podcast.PubDate = pubDate

			log.Info("Podcast:" + podcast.Title)

			if err := db.Podcasts.Create(podcast); err != nil {
				log.Error(err)
				continue
			}
		}

	}

}
