package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/server"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	_ "github.com/lib/pq"
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

	db := database.New(sqlx.MustConnect("postgres", *url))

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	log.Info("Starting web service...")

	s := server.New(db, log, &server.Config{
		StaticURL: staticURL,
		StaticDir: staticDir,
		SecretKey: "my-secret",
	})

	chain := alice.New(nosurf.NewPure).Then(s.Router())

	if err := http.ListenAndServe(":"+*port, chain); err != nil {
		panic(err)
	}

}
