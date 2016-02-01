package feeder

import (
	xmlx "github.com/danjac/podbaby/commands/Godeps/_workspace/src/github.com/jteeuwen/go-pkg-xmlx"
)

type MissingRssNodeError struct{}

func (err *MissingRssNodeError) Error() string { return "Failed to find rss/rdf node in XML." }

type Extension struct {
	Name      string
	Value     string
	Attrs     map[string]string
	Childrens map[string][]Extension
}

var days = map[string]int{
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
	"Sunday":    7,
}

func (this *Feed) readRss2(doc *xmlx.Document) (err error) {
	var foundChannels []*Channel
	var ch *Channel
	var i *Item
	var n *xmlx.Node
	var list, tl []*xmlx.Node
	const ns = "*"

	root := doc.SelectNode(ns, "rss")
	if root == nil {
		root = doc.SelectNode(ns, "RDF")
	}

	if root == nil {
		return &MissingRssNodeError{}
	}

	channels := root.SelectNodes(ns, "channel")
	for _, node := range channels {
		ch = new(Channel)
		foundChannels = append(foundChannels, ch)

		ch.Title = node.S(ns, "title")
		list = node.SelectNodes(ns, "link")
		ch.Links = make([]Link, len(list))

		for i, v := range list {
			if v.Name.Space == "http://www.w3.org/2005/Atom" && v.Name.Local == "link" {
				ch.Links[i].Href = v.As("", "href")
				ch.Links[i].Rel = v.As("", "rel")
				ch.Links[i].Type = v.As("", "type")
				ch.Links[i].HrefLang = v.As("", "hreflang")
			} else {
				ch.Links[i].Href = v.GetValue()
			}
		}

		ch.Description = node.S(ns, "description")
		ch.Language = node.S(ns, "language")
		ch.Copyright = node.S(ns, "copyright")
		ch.ManagingEditor = node.S(ns, "managingEditor")
		ch.WebMaster = node.S(ns, "webMaster")
		ch.PubDate = node.S(ns, "pubDate")
		ch.LastBuildDate = node.S(ns, "lastBuildDate")
		ch.Docs = node.S(ns, "docs")

		list = node.SelectNodes(ns, "category")
		ch.Categories = make([]*Category, len(list))
		for i, v := range list {
			ch.Categories[i] = new(Category)
			ch.Categories[i].Domain = v.As(ns, "domain")
			ch.Categories[i].Text = v.GetValue()
		}

		if n = node.SelectNode(ns, "generator"); n != nil {
			ch.Generator = Generator{}
			ch.Generator.Text = n.GetValue()
		}

		ch.TTL = node.I(ns, "ttl")
		ch.Rating = node.S(ns, "rating")

		list = node.SelectNodes(ns, "hour")
		ch.SkipHours = make([]int, len(list))
		for i, v := range list {
			ch.SkipHours[i] = v.I(ns, "hour")
		}

		list = node.SelectNodes(ns, "days")
		ch.SkipDays = make([]int, len(list))
		for i, v := range list {
			ch.SkipDays[i] = days[v.GetValue()]
		}

		if n = node.SelectNode(ns, "image"); n != nil {
			ch.Image.Title = n.S(ns, "title")
			ch.Image.Url = n.S(ns, "url")
			if ch.Image.Url == "" {
				ch.Image.Url = n.As(ns, "href")
			}
			ch.Image.Link = n.S(ns, "link")
			ch.Image.Width = n.I(ns, "width")
			ch.Image.Height = n.I(ns, "height")
			ch.Image.Description = n.S(ns, "description")
		}

		if n = node.SelectNode(ns, "cloud"); n != nil {
			ch.Cloud = Cloud{}
			ch.Cloud.Domain = n.As(ns, "domain")
			ch.Cloud.Port = n.Ai(ns, "port")
			ch.Cloud.Path = n.As(ns, "path")
			ch.Cloud.RegisterProcedure = n.As(ns, "registerProcedure")
			ch.Cloud.Protocol = n.As(ns, "protocol")
		}

		if n = node.SelectNode(ns, "textInput"); n != nil {
			ch.TextInput = Input{}
			ch.TextInput.Title = n.S(ns, "title")
			ch.TextInput.Description = n.S(ns, "description")
			ch.TextInput.Name = n.S(ns, "name")
			ch.TextInput.Link = n.S(ns, "link")
		}

		list = node.SelectNodes(ns, "item")
		if len(list) == 0 {
			list = doc.SelectNodes(ns, "item")
		}

		for _, item := range list {
			i = new(Item)
			i.Title = item.S(ns, "title")
			i.Description = item.S(ns, "description")

			tl = item.SelectNodes(ns, "link")
			for _, v := range tl {
				lnk := new(Link)
				if v.Name.Space == "http://www.w3.org/2005/Atom" && v.Name.Local == "link" {
					lnk.Href = v.As("", "href")
					lnk.Rel = v.As("", "rel")
					lnk.Type = v.As("", "type")
					lnk.HrefLang = v.As("", "hreflang")
				} else {
					lnk.Href = v.GetValue()
				}

				i.Links = append(i.Links, lnk)
			}

			if n = item.SelectNode(ns, "author"); n != nil {
				i.Author.Name = n.GetValue()

			} else if n = item.SelectNode(ns, "creator"); n != nil {
				i.Author.Name = n.GetValue()
			}

			i.Comments = item.S(ns, "comments")

			guid := item.S(ns, "guid")
			if len(guid) > 0 {
				i.Guid = &guid
			}

			i.PubDate = item.S(ns, "pubDate")

			tl = item.SelectNodes(ns, "category")
			for _, lv := range tl {
				cat := new(Category)
				cat.Domain = lv.As(ns, "domain")
				cat.Text = lv.GetValue()
				i.Categories = append(i.Categories, cat)
			}

			tl = item.SelectNodes(ns, "enclosure")
			for _, lv := range tl {
				enc := new(Enclosure)
				enc.Url = lv.As(ns, "url")
				enc.Length = lv.Ai64(ns, "length")
				enc.Type = lv.As(ns, "type")
				i.Enclosures = append(i.Enclosures, enc)
			}

			if src := item.SelectNode(ns, "source"); src != nil {
				i.Source = new(Source)
				i.Source.Url = src.As(ns, "url")
				i.Source.Text = src.GetValue()
			}

			tl = item.SelectNodes("http://purl.org/rss/1.0/modules/content/", "*")
			for _, lv := range tl {
				if lv.Name.Local == "encoded" {
					i.Content = new(Content)
					i.Content.Text = lv.String()
					break
				}
			}

			i.Extensions = make(map[string]map[string][]Extension)
			for _, lv := range item.Children {
				getExtensions(&i.Extensions, lv)
			}

			ch.Items = append(ch.Items, i)
		}

		ch.Extensions = make(map[string]map[string][]Extension)
		for _, v := range node.Children {
			getExtensions(&ch.Extensions, v)
		}

	}
	this.Channels = foundChannels
	return
}

