package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	rss "github.com/jteeuwen/go-pkg-rss"
	//"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
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

func fetchPodcasts(db *database.DB, url string) error {

	var rssChannel *rss.Channel

	chanHandler := func(feed *rss.Feed, newChannels []*rss.Channel) {
		rssChannel = newChannels[0]
	}

	var rssItems []*rss.Item

	itemHandler := func(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
		rssItems = append(rssItems, newItems...)
	}

	feed := rss.New(5, true, chanHandler, itemHandler)

	if err := feed.Fetch(url, nil); err != nil {
		return err
	}

	// tbd: check if channel already exists
	channel := &models.Channel{
		URL:         url,
		Title:       rssChannel.Title,
		Image:       rssChannel.Image.Url,
		Description: rssChannel.Description,
	}

	if err := db.Channels.Create(channel); err != nil {
		return err
	}

	// tbd: check pubdates: only insert if pub_date > MAX pub date of existing
	// items: make enclosure URL + channel ID unique

	for _, item := range rssItems {
		pubDate, _ := item.ParsedPubDate()

		pc := &models.Podcast{
			ChannelID:   channel.ID,
			Title:       item.Title,
			Description: item.Description,
			PubDate:     pubDate,
		}

		if len(item.Enclosures) == 0 {
			continue
		}

		pc.EnclosureURL = item.Enclosures[0].Url
		if err := db.Podcasts.Create(pc); err != nil {
			return err
		}

	}

	return nil
}

func main() {

	flag.Parse()

	db := database.New(sqlx.MustConnect("postgres", *url))
	router := mux.NewRouter()

	router.PathPrefix(staticURL).Handler(
		http.StripPrefix(staticURL, http.FileServer(http.Dir(staticDir))))

	render := render.New()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		csrfToken := nosurf.Token(r)
		ctx := map[string]string{
			"staticURL": staticURL,
			"csrfToken": csrfToken,
		}
		render.HTML(w, http.StatusOK, "index", ctx)
	})

	api := router.PathPrefix("/api/").Subrouter()
	auth := api.PathPrefix("/auth/").Subrouter()

	// load user on startup
	auth.HandleFunc("/currentuser/", func(w http.ResponseWriter, r *http.Request) {
		// log in here, set cookie, return username
		cookie, err := r.Cookie("userID")
		if err != nil {
			// check if cookie not found
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if cookie.Value == "" {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		user, err := db.Users.GetByID(cookie.Value)
		if err != nil {
			// check for no rows
			if err == sql.ErrNoRows {
				http.Error(w, "No user found", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, http.StatusOK, user)

	}).Methods("GET")

	auth.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {

		decoder := &decoders.Login{}

		if err := decoder.Decode(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := db.Users.GetByNameOrEmail(decoder.Identifier)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No user found", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(decoder.Password)); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid password", http.StatusBadRequest)
				return
			}
		}
		// login user

		cookie := &http.Cookie{
			Name:    "userID",
			Value:   strconv.FormatInt(user.ID, 10),
			Expires: time.Now().Add(time.Hour),
			//Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		// tbd: no need to return user!
		render.JSON(w, http.StatusOK, user)

	}).Methods("POST")

	auth.HandleFunc("/signup/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login

		decoder := &decoders.Signup{}

		if err := decoder.Decode(r); r != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if exists, _ := db.Users.IsEmail(decoder.Email); exists {
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}

		if exists, _ := db.Users.IsName(decoder.Name); exists {
			http.Error(w, "Name already taken", http.StatusBadRequest)
			return
		}

		// make new user

		password := []byte(decoder.Password)

		hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		encryptedPassword := string(hash)

		user := &models.User{
			Name:     decoder.Name,
			Email:    decoder.Email,
			Password: encryptedPassword,
		}

		if err := db.Users.Create(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// sign user in
		cookie := &http.Cookie{
			Name:    "userID",
			Value:   strconv.FormatInt(user.ID, 10),
			Expires: time.Now().Add(time.Hour),
			//Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		// tbd: no need to return user!
		render.JSON(w, http.StatusCreated, user)

	}).Methods("POST")

	auth.HandleFunc("/recoverpass/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login
	}).Methods("POST")

	auth.HandleFunc("/changepass/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login
	}).Methods("PATCH")

	auth.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login
	}).Methods("PUT")

	auth.HandleFunc("/logout/", func(w http.ResponseWriter, r *http.Request) {
		// delete cookie
		cookie := &http.Cookie{
			Name:    "userID",
			Value:   "",
			Expires: time.Now().Add(time.Hour),
			//Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}).Methods("DELETE")

	pc := api.PathPrefix("/podcasts/").Subrouter()

	pc.HandleFunc("/latest/", func(w http.ResponseWriter, r *http.Request) {
		podcasts, err := db.Podcasts.SelectAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, http.StatusOK, podcasts)

	}).Methods("GET")

	pc.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	pc.HandleFunc("/subscriptions/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	pc.HandleFunc("/subscriptions/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("POST")

	pc.HandleFunc("/subscriptions/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("DELETE")

	pc.HandleFunc("/bookmarks/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	pc.HandleFunc("/bookmarks/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("POST")

	pc.HandleFunc("/bookmarks/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("DELETE")

	pc.HandleFunc("/channels/", func(w http.ResponseWriter, r *http.Request) {

		// add channel
		decoder := &decoders.NewChannel{}

		if err := decoder.Decode(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		go func(db *database.DB, url string) {
			if err := fetchPodcasts(db, url); err != nil {
				fmt.Println("FEEDERROR", err)
			}
		}(db, decoder.URL)

		render.Text(w, http.StatusCreated, "New channel added")

	}).Methods("POST")

	pc.HandleFunc("/channels/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	//chain := alice.New(nosurf.NewPure).Then(router)
	chain := router

	if err := http.ListenAndServe(":"+*port, chain); err != nil {
		panic(err)
	}

}
