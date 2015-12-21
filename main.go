package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
	"net/http"
)

var (
	env  = flag.String("env", "prod", "environment ('prod' or 'dev')")
	port = flag.String("port", "5000", "server port")
)

const (
	staticURL    = "/static/"
	staticDir    = "./static/"
	devServerURL = "http://localhost:8080"
)

func main() {

	flag.Parse()

	router := mux.NewRouter()

	router.PathPrefix(staticURL).Handler(http.StripPrefix(staticURL,
		http.FileServer(http.Dir(staticDir))))

	render := render.New()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		ctx := map[string]string{
			"staticURL": staticURL,
		}
		render.HTML(w, http.StatusOK, "index", ctx)
	})

	chain := alice.New(nosurf.NewPure).Then(router)

	if err := http.ListenAndServe(":"+*port, chain); err != nil {
		panic(err)
	}

}
