package feedparser

import (
	"errors"
	rss "github.com/jteeuwen/go-pkg-rss"
)

type Result struct {
	channel *rss.Channel
	items   []*rss.Item
}

var InvalidFeed = errors.New("No channel found")

// fetches normalized podcast feed
func Fetch(url string) (*Result, error) {

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
