package feeder

import (
	"crypto/md5"
	"io"
	"time"
)

type Item struct {
	// RSS and Shared fields
	Title       string
	Links       []*Link
	Description string
	Author      Author
	Categories  []*Category
	Comments    string
	Enclosures  []*Enclosure
	Guid        *string
	PubDate     string
	Source      *Source

	// Atom specific fields
	Id           string
	Generator    *Generator
	Contributors []string
	Content      *Content
	Updated      string

	Extensions map[string]map[string][]Extension
}

func (i *Item) ParsedPubDate() (time.Time, error) {
	return parseTime(i.PubDate)
}

func (i *Item) Key() string {
	switch {
	case i.Guid != nil && len(*i.Guid) != 0:
		return *i.Guid
	case len(i.Id) != 0:
		return i.Id
	case len(i.Title) > 0 && len(i.PubDate) > 0:
		return i.Title + i.PubDate
	default:
		h := md5.New()
		io.WriteString(h, i.Description)
		return string(h.Sum(nil))
	}
}
