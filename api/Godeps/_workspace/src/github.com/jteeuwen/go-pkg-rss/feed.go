/*
 Author: jim teeuwen <jimteeuwen@gmail.com>
 Dependencies: go-pkg-xmlx (http://github.com/jteeuwen/go-pkg-xmlx)

 This package allows us to fetch Rss and Atom feeds from the internet.
 They are parsed into an object tree which is a hybrid of both the RSS and Atom
 standards.

 Supported feeds are:
 	- Rss v0.91, 0.91 and 2.0
 	- Atom 1.0

 The package allows us to maintain cache timeout management. This prevents us
 from querying the servers for feed updates too often and risk ip bams. Appart
 from setting a cache timeout manually, the package also optionally adheres to
 the TTL, SkipDays and SkipHours values specied in the feeds themselves.

 Note that the TTL, SkipDays and SkipHour fields are only part of the RSS spec.
 For Atom feeds, we use the CacheTimeout in the Feed struct.

 Because the object structure is a hybrid between both RSS and Atom specs, not
 all fields will be filled when requesting either an RSS or Atom feed. I have
 tried to create as many shared fields as possiblem but some of them simply do
 not occur in either the RSS or Atom spec.
*/
package feeder

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	xmlx "github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/jteeuwen/go-pkg-xmlx"
)

type UnsupportedFeedError struct {
	Type    string
	Version [2]int
}

func (err *UnsupportedFeedError) Error() string {
	return fmt.Sprintf("Unsupported feed: %s, version: %+v", err.Type, err.Version)
}

type ChannelHandlerFunc func(f *Feed, newchannels []*Channel)

func (h ChannelHandlerFunc) ProcessChannels(f *Feed, newchannels []*Channel) {
	if h != nil {
		h(f, newchannels)
	}
}

type ItemHandlerFunc func(f *Feed, ch *Channel, newitems []*Item)

func (h ItemHandlerFunc) ProcessItems(f *Feed, ch *Channel, newitems []*Item) {
	if h != nil {
		h(f, ch, newitems)
	}
}

type ChannelHandler interface {
	ProcessChannels(f *Feed, newchannels []*Channel)
}

type ItemHandler interface {
	ProcessItems(f *Feed, ch *Channel, newitems []*Item)
}

type Feed struct {
	// Custom cache timeout in minutes.
	CacheTimeout int

	// Make sure we adhere to the cache timeout specified in the feed. If
	// our CacheTimeout is higher than that, we will use that instead.
	EnforceCacheLimit bool

	// Type of feed. Rss, Atom, etc
	Type string

	// Version of the feed. Major and Minor.
	Version [2]int

	// Channels with content.
	Channels []*Channel

	// Url from which this feed was created.
	Url string

	// Database
	database *database

	// The channel handler
	channelHandler ChannelHandler

	// The item handler
	itemHandler ItemHandler

	// Last time content was fetched. Used in conjunction with CacheTimeout
	// to ensure we don't get content too often.
	lastupdate time.Time

	// On our next fetch *ONLY* (this will get reset to false afterwards),
	// ignore all cache settings and update frequency hints, and always fetch.
	ignoreCacheOnce bool

	// Custom user agent
	userAgent string
}

// New is a helper function to stay semi-compatible with
// the old code. Includes the database handler to ensure
// that this approach is functionally identical to the
// old database/handlers version.
func New(cachetimeout int, enforcecachelimit bool, ch ChannelHandlerFunc, ih ItemHandlerFunc) *Feed {
	db := NewDatabase()
	f := NewWithHandlers(cachetimeout, enforcecachelimit, NewDatabaseChannelHandler(db, ch), NewDatabaseItemHandler(db, ih))
	f.database = db
	return f
}

