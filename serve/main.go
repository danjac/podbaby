package main

import (
	"flag"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/api"
	"github.com/danjac/podbaby/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	env  = flag.String("env", "prod", "environment ('prod' or 'dev')")
	port = flag.String("port", "5000", "server port")
	url  = flag.String("url", "", "database connection url")
)

// should be settings
const (
	defaultStaticURL = "/static/"
	defaultStaticDir = "./static/"
	devStaticURL     = "http://localhost:8080/static/"
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

	var staticURL string
	if *env == "dev" {
		staticURL = devStaticURL
	} else {
		staticURL = defaultStaticURL
	}

	api := api.New(db, log, &api.Config{
		StaticURL: staticURL,
		StaticDir: defaultStaticDir,
		SecretKey: "my-secret",
	})

	if err := http.ListenAndServe(":"+*port, api.Handler()); err != nil {
		panic(err)
	}

}
