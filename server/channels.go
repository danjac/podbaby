package server

import (
	"net/http"

	"database/sql"
	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
)

func (s *Server) getChannelDetail(w http.ResponseWriter, r *http.Request) {
	channelID, _ := getInt64(r, "id")

	channel, err := s.DB.Channels.GetByID(channelID)
	if err != nil {
		s.abort(w, r, err)
		return
	}
	detail := &models.ChannelDetail{
		Channel: channel,
	}
	podcasts, err := s.DB.Podcasts.SelectByChannelID(channelID, getPage(r))
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
		s.abort(w, r, err)
		return
	}

	user, _ := getUser(r)

	channel, err := s.DB.Channels.GetByURL(decoder.URL)

	isNewChannel := false

	if err != nil {
		if err == sql.ErrNoRows {
			isNewChannel = true
		} else {
			s.abort(w, r, err)
			return
		}
	}

	tx, err := s.DB.Begin()
	if err != nil {
		s.abort(w, r, err)
		return
	}

	if isNewChannel {
		channel = &models.Channel{
			URL: decoder.URL,
		}
		if err := s.Feedparser.FetchChannel(channel); err != nil {
			if err == feedparser.ErrInvalidFeed {
				err = decoders.Errors{
					"url": "Sorry, we were unable to handle this feed, or the feed did not appear to contain any podcasts.",
				}
			}
			s.abort(w, r, err)
			return
		}
	}

	if err := s.DB.Subscriptions.Create(channel.ID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}

	if err := tx.Commit(); err != nil {
		s.abort(w, r, err)
		return
	}

	var status int
	if isNewChannel {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	s.Render.JSON(w, status, channel)
}
