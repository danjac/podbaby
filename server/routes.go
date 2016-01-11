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

	router.HandleFunc("/", s.indexPage)

	// OPML download

	router.Handle("/{prefix}.opml", s.requireAuth(s.getOPML)).Methods("GET")

	// API

	api := router.PathPrefix("/api/").Subrouter()

	// authentication

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/login/", s.login).Methods("POST")
	auth.HandleFunc("/signup/", s.signup).Methods("POST")
	auth.HandleFunc("/logout/", s.logout).Methods("DELETE")
	auth.HandleFunc("/recoverpass/", s.recoverPassword).Methods("POST")

	// user

	user := api.PathPrefix("/user/").Subrouter()

	user.HandleFunc("/name/", s.isName).Methods("GET")
	user.HandleFunc("/email/", s.isEmail).Methods("GET")

	user.Handle("/email/", s.requireAuth(s.changeEmail)).Methods("PATCH")
	user.Handle("/password/", s.requireAuth(s.changePassword)).Methods("PATCH")
	user.Handle("/", s.requireAuth(s.deleteAccount)).Methods("DELETE")

	// channels

	channels := api.PathPrefix("/channels/").Subrouter()

	channels.HandleFunc("/{id:[0-9]+}/", s.getChannelDetail).Methods("GET")
	channels.Handle("/", s.requireAuth(s.getChannels)).Methods("GET")
	channels.Handle("/", s.requireAuth(s.addChannel)).Methods("POST")

	// search

	search := api.PathPrefix("/search/").Subrouter()

	search.HandleFunc("/", s.search).Methods("GET")
	search.HandleFunc("/channel/{id:[0-9]+}/", s.searchChannel).Methods("GET")
	search.Handle("/bookmarks/", s.requireAuth(s.searchBookmarks)).Methods("GET")

	// subscriptions

	subs := api.PathPrefix("/subscriptions/").Subrouter()

	subs.Handle("/{id:[0-9]+}/", s.requireAuth(s.subscribe)).Methods("POST")
	subs.Handle("/{id:[0-9]+}/", s.requireAuth(s.unsubscribe)).Methods("DELETE")

	// bookmarks

	bookmarks := api.PathPrefix("/bookmarks/").Subrouter()

	bookmarks.Handle("/", s.requireAuth(s.getBookmarks)).Methods("GET")
	bookmarks.Handle("/{id:[0-9]+}/", s.requireAuth(s.addBookmark)).Methods("POST")
	bookmarks.Handle("/{id:[0-9]+}/", s.requireAuth(s.removeBookmark)).Methods("DELETE")

	// plays

	plays := api.PathPrefix("/plays/").Subrouter()

	plays.Handle("/", s.requireAuth(s.getPlays)).Methods("GET")
	plays.Handle("/", s.requireAuth(s.deleteAllPlays)).Methods("DELETE")
	plays.Handle("/{id:[0-9]+}/", s.requireAuth(s.addPlay)).Methods("POST")

	// podcasts

	podcasts := api.PathPrefix("/podcasts/").Subrouter()

	podcasts.HandleFunc("/detail/{id:[0-9]+}/", s.getPodcast).Methods("GET")
	podcasts.HandleFunc("/latest/", s.getLatestPodcasts).Methods("GET")

	return router
}
