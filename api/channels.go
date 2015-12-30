package api

import (
	"net/http"

	"github.com/danjac/podbaby/decoders"
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

	channel := &models.Channel{
		URL: decoder.URL,
	}

	if err := s.DB.Channels.Create(channel); err != nil {
		s.abort(w, r, err)
		return
	}

	if err := s.DB.Subscriptions.Create(channel.ID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}

	go func(ch *models.Channel, user *models.User) {
		if err := s.Feedparser.FetchChannel(ch); err != nil {
			s.Log.Error(err)
			return
		}
		if err := s.DB.Subscriptions.Create(ch.ID, user.ID); err != nil {
			s.Log.Error(err)
		}
	}(channel, user)

	s.Render.Text(w, http.StatusCreated, "New channel")
}
