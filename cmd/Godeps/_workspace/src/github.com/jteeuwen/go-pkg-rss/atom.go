package feeder

import xmlx "github.com/danjac/podbaby/cmd/Godeps/_workspace/src/github.com/jteeuwen/go-pkg-xmlx"

func (this *Feed) readAtom(doc *xmlx.Document) (err error) {
	ns := "http://www.w3.org/2005/Atom"
	channels := doc.SelectNodes(ns, "feed")

	var foundChannels []*Channel
	var ch *Channel
	var i *Item
	var tn *xmlx.Node
	var list []*xmlx.Node

	for _, node := range channels {
		ch = new(Channel)
		foundChannels = append(foundChannels, ch)

		ch.Title = node.S(ns, "title")
		ch.LastBuildDate = node.S(ns, "updated")
		ch.Id = node.S(ns, "id")
		ch.Rights = node.S(ns, "rights")

		list = node.SelectNodes(ns, "link")
		ch.Links = make([]Link, len(list))
		for i, v := range list {
			ch.Links[i].Href = v.As("", "href")
			ch.Links[i].Rel = v.As("", "rel")
			ch.Links[i].Type = v.As("", "type")
			ch.Links[i].HrefLang = v.As("", "hreflang")
		}

		if tn = node.SelectNode(ns, "subtitle"); tn != nil {
			ch.SubTitle = SubTitle{}
			ch.SubTitle.Type = tn.As("", "type")
			ch.SubTitle.Text = tn.GetValue()
		}

		if tn = node.SelectNode(ns, "generator"); tn != nil {
			ch.Generator = Generator{}
			ch.Generator.Uri = tn.As("", "uri")
			ch.Generator.Version = tn.As("", "version")
			ch.Generator.Text = tn.GetValue()
		}

		if tn = node.SelectNode(ns, "author"); tn != nil {
			ch.Author = Author{}
			ch.Author.Name = tn.S("", "name")
			ch.Author.Uri = tn.S("", "uri")
			ch.Author.Email = tn.S("", "email")
		}

		list = node.SelectNodes(ns, "entry")

		for _, item := range list {
			i = new(Item)
			i.Title = item.S(ns, "title")
			i.Id = item.S(ns, "id")
			i.PubDate = item.S(ns, "published")
			i.Updated = item.S(ns, "updated")
			i.Description = item.S(ns, "summary")

			links := item.SelectNodes(ns, "link")
			for _, lv := range links {
				if lv.As(ns, "rel") == "enclosure" {
					enc := new(Enclosure)
					enc.Url = lv.As("", "href")
					enc.Type = lv.As("", "type")
					i.Enclosures = append(i.Enclosures, enc)
				} else {
					lnk := new(Link)
					lnk.Href = lv.As("", "href")
					lnk.Rel = lv.As("", "rel")
					lnk.Type = lv.As("", "type")
					lnk.HrefLang = lv.As("", "hreflang")
					i.Links = append(i.Links, lnk)
				}
			}

			list = item.SelectNodes(ns, "contributor")
			for _, cv := range list {
				i.Contributors = append(i.Contributors, cv.S("", "name"))
			}

			list = item.SelectNodes(ns, "category")
			for _, cv := range list {
				cat := new(Category)
				cat.Domain = ""
				cat.Text = cv.As("", "term")
				i.Categories = append(i.Categories, cat)
			}

			if tn = item.SelectNode(ns, "content"); tn != nil {
				i.Content = new(Content)
				i.Content.Type = tn.As("", "type")
				i.Content.Lang = tn.S("xml", "lang")
				i.Content.Base = tn.S("xml", "base")
				i.Content.Text = tn.GetValue()
			}

			if tn = item.SelectNode(ns, "author"); tn != nil {
				i.Author = Author{}
				i.Author.Name = tn.S(ns, "name")
				i.Author.Uri = tn.S(ns, "uri")
				i.Author.Email = tn.S(ns, "email")
			}

			ch.Items = append(ch.Items, i)
		}
	}
	this.Channels = foundChannels
	return
}
