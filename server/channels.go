package server

import (
	"net/http"

	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
)

func getChannelsByCategory(s *Server, w http.ResponseWriter, r *http.Request) error {

	categoryID, _ := getID(r)
	channels, err := s.DB.Channels.SelectByCategoryID(categoryID)
	if err != nil {
		return err
	}
	return s.Render.JSON(w, http.StatusOK, channels)
}

func getRecommendations(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, ok := getUser(r)
	var (
		channels []models.Channel
		err      error
	)
	if ok {
		channels, err = s.DB.Channels.SelectRecommendedByUserID(user.ID)
	} else {
		channels, err = s.DB.Channels.SelectRecommended()
	}

	if err != nil {
		return err
	}

	return s.Render.JSON(w, http.StatusOK, channels)
}

func getChannelDetail(s *Server, w http.ResponseWriter, r *http.Request) error {
	channelID, _ := getID(r)

	channel, err := s.DB.Channels.GetByID(channelID)
	if err != nil {
		return err
	}

	detail := &models.ChannelDetail{
		Channel: channel,
	}

	categories, err := s.DB.Categories.SelectByChannelID(channelID)
	if err != nil {
		return err
	}

	detail.Categories = categories

	related, err := s.DB.Channels.SelectRelated(channelID)
	if err != nil {
		return err
	}

	detail.Related = related

	podcasts, err := s.DB.Podcasts.SelectByChannelID(channelID, getPage(r))
	if err != nil {
		return err
	}
	for _, pc := range podcasts.Podcasts {
		pc.Name = channel.Title
		pc.Image = channel.Image
		pc.ChannelID = channel.ID
		detail.Podcasts = append(detail.Podcasts, pc)
	}
	detail.Page = podcasts.Page
	return s.Render.JSON(w, http.StatusOK, detail)
}

func getChannels(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	channels, err := s.DB.Channels.SelectSubscribed(user.ID)
	if err != nil {
		return err
	}
	return s.Render.JSON(w, http.StatusOK, channels)
}

func addChannel(s *Server, w http.ResponseWriter, r *http.Request) error {

	decoder := &decoders.NewChannel{}

	if err := decoders.Decode(r, decoder); err != nil {
		return err
	}

	user, _ := getUser(r)
	/*

	   conn := store.FromContext(ctx)
	   channelsDB := conn.Channels()

	   channels, err := channelsDB.GetByURL(conn, decoder.URL)

	   tx, err := conn.Begin()

	   err := channelsDB.Create(tx)

	*/
	channel, err := s.DB.Channels.GetByURL(decoder.URL)
	isNewChannel := false

	if err != nil {
		if isErrNoRows(err) {
			isNewChannel = true
		} else {
			return err
		}
	}

	if isNewChannel {

		channel = &models.Channel{
			URL: decoder.URL,
		}

		if err := s.Feedparser.Fetch(channel); err != nil {
			if err == feedparser.ErrInvalidFeed {
				err = decoders.Errors{
					"url": "Sorry, we were unable to handle this feed, or the feed did not appear to contain any podcasts.",
				}
			}
			return err
		}
		tx, err := s.DB.Begin()
		if err != nil {
			return err
		}

		defer func() {
			_ = tx.Rollback()
		}()

		if err := tx.Create(channel); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}

	}

	if err := s.DB.Subscriptions.Create(channel.ID, user.ID); err != nil {
		return err
	}

	if isNewChannel {
		go func(channel *models.Channel) {

			tx, err := s.DB.Channels.Begin()
			if err != nil {
				s.Log.Error(err)
				return
			}

			defer func() {
				_ = tx.Rollback()
			}()

			if err := tx.AddPodcasts(channel); err != nil {
				s.Log.Error(err)
				return
			}

			if err := tx.AddCategories(channel); err != nil {
				s.Log.Error(err)
				return
			}

			if err := tx.Commit(); err != nil {
				s.Log.Error(err)
			}

		}(channel)

	}

	var status int
	if isNewChannel {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	return s.Render.JSON(w, status, channel)
}
