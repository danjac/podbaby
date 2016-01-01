package main

/*
This is a minimal sample application, demonstrating how to set up an RSS feed
for regular polling of new channels/items.

Build & run with:

 $ go run handlerexample.go

*/

import (
	"fmt"
	rss "github.com/danjac/podbaby/Godeps/_workspace/src/github.com/jteeuwen/go-pkg-rss"
	"os"
	"time"
)

func main() {
	// This sets up a new feed and polls it for new channels/items.
	// Invoke it with 'go PollFeed(...)' to have the polling performed in a
	// separate goroutine, so you can continue with the rest of your program.
	PollFeed("http://blog.case.edu/news/feed.atom", 5)
}

func PollFeed(uri string, timeout int) {

	handlers := &MyHandlers{}
	feed := rss.NewWithHandlers(timeout, true, handlers, handlers)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		<-time.After(time.Duration(10 * time.Second))
	}
}

type MyHandlers struct{}

func (m *MyHandlers) ProcessChannels(feed *rss.Feed, newchannels []*rss.Channel) {
	fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (m *MyHandlers) ProcessItems(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	fmt.Printf("%d new rad item(s) in %s\n", len(newitems), feed.Url)
}
