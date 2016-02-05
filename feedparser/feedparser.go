package feedparser

import (
	"errors"
	"github.com/danjac/podbaby/feedparser/Godeps/_workspace/src/github.com/jinzhu/now"
	rss "github.com/danjac/podbaby/feedparser/Godeps/_workspace/src/github.com/jteeuwen/go-pkg-rss"
	"github.com/danjac/podbaby/models"
	"strings"
	"time"
)

const iTunesDTD = "http://www.itunes.com/dtds/podcast-1.0.dtd"

type keywordSet map[string]string

func (k keywordSet) add(kw string) {
	k[kw] = kw
}

func (k keywordSet) slice() []string {
	i := 0
	rv := make([]string, len(k))
	for key := range k {
		rv[i] = key
		i++
	}
	return rv
}

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

// ErrInvalidFeed is returned if RSS/Atom is invalid, missing or empty
var ErrInvalidFeed = errors.New("Invalid feed")

type result struct {
	channel *rss.Channel
	items   []*rss.Item
}

func (r *result) getWebsiteURL() string {
	// ensure we just get the top non-RSS link
	for _, link := range r.channel.Links {
		isItemLink := false
		if link.Type == "" && link.Rel == "" {
			for _, item := range r.items {
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

func (r *result) getCategories() []string {
	categories := []string{}
	for _, ext := range r.channel.Extensions {
		for _, exts := range ext {
			for _, item := range exts {
				categories = append(categories, getCategories(&item)...)
			}
		}
	}
	return categories
}

// Feedparser handles fetching of RSS/Atom podcast feeds
type Feedparser interface {
	Fetch(*models.Channel) error
}

type feedparserImpl struct{}

// New returns a Feedparser instance
func New() Feedparser {
	return &feedparserImpl{}
}

func (f *feedparserImpl) Fetch(channel *models.Channel) error {

	result, err := fetch(channel.URL)

	if err != nil {
		return err
	}

	channel.Title = result.channel.Title
	channel.Image = result.channel.Image.Url
	channel.Description = result.channel.Description

	website := result.getWebsiteURL()

	if website != "" {
		channel.Website.String = website
		channel.Website.Valid = true
	}

	channel.Categories = result.getCategories()

	// we just want unique categories
	keywords := make(keywordSet)

	for _, c := range result.channel.Categories {
		keywords.add(c.Text)
	}

	channel.Keywords.String = strings.Join(keywords.slice(), " ")
	channel.Keywords.Valid = true

	var podcasts []*models.Podcast

	for _, item := range result.items {

		podcast := &models.Podcast{
			Title:       item.Title,
			Description: item.Description,
		}

		podcast.EnclosureURL = item.Enclosures[0].Url

		if item.Guid == nil {
			// use pub date + URL as standin Guid

			podcast.GUID = item.PubDate + ":" + podcast.EnclosureURL

		} else {
			podcast.GUID = *item.Guid
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

		podcasts = append(podcasts, podcast)

	}

	channel.Podcasts = podcasts
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

func getCategories(ext *rss.Extension) []string {
	categories := []string{}
	if ext.Name != "category" {
		return categories
	}
	s := strings.Trim(ext.Attrs["text"], " ")
	if s != "" {
		categories = append(categories, s)
	}
	for _, child := range ext.Childrens {
		for _, item := range child {
			categories = append(categories, getCategories(&item)...)
		}
	}
	return categories
}

func fetch(url string) (*result, error) {

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

	result := &result{channels[0], items}
	return result, nil
}