func getExtensions(extensionsX *map[string]map[string][]Extension, node *xmlx.Node) {
	extentions := *extensionsX

	extension, noErrors := getExtension(node)
	if noErrors {
		if len(extentions[node.Name.Space]) == 0 {
			extentions[node.Name.Space] = make(map[string][]Extension, 0)
		}
		if len(extentions[node.Name.Space][node.Name.Local]) == 0 {
			extentions[node.Name.Space][node.Name.Local] = make([]Extension, 0)
		}
		extentions[node.Name.Space][node.Name.Local] = append(extentions[node.Name.Space][node.Name.Local], extension)
	}
}

func getExtension(node *xmlx.Node) (Extension, bool) {
	var extension Extension
	if node.Name.Space != "" {
		extension = Extension{Name: node.Name.Local, Value: node.GetValue()}
		extension.Attrs = make(map[string]string)
		extension.Childrens = make(map[string][]Extension, 0)
		for _, x := range node.Attributes {
			extension.Attrs[x.Name.Local] = x.Value
		}
		for _, y := range node.Children {
			children, ok := getExtension(y)
			if ok {
				extension.Childrens[y.Name.Local] = append(extension.Childrens[y.Name.Local], children)
			}
		}
		return extension, true
	} else {
		return extension, false
	}
}
