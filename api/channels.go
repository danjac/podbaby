package api

import (
	"net/http"

	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
)

func (s *Server) getChannelDetail(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channelID, err := getInt64(r, "id")
	if err != nil {
		s.abort(w, r, err)
		return
	}
	channel, err := s.DB.Channels.GetByID(channelID, user.ID)
	if err != nil {
		s.abort(w, r, err)
		return
	}
	detail := &models.ChannelDetail{
		Channel: channel,
	}
	podcasts, err := s.DB.Podcasts.SelectByChannelID(channelID, user.ID, getPage(r))
	if err != nil {
		s.abort(w, r, err)
		return
	}
	for _, pc := range podcasts.Podcasts {
		pc.Name = channel.Title
		pc.Image = channel.Image
		pc.ChannelID = channel.ID
		detail.Podcasts = append(detail.Podcasts, pc)
	}
	detail.Page = podcasts.Page
	s.Render.JSON(w, http.StatusOK, detail)
}

func (s *Server) getChannels(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channels, err := s.DB.Channels.SelectSubscribed(user.ID)
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, channels)
}

func (s *Server) addChannel(w http.ResponseWriter, r *http.Request) {

	decoder := &decoders.NewChannel{}

	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	user, _ := getUser(r)

	go func(url string, userID int64) {

		result, err := feedparser.Fetch(url)

		if err != nil {
			s.Log.Error(err)
			return
		}

		channel := &models.Channel{
			URL:         url,
			Title:       result.Channel.Title,
			Image:       result.Channel.Image.Url,
			Description: result.Channel.Description,
		}

		if err := s.DB.Channels.Create(channel); err != nil {
			s.Log.Error(err)
			return
		}

		if err := s.DB.Subscriptions.Create(channel.ID, userID); err != nil {
			s.Log.Error(err)
			return
		}

		for _, item := range result.Items {
			podcast := &models.Podcast{
				ChannelID:   channel.ID,
				Title:       item.Title,
				Description: item.Description,
			}
			if len(item.Enclosures) == 0 {
				s.Log.Debug("Item has no enclosures")
				continue
			}
			podcast.EnclosureURL = item.Enclosures[0].Url
			pubDate, _ := item.ParsedPubDate()
			podcast.PubDate = pubDate

			if err := s.DB.Podcasts.Create(podcast); err != nil {
				s.Log.Error(err)
				return
			}
		}

	}(decoder.URL, user.ID)

	s.Render.Text(w, http.StatusCreated, "New channel")
}
