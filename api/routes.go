package api

import (
	"github.com/danjac/podbaby/config"
	"github.com/labstack/echo"
	"net/http"
)

func withRoutes(e *echo.Echo, cfg *config.Config) {

	// front page

	c.Get("/", indexPage)

	// API

	api := e.Group("/api/")

	// PUBLIC ROUTES

	// authentication

	auth := api.Group("/auth/")
	auth.Post("/login/", login)
	auth.Post("/signup/", signup)
	auth.Post("/recoverpass/", recoverpass)
	auth.Get("/name/", isName)
	auth.Get("/email/", isEmail)
	auth.Delete("/logout/", logout)

	// channels

	channels := api.Group("/channels/")
	channels.Get("/:id/", getChannelDetail)
	channels.Get("/category/:id/", getChannelsByCategory)

	search := api.Group("/search/")
	search.Get("/", searchAll)
	search.Get("/channels/:id/", searchChannel)

	// podcasts

	podcasts := api.Group("/podcasts/")
	podcasts.Get("/detail/:id/", getPodcast)
	podcasts.Get("/latest/", getLatestPodcasts)

	// MEMBERS ONLY

	member := api.Group("/member/")
	member.Use(authorize)

	member.Get("/recommended/", getRecommendations)
	member.Post("/new/", addChannel)

	// settings

	settings := member.Group("/settings/")
	settings.Use(authMiddleware(true))

	settings.Get("/email/", isEmail)
	settings.Patch("/email/", changeEmail)
	settings.Delete("/", deleteAccount)

	// subscriptions

	subs = member.Group("/subscriptions/")
	subs.Get("/", getChannels)
	subs.Get("/:prefix.opml", getOPML)
	subs.Post("/:id/", subscribe)
	subs.Delete("/:id/", unsubscribe)

	// bookmarks

	bookmarks := member.Group("/bookmarks/")
	bookmarks.Get("/", getBookmarks)
	bookmarks.Post("/:id/", addBookmark)
	bookmarks.Delete("/:id/", removeBookmark)

	// plays

	plays := member.Group("/plays/")
	plays.Get("/", getPlays)
	plays.Post("/", addPlay)
	plays.Delete("/", deleteAllPlays)

}
