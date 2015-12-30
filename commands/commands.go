package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/api"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// should be settings
const (
	defaultStaticURL = "/static/"
	defaultStaticDir = "./static/"
	devStaticURL     = "http://localhost:8080/static/"
)

// Serve runs the webserver
func Serve(url string, port int, secretKey, env string) {

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	log.Info("Starting web service...")

	db := database.New(sqlx.MustConnect("postgres", url))
	defer db.Close()

	var staticURL string
	if env == "dev" {
		staticURL = devStaticURL
	} else {
		staticURL = defaultStaticURL
	}

	api := api.New(db, log, &api.Config{
		StaticURL: staticURL,
		StaticDir: defaultStaticDir,
		SecretKey: secretKey,
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), api.Handler()); err != nil {
		panic(err)
	}

}

// Fetch retrieves latest podcasts
func Fetch(url string) {

	db := database.New(sqlx.MustConnect("postgres", url))
	defer db.Close()

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

		log.Info("Channel:"+channel.Title, " podcasts:", len(result.Items))

		channel.Title = result.Channel.Title
		channel.Image = result.Channel.Image.Url
		channel.Description = result.Channel.Description

		// we just want unique categories
		categoryMap := make(map[string]string)

		for _, category := range result.Channel.Categories {
			categoryMap[category.Text] = category.Text
		}

		var categories []string
		for _, category := range categoryMap {
			categories = append(categories, category)
		}

		channel.Categories.String = strings.Join(categories, " ")
		channel.Categories.Valid = true

		if err := db.Channels.Create(&channel); err != nil {
			log.Error(err)
			continue
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

			//log.Info("Podcast:" + podcast.Title)

			if err := db.Podcasts.Create(podcast); err != nil {
				log.Error(err)
				continue
			}
		}

	}

}
