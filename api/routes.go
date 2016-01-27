package api

import "github.com/labstack/echo"

func withRoutes(e *echo.Echo) {

	// front page

	e.Get("/", indexPage)

	// API

	api := e.Group("/api/")

	api.Get("search/", searchAll)

	// PUBLIC ROUTES

	// authentication

	auth := api.Group("auth/")

	auth.Post("login/", login)
	auth.Post("signup/", signup)
	auth.Post("recoverpass/", recoverPassword)
	auth.Get("name/", isName)
	auth.Get("email/", isEmail)
	auth.Delete("logout/", logout)

	// channels

	channels := api.Group("channels/")

	channels.Get("category/:id/", getChannelsByCategory)
	channels.Get("recommended/", getRecommendations)
	channels.Get(":id/", getChannelDetail)
	channels.Get(":id/search/", searchChannel)

	// podcasts

	podcasts := api.Group("podcasts/")

	podcasts.Get("detail/:id/", getPodcast)
	podcasts.Get("latest/", getLatestPodcasts)

	// MEMBERS ONLY

	member := api.Group("member/")

	member.Use(authorize())

	member.Post("new/", addChannel)

	// settings

	settings := member.Group("settings/")

	settings.Patch("email/", changeEmail)
	settings.Patch("password/", changePassword)
	settings.Delete("", deleteAccount)

	// subscriptions

	subs := member.Group("subscriptions/")

	subs.Get("", getSubscriptions)
	subs.Get(":prefix.opml", getOPML)
	subs.Post(":id/", subscribe)
	subs.Delete(":id/", unsubscribe)

	// bookmarks

	bookmarks := member.Group("bookmarks/")

	bookmarks.Get("", getBookmarks)
	bookmarks.Get("search/", searchBookmarks)
	bookmarks.Post(":id/", addBookmark)
	bookmarks.Delete(":id/", removeBookmark)

	// plays

	plays := member.Group("plays/")

	plays.Get("", getPlays)
	plays.Post(":id/", addPlay)
	plays.Delete("", deleteAllPlays)

}
