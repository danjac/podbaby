package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
)

func getChannelsByCategory(c *echo.Context) error {

	categoryID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	var (
		channels = []models.Channel{}
		key      = fmt.Sprintf("channels:category:%v", categoryID)
		cache    = getCache(c)
		timeout  = time.Hour
	)

	if err := cache.Get(key, timeout, &channels, func() error {
		store := getStore(c)
		return store.Channels().SelectByCategoryID(store.Conn(), &channels, categoryID)
	}); err != nil {
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
		err = store.Channels().SelectRecommendedByUserID(conn, &channels, user.ID)
	} else {
		err = getCache(c).Get("recommendations", time.Hour*24, &channels, func() error {
			return store.Channels().SelectRecommended(conn, &channels)
		})
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
		page    = getPage(c)
		cache   = getCache(c)
		key     = fmt.Sprintf("channel:%v:page:%v", channelID, page)
		timeout = time.Minute * 30
		detail  = &models.ChannelDetail{}
	)

	if err := cache.Get(key, timeout, detail, func() error {

		var (
			store         = getStore(c)
			conn          = store.Conn()
			channelStore  = store.Channels()
			podcastStore  = store.Podcasts()
			categoryStore = store.Categories()
		)

		detail.Channel = &models.Channel{}

		if err := channelStore.GetByID(conn, detail.Channel, channelID); err != nil {
			return err
		}

		if err := categoryStore.SelectByChannelID(conn, &detail.Categories, channelID); err != nil {
			return err
		}

		if err := channelStore.SelectRelated(conn, &detail.Related, channelID); err != nil {
			return err
		}

		podcasts := &models.PodcastList{}

		if err := podcastStore.SelectByChannel(conn, podcasts, detail.Channel, page); err != nil {
			return err
		}

		for _, pc := range podcasts.Podcasts {
			pc.Name = detail.Channel.Title
			pc.Image = detail.Channel.Image
			pc.ChannelID = channelID
			detail.Podcasts = append(detail.Podcasts, pc)
		}
		detail.Page = podcasts.Page
		return nil
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, detail)
}

func getSubscriptions(c *echo.Context) error {
	var (
		user  = getUser(c)
		store = getStore(c)
	)
	var channels []models.Channel
	if err := store.Channels().SelectSubscribed(store.Conn(), &channels, user.ID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, channels)
}

func addChannel(c *echo.Context) error {

	var (
		v            = newValidator(c)
		store        = getStore(c)
		conn         = store.Conn()
		channelStore = store.Channels()
		user         = getUser(c)
		f            = getFeedparser(c)
	)

	decoder := &newChannelDecoder{}

	if ok, err := v.handle(decoder); !ok {
		return err
	}

	channel := &models.Channel{}
	isNewChannel := false

	if err := channelStore.GetByURL(conn, channel, decoder.URL); err != nil {
		if err == sql.ErrNoRows {
			isNewChannel = true
		} else {
			return err
		}
	}

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if isNewChannel {

		channel = &models.Channel{
			URL: decoder.URL,
		}

		if err := f.Fetch(channel); err != nil {
			if err == feedparser.ErrInvalidFeed {
				return v.invalid(
					"url",
					"Sorry, we were unable to handle this feed, or the feed did not appear to contain any podcasts.",
				).render()
			}
			return err
		}

		if err := channelStore.Create(tx, channel); err != nil {
			return err
		}

		if err := channelStore.AddPodcasts(tx, channel); err != nil {
			return err
		}

		if err := channelStore.AddCategories(tx, channel); err != nil {
			return err
		}

	}

	if err := store.Subscriptions().Create(tx, channel.ID, user.ID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	var status int
	if isNewChannel {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	return c.JSON(status, channel)
}
