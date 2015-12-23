package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"net/http"
)

var (
	env  = flag.String("env", "prod", "environment ('prod' or 'dev')")
	port = flag.String("port", "5000", "server port")
	url  = flag.String("url", "", "database connection url")
)

// should be settings
const (
	staticURL    = "/static/"
	staticDir    = "./static/"
	devServerURL = "http://localhost:8080"
)

func main() {

	flag.Parse()

	sqlx.MustConnect("postgres", *url)
	router := mux.NewRouter()

	router.PathPrefix(staticURL).Handler(
		http.StripPrefix(staticURL, http.FileServer(http.Dir(staticDir))))

	render := render.New()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]string{
			"staticURL": staticURL,
		}
		render.HTML(w, http.StatusOK, "index", ctx)
	})

	api := router.PathPrefix("/api/").Subrouter()
	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		// log in here, set cookie, return user

	}).Methods("POST")

	auth.HandleFunc("/signup/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login

	}).Methods("POST")

	auth.HandleFunc("/recoverpass/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login
	}).Methods("POST")

	auth.HandleFunc("/changepass/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login
	}).Methods("PATCH")

	auth.HandleFunc("/changeid/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login
	}).Methods("PATCH")

	auth.HandleFunc("/logout/", func(w http.ResponseWriter, r *http.Request) {
		// delete cookie
	}).Methods("POST")

	pc := api.PathPrefix("/podcasts/").Subrouter()

	pc.HandleFunc("/latest/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	pc.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	pc.HandleFunc("/subscriptions/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	pc.HandleFunc("/subscriptions/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("POST")

	pc.HandleFunc("/subscriptions/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("DELETE")

	pc.HandleFunc("/pins/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	pc.HandleFunc("/pins/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("POST")

	pc.HandleFunc("/unpin/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("DELETE")

	pc.HandleFunc("/channels/", func(w http.ResponseWriter, r *http.Request) {
		// add channel
	}).Methods("POST")

	pc.HandleFunc("/channels/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	chain := alice.New(nosurf.NewPure).Then(router)

	if err := http.ListenAndServe(":"+*port, chain); err != nil {
		panic(err)
	}

}
