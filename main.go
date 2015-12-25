package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/server"
	"github.com/jmoiron/sqlx"
	//rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	_ "github.com/lib/pq"
	"net/http"
)

var (
	env  = flag.String("env", "prod", "environment ('prod' or 'dev')")
	port = flag.String("port", "5000", "server port")
	url  = flag.String("url", "", "database connection url")
)

// should be settings
const (
	staticURL    = "/static/"
	staticDir    = "./static/"
	devServerURL = "http://localhost:8080"
)

/*
func fetchPodcasts(db *database.DB, url string) error {

	var rssChannel *rss.Channel

	chanHandler := func(feed *rss.Feed, newChannels []*rss.Channel) {
		rssChannel = newChannels[0]
	}

	var rssItems []*rss.Item

	itemHandler := func(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
		rssItems = append(rssItems, newItems...)
	}

	feed := rss.New(5, true, chanHandler, itemHandler)

	if err := feed.Fetch(url, nil); err != nil {
		return err
	}

	// tbd: check if channel already exists
	channel := &models.Channel{
		URL:         url,
		Title:       rssChannel.Title,
		Image:       rssChannel.Image.Url,
		Description: rssChannel.Description,
	}

	if err := db.Channels.Create(channel); err != nil {
		return err
	}

	// tbd: check pubdates: only insert if pub_date > MAX pub date of existing
	// items: make enclosure URL + channel ID unique

	for _, item := range rssItems {
		pubDate, _ := item.ParsedPubDate()

		pc := &models.Podcast{
			ChannelID:   channel.ID,
			Title:       item.Title,
			Description: item.Description,
			PubDate:     pubDate,
		}

		if len(item.Enclosures) == 0 {
			continue
		}

		pc.EnclosureURL = item.Enclosures[0].Url
		if err := db.Podcasts.Create(pc); err != nil {
			return err
		}

	}

	return nil
}
*/

func main() {

	flag.Parse()

	db := database.New(sqlx.MustConnect("postgres", *url))

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	log.Info("Starting web service...")

	s := server.New(db, log, &server.Config{
		StaticURL: staticURL,
		StaticDir: staticDir,
	})

	chain := alice.New(nosurf.NewPure).Then(s.Router())

	if err := http.ListenAndServe(":"+*port, chain); err != nil {
		panic(err)
	}

}