// NewWithHandler creates a new feed with handlers.
// People should use this approach from now on.
func NewWithHandlers(cachetimeout int, enforcecachelimit bool, ch ChannelHandler, ih ItemHandler) *Feed {
	v := new(Feed)
	v.CacheTimeout = cachetimeout
	v.EnforceCacheLimit = enforcecachelimit
	v.Type = "none"
	v.channelHandler = ch
	v.itemHandler = ih
	return v
}

// This returns a timestamp of the last time the feed was updated.
func (this *Feed) LastUpdate() time.Time {
	return this.lastupdate
}

// Until the next *successful* fetching of the feed's content, the
// fetcher will ignore all cache values and update interval hints,
// and always attempt to retrieve a fresh copy of the feed.
func (this *Feed) IgnoreCacheOnce() {
	this.ignoreCacheOnce = true
}

// Fetch retrieves the feed's latest content if necessary.
//
// The charset parameter overrides the xml decoder's CharsetReader.
// This allows us to specify a custom character encoding conversion
// routine when dealing with non-utf8 input. Supply 'nil' to use the
// default from Go's xml package.
//
// This is equivalent to calling FetchClient with http.DefaultClient
func (this *Feed) Fetch(uri string, charset xmlx.CharsetFunc) (err error) {
	return this.FetchClient(uri, http.DefaultClient, charset)
}

// Fetch retrieves the feed's latest content if necessary.
//
// The charset parameter overrides the xml decoder's CharsetReader.
// This allows us to specify a custom character encoding conversion
// routine when dealing with non-utf8 input. Supply 'nil' to use the
// default from Go's xml package.
//
// The client parameter allows the use of arbitrary network connections, for
// example the Google App Engine "URL Fetch" service.
func (this *Feed) FetchClient(uri string, client *http.Client, charset xmlx.CharsetFunc) (err error) {
	if !this.CanUpdate() {
		return
	}

	this.lastupdate = time.Now().UTC()
	this.Url = uri
	doc := xmlx.New()

	if len(this.userAgent) > 1 {
		doc.SetUserAgent(this.userAgent)
	}

	if err = doc.LoadUriClient(uri, client, charset); err != nil {
		return
	}

	if err = this.makeFeed(doc); err == nil {
		// Only if fetching and parsing succeeded.
		this.ignoreCacheOnce = false
	}

	return
}

// Fetch retrieves the feed's content from the []byte
//
// The charset parameter overrides the xml decoder's CharsetReader.
// This allows us to specify a custom character encoding conversion
// routine when dealing with non-utf8 input. Supply 'nil' to use the
// default from Go's xml package.
func (this *Feed) FetchBytes(uri string, content []byte, charset xmlx.CharsetFunc) (err error) {
	this.Url = uri

	doc := xmlx.New()

	if err = doc.LoadBytes(content, charset); err != nil {
		return
	}

	return this.makeFeed(doc)
}

func (this *Feed) makeFeed(doc *xmlx.Document) (err error) {
	// Extract type and version of the feed so we can have the appropriate
	// function parse it (rss 0.91, rss 0.92, rss 2, atom etc).
	this.Type, this.Version = this.GetVersionInfo(doc)

	if ok := this.testVersions(); !ok {
		return &UnsupportedFeedError{Type: this.Type, Version: this.Version}
	}

	if err = this.buildFeed(doc); err != nil || len(this.Channels) == 0 {
		return
	}

	// reset cache timeout values according to feed specified values (TTL)
	if this.EnforceCacheLimit && this.CacheTimeout < this.Channels[0].TTL {
		this.CacheTimeout = this.Channels[0].TTL
	}

	this.notifyListeners()

	return
}

func (this *Feed) notifyListeners() {
	for _, channel := range this.Channels {
		if len(channel.Items) > 0 && this.itemHandler != nil {
			this.itemHandler.ProcessItems(this, channel, channel.Items)
		}
	}

	if len(this.Channels) > 0 && this.channelHandler != nil {
		this.channelHandler.ProcessChannels(this, this.Channels)
	}
}

