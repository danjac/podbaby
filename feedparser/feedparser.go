package feedparser

import (
	"errors"
	"github.com/danjac/podbaby/models"
	"github.com/jinzhu/now"
	rss "github.com/jteeuwen/go-pkg-rss"
	"strings"
	"time"
)

func init() {
	// support for additional pub date formats we've found
	formats := []string{
		"Mon, 02 Jan 2006 15:04:05 +0000",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 5 January 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 January 2006 15:04:05 MST",
		"Mon, 02 January 2006 15:04:05 +0000",
		"Mon, 02 January 2006 15:04:05 -0700",
		"Mon, 2 January 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04 +0000",
		"Mon, 02 Jan 2006 15:04 +0000",
		"Monday, 2 January 2006 15:04:05 MST",
		"Monday, 2 Jan 2006 15:04:05 MST",
		"2 January 2006 15:04:05 MST",
		"2 Jan 2006 15:04:05 MST",
	}
	for _, format := range formats {
		now.TimeFormats = append(now.TimeFormats, format)
	}
}

var ErrInvalidFeed = errors.New("Invalid feed")

type ChannelHandler func(*models.Channel) error
type PodcastHandler func(*models.Podcast) error

type Result struct {
	Channel *rss.Channel
	Items   []*rss.Item
}

func (result *Result) getWebsiteURL() string {
	// ensure we just get the top non-RSS link
	for _, link := range result.Channel.Links {
		isItemLink := false
		if link.Type == "" && link.Rel == "" {
			for _, item := range result.Items {
				for _, itemLink := range item.Links {
					if itemLink.Href == link.Href {
						isItemLink = true
					}
				}
			}
			if !isItemLink {
				return link.Href
			}
		}
	}
	return ""
}

type Feedparser interface {
	FetchChannel(*models.Channel) error
	FetchAll([]models.Channel) error
}

type defaultFeedparserImpl struct {
	ch ChannelHandler
	ph PodcastHandler
}

func New(ch ChannelHandler, ph PodcastHandler) Feedparser {
	return &defaultFeedparserImpl{ch, ph}
}

func (f *defaultFeedparserImpl) FetchChannel(channel *models.Channel) error {

	result, err := fetch(channel.URL)

	if err != nil {
		return err
	}

	channel.Title = result.Channel.Title
	channel.Image = result.Channel.Image.Url
	channel.Description = result.Channel.Description

	website := result.getWebsiteURL()

	if website != "" {
		channel.Website.String = website
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

	if err := f.ch(channel); err != nil {
		return err
	}

	for _, item := range result.Items {

		podcast := &models.Podcast{
			ChannelID:   channel.ID,
			Title:       item.Title,
			Description: item.Description,
		}

		podcast.EnclosureURL = item.Enclosures[0].Url

		if item.Guid == nil {
			// use pub date + URL as standin Guid

			podcast.Guid = item.PubDate + ":" + podcast.EnclosureURL

		} else {
			podcast.Guid = *item.Guid
		}

		if item.Source != nil {
			podcast.Source = item.Source.Url
		}

		var pubDate time.Time

		// try using the builtin RSS parser first
		if pubDate, err = item.ParsedPubDate(); err != nil {
			// try some other parsers
			pubDate, err = now.Parse(item.PubDate)
			// pubdate will be "empty", we'll have to live with that
		}
		podcast.PubDate = pubDate

		if err = f.ph(podcast); err != nil {
			return err
		}
	}

	return nil

}

func (f *defaultFeedparserImpl) FetchAll(channels []models.Channel) error {

	for _, channel := range channels {
		if err := f.FetchChannel(&channel); err != nil {
			continue
		}
	}

	return nil

}

func isPlayable(mediaType string) bool {

	prefixes := []string{
		"audio",
		"video",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(mediaType, prefix) {
			return true
		}
	}

	return false
}

func fetch(url string) (*Result, error) {

	var channels []*rss.Channel

	chanHandler := func(feed *rss.Feed, newChannels []*rss.Channel) {
		channels = append(channels, newChannels...)
	}

	var items []*rss.Item

	itemHandler := func(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
		for _, item := range newItems {
			// only include items with enclosures
			for _, enclosure := range item.Enclosures {
				if isPlayable(enclosure.Type) {
					items = append(items, item)
				}
			}
		}
	}

	feed := rss.New(5, true, chanHandler, itemHandler)

	if err := feed.Fetch(url, nil); err != nil {
		return nil, ErrInvalidFeed
	}

	if len(channels) == 0 {
		return nil, ErrInvalidFeed
	}

	if len(items) == 0 {
		return nil, ErrInvalidFeed
	}

	result := &Result{channels[0], items}
	return result, nil
}
