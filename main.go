package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/asaskevich/govalidator"
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

type User struct {
	ID       int64  `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type Channel struct {
	ID          int64  `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image"`
	URL         string `db:"url" json:"url"`
}

type Podcast struct {
	ChannelID    int64     `db:"channel_id"`
	Title        string    `db:"title" json:"title"`
	Description  string    `db:"description" json:"description"`
	EnclosureURL string    `db:"enclosure_url" json:"enclosureUrl"`
	PubDate      time.Time `db:"pub_date" json:"pubDate"`
}

type Signup struct {
	Name     string `json:"name",valid:"required"`
	Email    string `json:"email",valid:"email,required"`
	Password string `json:"password",valid:"required"`
}

type Login struct {
	Identifier string `json:"identifier",valid:"required"`
	Password   string `json:"password",valid:"required"`
}

type NewChannel struct {
	URL string `json:"url",valid:"url,required"`
}

func fetchPodcasts(db *sqlx.DB, url string) error {

	channel := &Channel{
		URL: url,
	}

	query, args, err := sqlx.Named("INSERT INTO channels (url)  VALUES (:url) RETURNING id", channel)
	if err != nil {
		return err
	}

	if err := db.QueryRow(db.Rebind(query), args...).Scan(&channel.ID); err != nil {
		return err
	}

	isChannelUpdated := false
	var rssErr error
	chanHandler := func(feed *rss.Feed, newChannels []*rss.Channel) {
		for _, ch := range newChannels {
			if isChannelUpdated {
				return
			}

			channel.Title = ch.Title
			channel.Image = ch.Image.Url
			channel.Description = ch.Description

			query, args, err := sqlx.Named("UPDATE channels SET title=:title, image=:image, description=:description WHERE id=:id", channel)
			if err != nil {
				rssErr = err
				return
			}
			fmt.Println(query, args)
			if _, err := db.Exec(db.Rebind(query), args...); err != nil {
				rssErr = err
				return
			}
			isChannelUpdated = true
		}
	}

	itemHandler := func(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
		fmt.Println("Items for Channel:", ch.Title)
		for _, item := range newItems {
			fmt.Println("Item:", item.Title)
			pubDate, _ := item.ParsedPubDate()
			pc := &Podcast{
				ChannelID:   channel.ID,
				Title:       item.Title,
				Description: item.Description,
				PubDate:     pubDate,
			}
			for _, enclosure := range item.Enclosures {
				pc.EnclosureURL = enclosure.Url
				break
			}
			if pc.EnclosureURL == "" {
				continue
			}
			query, args, err := sqlx.Named("INSERT INTO podcasts (channel_id, title, description, enclosure_url, pub_date) VALUES(:channel_id, :title, :description, :enclosure_url, :pub_date)", pc)
			if err != nil {
				rssErr = err
				return
			}
			fmt.Println(query, args)
			if _, err := db.Exec(db.Rebind(query), args...); err != nil {
				rssErr = err
				return
			}

		}
	}

	feed := rss.New(5, true, chanHandler, itemHandler)
	if err := feed.Fetch(url, nil); err != nil {
		return err
	}
	return rssErr
}

func main() {

	flag.Parse()

	db := sqlx.MustConnect("postgres", *url)
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
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user := &User{}
		if err := db.Get(user, "SELECT id, name, email FROM users WHERE id=$1", cookie.Value); err != nil {
			// check for no rows
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, http.StatusOK, user)

	}).Methods("GET")

	auth.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		login := &Login{}

		if err := json.NewDecoder(r.Body).Decode(login); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := govalidator.ValidateStruct(login); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// find the user
		user := &User{}

		if err := db.Get(user, "SELECT id, name FROM users WHERE email=$1 or name=$1", login.Identifier); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No user found", http.StatusBadRequest)
				return
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
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

		defer r.Body.Close()

		signup := &Signup{}
		if err := json.NewDecoder(r.Body).Decode(signup); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := govalidator.ValidateStruct(signup); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if email exists
		var num int64

		if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", signup.Email).Scan(&num); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if num != 0 {
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}

		// make new user

		password := []byte(signup.Password)

		hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		encryptedPassword := string(hash)

		user := &User{
			Name:     signup.Name,
			Email:    signup.Email,
			Password: encryptedPassword,
		}

		sql := "INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id"

		if err := db.QueryRow(sql, user.Name, user.Email, user.Password).Scan(&user.ID); err != nil {
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
		/*
			s := &NewChannel{}
			defer r.Body.Close()

			if err := json.NewDecoder(r.Body).Decode(s); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Println("ADDING NEW CHANNEL....", s.URL)

			if _, err := govalidator.ValidateStruct(s); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		*/

		channel := &Channel{
			URL: r.FormValue("url"),
		}

		url := r.FormValue("url")

		go func(db *sqlx.DB, url string) {
			if err := fetchPodcasts(db, url); err != nil {
				fmt.Println("FEEDERROR", err)
			}
		}(db, url)

		render.JSON(w, http.StatusOK, channel)

	}).Methods("POST")

	pc.HandleFunc("/channels/{id}/", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	//chain := alice.New(nosurf.NewPure).Then(router)
	chain := router

	if err := http.ListenAndServe(":"+*port, chain); err != nil {
		panic(err)
	}

}
