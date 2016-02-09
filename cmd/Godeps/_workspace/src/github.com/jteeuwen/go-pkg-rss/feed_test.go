package feeder

import (
	"io/ioutil"
	"testing"
)

var items []*Item

func TestFeed(t *testing.T) {
	urilist := []string{
		//"http://cyber.law.harvard.edu/rss/examples/sampleRss091.xml", // Non-utf8 encoding.
		"http://store.steampowered.com/feeds/news.xml", // This feed violates the rss spec.
		"http://cyber.law.harvard.edu/rss/examples/sampleRss092.xml",
		"http://cyber.law.harvard.edu/rss/examples/rss2sample.xml",
		"http://blog.case.edu/news/feed.atom",
	}

	var feed *Feed
	var err error

	for _, uri := range urilist {
		feed = New(5, true, chanHandler, itemHandler)

		if err = feed.Fetch(uri, nil); err != nil {
			t.Errorf("%s >>> %s", uri, err)
			return
		}
	}
}

func Test_NoHandlers(t *testing.T) {
	feed := New(1, true, nil, nil)
	content, _ := ioutil.ReadFile("testdata/initial.atom")
	err := feed.FetchBytes("http://example.com", content, nil)
	if err != nil {
		t.Error(err)
	}
}

func Test_NewItem(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/initial.atom")
	feed := New(1, true, chanHandler, itemHandler)
	err := feed.FetchBytes("http://example.com", content, nil)
	if err != nil {
		t.Error(err)
	}

	content, _ = ioutil.ReadFile("testdata/initial_plus_one_new.atom")
	feed.FetchBytes("http://example.com", content, nil)
	expected := "Second title"
	if len(items) != 1 {
		t.Errorf("Expected %s new item, got %s", 1, len(items))
	}

	if expected != items[0].Title {
		t.Errorf("Expected %s, got %s", expected, items[0].Title)
	}
}

func Test_AtomAuthor(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/idownload.atom")
	if err != nil {
		t.Errorf("unable to load file")
	}
	feed := New(1, true, chanHandler, itemHandler)
	err = feed.FetchBytes("http://example.com", content, nil)

	item := feed.Channels[0].Items[0]
	expected := "Cody Lee"
	if item.Author.Name != expected {
		t.Errorf("Expected author to be %s but found %s", expected, item.Author.Name)
	}
}

func Test_RssAuthor(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/boing.rss")
	feed := New(1, true, chanHandler, itemHandler)
	feed.FetchBytes("http://example.com", content, nil)

	item := feed.Channels[0].Items[0]
	expected := "Cory Doctorow"
	if item.Author.Name != expected {
		t.Errorf("Expected author to be %s but found %s", expected, item.Author.Name)
	}
}

func Test_ItemExtensions(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/extension.rss")
	feed := New(1, true, chanHandler, itemHandler)
	feed.FetchBytes("http://example.com", content, nil)

	edgarExtensionxbrlFiling := feed.Channels[0].Items[0].Extensions["http://www.sec.gov/Archives/edgar"]["xbrlFiling"][0].Childrens
	companyExpected := "Cellular Biomedicine Group, Inc."
	companyName := edgarExtensionxbrlFiling["companyName"][0]
	if companyName.Value != companyExpected {
		t.Errorf("Expected company to be %s but found %s", companyExpected, companyName.Value)
	}

	files := edgarExtensionxbrlFiling["xbrlFiles"][0].Childrens["xbrlFile"]
	fileSizeExpected := 10
	if len(files) != 10 {
		t.Errorf("Expected files size to be %s but found %s", fileSizeExpected, len(files))
	}

	file := files[0]
	fileExpected := "cbmg_10qa.htm"
	if file.Attrs["file"] != fileExpected {
		t.Errorf("Expected file to be %s but found %s", fileExpected, len(file.Attrs["file"]))
	}
}

func Test_ChannelExtensions(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/extension.rss")
	feed := New(1, true, chanHandler, itemHandler)
	feed.FetchBytes("http://example.com", content, nil)

	channel := feed.Channels[0]
	itunesExtentions := channel.Extensions["http://www.itunes.com/dtds/podcast-1.0.dtd"]

	authorExptected := "The Author"
	ownerEmailExpected := "test@rss.com"
	categoryExpected := "Politics"
	imageExptected := "http://golang.org/doc/gopher/project.png"

	if itunesExtentions["author"][0].Value != authorExptected {
		t.Errorf("Expected author to be %s but found %s", authorExptected, itunesExtentions["author"][0].Value)
	}

	if itunesExtentions["owner"][0].Childrens["email"][0].Value != ownerEmailExpected {
		t.Errorf("Expected owner email to be %s but found %s", ownerEmailExpected, itunesExtentions["owner"][0].Childrens["email"][0].Value)
	}

	if itunesExtentions["category"][0].Attrs["text"] != categoryExpected {
		t.Errorf("Expected category text to be %s but found %s", categoryExpected, itunesExtentions["category"][0].Attrs["text"])
	}

	if itunesExtentions["image"][0].Attrs["href"] != imageExptected {
		t.Errorf("Expected image href to be %s but found %s", imageExptected, itunesExtentions["image"][0].Attrs["href"])
	}
}

func Test_CData(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/iosBoardGameGeek.rss")
	feed := New(1, true, chanHandler, itemHandler)
	feed.FetchBytes("http://example.com", content, nil)

	item := feed.Channels[0].Items[0]
	expected := `<p>abc<div>"def"</div>ghi`
	if item.Description != expected {
		t.Errorf("Expected item.Description to be [%s] but item.Description=[%s]", expected, item.Description)
	}
}

func Test_Link(t *testing.T) {
	content, _ := ioutil.ReadFile("testdata/nytimes.rss")
	feed := New(1, true, chanHandler, itemHandler)
	feed.FetchBytes("http://example.com", content, nil)

	channel := feed.Channels[0]
	item := channel.Items[0]

	channelLinkExpected := "http://www.nytimes.com/services/xml/rss/nyt/US.xml"
	itemLinkExpected := "http://www.nytimes.com/2014/01/18/technology/in-keeping-grip-on-data-pipeline-obama-does-little-to-reassure-industry.html?partner=rss&emc=rss"

	if channel.Links[0].Href != channelLinkExpected {
		t.Errorf("Expected author to be %s but found %s", channelLinkExpected, channel.Links[0].Href)
	}

	if item.Links[0].Href != itemLinkExpected {
		t.Errorf("Expected author to be %s but found %s", itemLinkExpected, item.Links[0].Href)
	}
}

func chanHandler(feed *Feed, newchannels []*Channel) {
	println(len(newchannels), "new channel(s) in", feed.Url)
}

func itemHandler(feed *Feed, ch *Channel, newitems []*Item) {
	items = newitems
	println(len(newitems), "new item(s) in", ch.Title, "of", feed.Url)
}