// This function returns true or false, depending on whether the CacheTimeout
// value has expired or not. Additionally, it will ensure that we adhere to the
// RSS spec's SkipDays and SkipHours values (if Feed.EnforceCacheLimit is set to
// true). If this function returns true, you can be sure that a fresh feed
// update will be performed.
func (this *Feed) CanUpdate() bool {
	if this.ignoreCacheOnce {
		// Even though ignoreCacheOnce is only good for one fetch, we only reset
		// it after a successful fetch, so CanUpdate() has no side-effects, and
		// can be called repeatedly before performing the actual fetch.
		return true
	}

	// Make sure we are not within the specified cache-limit.
	// This ensures we don't request data too often.
	if this.SecondsTillUpdate() > 0 {
		return false
	}

	utc := time.Now().UTC()

	// If skipDays or skipHours are set in the RSS feed, use these to see if
	// we can update.
	if len(this.Channels) == 1 && this.Type == "rss" {
		if this.EnforceCacheLimit && len(this.Channels[0].SkipDays) > 0 {
			for _, v := range this.Channels[0].SkipDays {
				if time.Weekday(v) == utc.Weekday() {
					return false
				}
			}
		}

		if this.EnforceCacheLimit && len(this.Channels[0].SkipHours) > 0 {
			for _, v := range this.Channels[0].SkipHours {
				if v == utc.Hour() {
					return false
				}
			}
		}
	}

	return true
}

// Returns the number of seconds needed to elapse
// before the feed should update.
func (this *Feed) SecondsTillUpdate() int64 {
	utc := time.Now().UTC()
	elapsed := utc.Sub(this.lastupdate)
	return int64(this.CacheTimeout*60) - int64(elapsed.Seconds())
}

// Returns the duration needed to elapse before the feed should update.
func (this *Feed) TillUpdate() (time.Duration, error) {
	return time.ParseDuration(fmt.Sprintf("%ds", this.SecondsTillUpdate()))
}

func (this *Feed) buildFeed(doc *xmlx.Document) (err error) {
	switch this.Type {
	case "rss":
		err = this.readRss2(doc)
	case "atom":
		err = this.readAtom(doc)
	}
	return
}

func (this *Feed) testVersions() bool {
	switch this.Type {
	case "rss":
		if this.Version[0] > 2 || (this.Version[0] == 2 && this.Version[1] > 0) {
			return false
		}

	case "atom":
		if this.Version[0] > 1 || (this.Version[0] == 1 && this.Version[1] > 0) {
			return false
		}

	default:
		return false
	}

	return true
}

// Returns the type of the feed, ie. "atom" or "rss", and the version number as an array.
// The first item in the array is the major and the second the minor version number.
func (this *Feed) GetVersionInfo(doc *xmlx.Document) (ftype string, fversion [2]int) {
	var node *xmlx.Node

	if node = doc.SelectNode("http://www.w3.org/2005/Atom", "feed"); node != nil {
		ftype = "atom"
		fversion = [2]int{1, 0}
		return
	}

	if node = doc.SelectNode("", "rss"); node != nil {
		ftype = "rss"
		major := 0
		minor := 0
		version := node.As("", "version")
		p := strings.Index(version, ".")
		if p != -1 {
			major, _ = strconv.Atoi(version[0:p])
			minor, _ = strconv.Atoi(version[p+1 : len(version)])
		}
		fversion = [2]int{major, minor}
		return
	}

	// issue#5: Some documents have an RDF root node instead of rss.
	if node = doc.SelectNode("http://www.w3.org/1999/02/22-rdf-syntax-ns#", "RDF"); node != nil {
		ftype = "rss"
		fversion = [2]int{1, 1}
		return
	}

	ftype = "unknown"
	fversion = [2]int{0, 0}
	return
}

func (this *Feed) SetUserAgent(s string) {
	this.userAgent = s
}
