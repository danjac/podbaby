package main

import (
	"encoding/json"
	"flag"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
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

type Signup struct {
	Name     string `db:"name",json:"name",valid:"required"`
	Email    string `db:"email",json:"email",valid:"email,required"`
	Password string `db:"password",json:"password",valid:"required"`
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

	auth.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		// log in here, set cookie, return user

	}).Methods("POST")

	auth.HandleFunc("/signup/", func(w http.ResponseWriter, r *http.Request) {
		// return new user, login

		defer r.Body.Close()

		var signup = &Signup{}
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
			Name:     "auth",
			Value:    strconv.FormatInt(user.ID, 10),
			Expires:  time.Now().Add(time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		render.JSON(w, http.StatusCreated, user)

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
