package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) configureRoutes() http.Handler {

	router := mux.NewRouter()

	// static routes

	router.PathPrefix(s.Config.StaticURL).Handler(
		http.StripPrefix(s.Config.StaticURL,
			http.FileServer(http.Dir(s.Config.StaticDir))))

	// front page

	router.Handle("/", s.authOptionalHandler(indexPage))

	// OPML download

	router.Handle("/{prefix}.opml", s.authRequiredHandler(getOPML)).Methods("GET")

	// API

	api := router.PathPrefix("/api/").Subrouter()

	// authentication

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.Handle("/login/", s.authIgnoreHandler(login)).Methods("POST")
	auth.Handle("/signup/", s.authIgnoreHandler(signup)).Methods("POST")
	auth.Handle("/logout/", s.authIgnoreHandler(logout)).Methods("DELETE")
	auth.Handle("/recoverpass/", s.authIgnoreHandler(recoverPassword)).Methods("POST")

	// user

	user := api.PathPrefix("/user/").Subrouter()

	user.Handle("/name/", s.authIgnoreHandler(isName)).Methods("GET")
	user.Handle("/email/", s.authOptionalHandler(isEmail)).Methods("GET")

	user.Handle("/email/", s.authRequiredHandler(changeEmail)).Methods("PATCH")
	user.Handle("/password/", s.authRequiredHandler(changePassword)).Methods("PATCH")
	user.Handle("/", s.authRequiredHandler(deleteAccount)).Methods("DELETE")
	// user.Handle("/", handler{s, authRequired, deleteAccount}).Methods("DELETE")

	// channels

	channels := api.PathPrefix("/channels/").Subrouter()

	channels.Handle("/{id:[0-9]+}/", s.authIgnoreHandler(getChannelDetail)).Methods("GET")
	channels.Handle("/", s.authRequiredHandler(getChannels)).Methods("GET")
	channels.Handle("/", s.authRequiredHandler(addChannel)).Methods("POST")

	// search

	search := api.PathPrefix("/search/").Subrouter()

	search.Handle("/", s.authIgnoreHandler(searchAll)).Methods("GET")
	search.Handle("/channel/{id:[0-9]+}/", s.authIgnoreHandler(searchChannel)).Methods("GET")
	search.Handle("/bookmarks/", s.authRequiredHandler(searchBookmarks)).Methods("GET")

	// subscriptions

	subs := api.PathPrefix("/subscriptions/").Subrouter()

	subs.Handle("/{id:[0-9]+}/", s.authRequiredHandler(subscribe)).Methods("POST")
	subs.Handle("/{id:[0-9]+}/", s.authRequiredHandler(unsubscribe)).Methods("DELETE")

	// bookmarks

	bookmarks := api.PathPrefix("/bookmarks/").Subrouter()

	bookmarks.Handle("/", s.authRequiredHandler(getBookmarks)).Methods("GET")
	bookmarks.Handle("/{id:[0-9]+}/", s.authRequiredHandler(addBookmark)).Methods("POST")
	bookmarks.Handle("/{id:[0-9]+}/", s.authRequiredHandler(removeBookmark)).Methods("DELETE")

	// plays

	plays := api.PathPrefix("/plays/").Subrouter()

	plays.Handle("/", s.authRequiredHandler(getPlays)).Methods("GET")
	plays.Handle("/", s.authRequiredHandler(deleteAllPlays)).Methods("DELETE")
	plays.Handle("/{id:[0-9]+}/", s.authRequiredHandler(addPlay)).Methods("POST")

	// podcasts

	podcasts := api.PathPrefix("/podcasts/").Subrouter()

	podcasts.Handle("/detail/{id:[0-9]+}/", s.authIgnoreHandler(getPodcast)).Methods("GET")
	podcasts.Handle("/latest/", s.authOptionalHandler(getLatestPodcasts)).Methods("GET")

	return router
}
