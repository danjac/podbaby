package feedparser

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/models"
	rss "github.com/jteeuwen/go-pkg-rss"
	"strings"
)

type Result struct {
	Channel *rss.Channel
	Items   []*rss.Item
}

func (result *Result) GetWebsiteLink() string {
	for _, link := range result.Channel.Links {
		found := false
		if link.Type == "" && link.Rel == "" {
			for _, item := range result.Items {
				for _, itemLink := range item.Links {
					if itemLink.Href == link.Href {
						found = true
					}
				}
			}
			if !found {
				return link.Href
			}
		}
	}
	return ""
}

type Feedparser interface {
	FetchChannel(*models.Channel) error
	FetchAll() error
}

type defaultFeedparserImpl struct {
	DB  *database.DB
	Log *logrus.Logger
}

func New(db *database.DB, log *logrus.Logger) Feedparser {
	return &defaultFeedparserImpl{db, log}
}

func (f *defaultFeedparserImpl) FetchChannel(channel *models.Channel) error {

	result, err := fetch(channel.URL)

	if err != nil {
		return err
	}

	// update channel

	f.Log.Info("Channel:"+channel.Title, " podcasts:", len(result.Items))

	channel.Title = result.Channel.Title
	channel.Image = result.Channel.Image.Url
	channel.Description = result.Channel.Description

	link := result.GetWebsiteLink()
	if link != "" {
		channel.Website.String = result.GetWebsiteLink()
		channel.Website.Valid = true
	}

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

	if err := f.DB.Channels.Create(channel); err != nil {
		return err
	}

	for _, item := range result.Items {
		podcast := &models.Podcast{
			ChannelID:   channel.ID,
			Title:       item.Title,
			Description: item.Description,
		}
		if len(item.Enclosures) == 0 {
			continue
		}
		podcast.EnclosureURL = item.Enclosures[0].Url
		pubDate, _ := item.ParsedPubDate()
		podcast.PubDate = pubDate

		//log.Info("Podcast:" + podcast.Title)

		if err := f.DB.Podcasts.Create(podcast); err != nil {
			return err
		}
	}

	return nil

}

func (f *defaultFeedparserImpl) FetchAll() error {

	f.Log.Info("Starting podcast fetching...")

	channels, err := f.DB.Channels.SelectAll()

	if err != nil {
		return err
	}

	for _, channel := range channels {
		if err := f.FetchChannel(&channel); err != nil {
			f.Log.Error(err)
			continue
		}
	}

	return nil

}

var InvalidFeed = errors.New("No channel found")

func fetch(url string) (*Result, error) {

	var channels []*rss.Channel

	chanHandler := func(feed *rss.Feed, newChannels []*rss.Channel) {
		channels = append(channels, newChannels...)
	}

	var items []*rss.Item

	itemHandler := func(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
		items = append(items, newItems...)
	}

	feed := rss.New(5, true, chanHandler, itemHandler)

	if err := feed.Fetch(url, nil); err != nil {
		return nil, err
	}

	if len(channels) == 0 {
		return nil, InvalidFeed
	}

	result := &Result{channels[0], items}
	return result, nil
}
