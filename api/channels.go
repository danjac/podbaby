package api

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

	channelID, err := getInt64(c, "id") // getIntOr404(c, "id") ???
	if err != nil {
		return err
	}

	var (
		store = storeFromContext(c)
		conn  = store.Conn()
	)

	channel, err := store.Channels().GetByID(conn, channelID)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		return err
	}

	detail := &models.ChannelDetail{
		Channel: channel,
	}

	categories, err := store.Categories().SelectByChannelID(conn, channelID)
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

func getSubscriptions(c *echo.Context) error {
	var (
		user  = userFromContext(c)
		store = storeFromContext(c)
	)
	channels, err := store.Channels().SelectSubscribed(store.Conn(), user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, channels)
}

func addChannel(c *echo.Context) error {

	decoder := &decoders.NewChannel{}
	if err, ok := decoders.Decode(c, decoder); !ok {
		return err
	}

	var (
		store        = storeFromContext(c)
		conn         = store.Conn()
		channelStore = store.Channels()
		user         = userFromContext(c)
	)
	channel, err := channelStore.GetByURL(conn, decoder.URL)
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
				errors := decoders.Errors{
					"url": "Sorry, we were unable to handle this feed, or the feed did not appear to contain any podcasts.",
				}
				return errors.Render(c)
			}
			return err
		}
		tx, err := conn.Begin()
		if err != nil {
			return err
		}

		defer func() {
			_ = tx.Rollback()
		}()

		if err := channelStore.Create(tx, channel); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}

	}

	if err := store.Subscriptions().Create(conn, channel.ID, user.ID); err != nil {
		return err
	}

	if isNewChannel {
		go func(c *echo.Context, channel *models.Channel) {

			var (
				store        = storeFromContext(c)
				channelStore = store.Channels()
				log          = c.Echo().Logger()
			)
			tx, err := store.Conn().Begin()
			if err != nil {
				log.Error(err)
				return
			}

			defer func() {
				_ = tx.Rollback()
			}()

			if err := channelStore.AddPodcasts(tx, channel); err != nil {
				log.Error(err)
				return
			}

			if err := channelStore.AddCategories(tx, channel); err != nil {
				log.Error(err)
				return
			}

			if err := tx.Commit(); err != nil {
				log.Error(err)
			}

		}(c, channel)

	}

	var status int
	if isNewChannel {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	return c.JSON(status, channel)
}
