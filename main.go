package main

import (
	"flag"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
	"net/http"
	"time"
)

var (
	env  = flag.String("env", "prod", "environment ('prod' or 'dev')")
	port = flag.String("port", "5000", "server port")
)

// should be settings
const (
	staticURL    = "/static/"
	staticDir    = "./static/"
	devServerURL = "http://localhost:8080"
)

func main() {

	flag.Parse()

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

	router.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims["id"] = "1234567"
		token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, err := token.SignedString([]byte("My Secret"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, tokenString)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(tokenString))
	})

	router.HandleFunc("/secure/", func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user")
		fmt.Fprintf(w, "This is an authenticated request:\n")
		if user == nil {
			fmt.Fprintf(w, "No token set")
		} else {
			fmt.Fprintf(w, "Claim content:\n")
			for k, v := range user.(*jwt.Token).Claims {
				fmt.Fprintf(w, "%s: \t%#v", k, v)
			}
		}
	})

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		CredentialsOptional: true,
		SigningMethod:       jwt.SigningMethodHS256,
	})

	chain := alice.New(nosurf.NewPure, jwtMiddleware.Handler).Then(router)

	if err := http.ListenAndServe(":"+*port, chain); err != nil {
		panic(err)
	}

}
