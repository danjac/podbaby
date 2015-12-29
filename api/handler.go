package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
)

func (s *Server) Handler() http.Handler {

	router := mux.NewRouter()
	// static routes

	router.PathPrefix(s.Config.StaticURL).Handler(
		http.StripPrefix(s.Config.StaticURL,
			http.FileServer(http.Dir(s.Config.StaticDir))))

	// front page
	router.HandleFunc("/", s.indexPage)

	// API

	api := router.PathPrefix("/api/").Subrouter()

	// authentication

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/login/", s.login).Methods("POST")
	auth.HandleFunc("/signup/", s.signup).Methods("POST")
	auth.HandleFunc("/logout/", s.logout).Methods("DELETE")

	// user

	user := api.PathPrefix("/user/").Subrouter()

	user.Handle("/email/", s.requireAuth(s.changeEmail)).Methods("PATCH")

	// channels

	channels := api.PathPrefix("/channels/").Subrouter()

	channels.Handle("/", s.requireAuth(s.getChannels)).Methods("GET")
	channels.Handle("/", s.requireAuth(s.addChannel)).Methods("POST")
	channels.Handle("/{id}/", s.requireAuth(s.getChannelDetail)).Methods("GET")

	// search

	search := api.PathPrefix("/search/").Subrouter()

	search.Handle("/", s.requireAuth(s.search)).Methods("GET")

	// subscriptions

	subs := api.PathPrefix("/subscriptions/").Subrouter()

	subs.Handle("/{id}/", s.requireAuth(s.subscribe)).Methods("POST")
	subs.Handle("/{id}/", s.requireAuth(s.unsubscribe)).Methods("DELETE")

	// bookmarks

	bookmarks := api.PathPrefix("/bookmarks/").Subrouter()

	bookmarks.Handle("/", s.requireAuth(s.getBookmarks)).Methods("GET")
	bookmarks.Handle("/{id}/", s.requireAuth(s.addBookmark)).Methods("POST")
	bookmarks.Handle("/{id}/", s.requireAuth(s.removeBookmark)).Methods("DELETE")

	// podcasts

	podcasts := api.PathPrefix("/podcasts/").Subrouter()

	podcasts.Handle("/latest/", s.requireAuth(s.getLatestPodcasts)).Methods("GET")

	return nosurf.NewPure(router)
}
