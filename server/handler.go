package server

import (
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"net/http"
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

	// channels

	channels := api.PathPrefix("/channels/").Subrouter()
	channels.Handle("/", s.requireAuth(s.getChannels)).Methods("GET")
	channels.Handle("/", s.requireAuth(s.addChannel)).Methods("POST")

	// podcasts

	podcasts := api.PathPrefix("/podcasts/").Subrouter()
	podcasts.Handle("/latest/", s.requireAuth(s.getLatestPodcasts)).Methods("GET")

	return nosurf.NewPure(router)
}
