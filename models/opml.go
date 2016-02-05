package models

import "encoding/xml"

// Outline is info on a single feed inside an OPML doc
type Outline struct {
	XMLName xml.Name `xml:"outline"`
	Type    string   `xml:"type,attr"`
	Title   string   `xml:"title,attr"`
	Text    string   `xml:"text,attr"`
	URL     string   `xml:"xmlUrl,attr"`
	HTMLURL string   `xml:"htmlUrl,attr,omitempty"`
}

// OPML encapsulates an OPML XML document
type OPML struct {
	XMLName  xml.Name   `xml:"opml"`
	Version  string     `xml:"version,attr"`
	Title    string     `xml:"head>title"`
	Outlines []*Outline `xml:"body>outline"`
}
