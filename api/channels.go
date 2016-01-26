package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
)

func getChannelsByCategory(c *echo.Context) error {

	categoryID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	store := getStore(c)

	channels, err := store.Channels().SelectByCategoryID(store.Conn(), categoryID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, channels)
}

func getRecommendations(c *echo.Context) error {

	var (
		channels []models.Channel
		err      error
		store    = getStore(c)
		conn     = store.Conn()
	)

	user, _ := authenticate(c)
	if user != nil {
		channels, err = store.Channels().SelectRecommendedByUserID(conn, user.ID)
	} else {
		channels, err = store.Channels().SelectRecommended(conn)
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, channels)
}

func getChannelDetail(c *echo.Context) error {

	channelID, err := getIntOr404(c, "id")

	if err != nil {
		return err
	}

	var (
		store         = getStore(c)
		conn          = store.Conn()
		channelStore  = store.Channels()
		podcastStore  = store.Podcasts()
		categoryStore = store.Categories()
	)

	channel, err := channelStore.GetByID(conn, channelID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return err
	}

	detail := &models.ChannelDetail{
		Channel: channel,
	}

	categories, err := categoryStore.SelectByChannelID(conn, channelID)
	if err != nil {
		return err
	}

	detail.Categories = categories

	related, err := channelStore.SelectRelated(conn, channelID)
	if err != nil {
		return err
	}

	detail.Related = related

	podcasts, err := podcastStore.SelectByChannelID(conn, channelID, getPage(c))
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
	return c.JSON(http.StatusOK, detail)
}

func getSubscriptions(c *echo.Context) error {
	var (
		user  = getUser(c)
		store = getStore(c)
	)
	channels, err := store.Channels().SelectSubscribed(store.Conn(), user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, channels)
}

func addChannel(c *echo.Context) error {

	var (
		validator    = newValidator(c)
		store        = getStore(c)
		conn         = store.Conn()
		channelStore = store.Channels()
		user         = getUser(c)
		f            = getFeedparser(c)
	)

	log.Println("OK")

	decoder := &newChannelDecoder{}

	if ok, err := validator.validate(decoder); !ok {
		log.Println("BAD ERROR", err)
		return err
	}

	channel, err := channelStore.GetByURL(conn, decoder.URL)
	isNewChannel := false

	if err != nil {
		if err == sql.ErrNoRows {
			isNewChannel = true
		} else {
			log.Println("BAD SQL ERROR:", err)
			return err
		}
	}
	log.Println("NEWCHANNEL?", isNewChannel)

	if isNewChannel {

		channel = &models.Channel{
			URL: decoder.URL,
		}

		log.Println("OK2")

		if err := f.Fetch(channel); err != nil {
			if err == feedparser.ErrInvalidFeed {
				return validator.invalid(
					"url",
					"Sorry, we were unable to handle this feed, or the feed did not appear to contain any podcasts.",
				).render()
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
				store        = getStore(c)
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
