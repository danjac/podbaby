package models

import "encoding/xml"

type Outline struct {
	XMLName xml.Name `xml:"outline"`
	Type    string   `xml:"type,attr"`
	Title   string   `xml:"title,attr"`
	Text    string   `xml:"text,attr"`
	URL     string   `xml:"xmlUrl,attr"`
	HtmlURL string   `xml:"htmlUrl,attr,omitempty"`
}

type OPML struct {
	XMLName  xml.Name   `xml:"opml"`
	Version  string     `xml:"version,attr"`
	Title    string     `xml:"head>title"`
	Outlines []*Outline `xml:"body>outline"`
}
