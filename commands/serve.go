package commands

import (
	"fmt"
	"github.com/danjac/podbaby/api"
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/store"
	"log"
	"net/http"
)

// Serve runs the webserver
func Serve(cfg *config.Config) {

	store, err := store.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	defer store.Conn().Close()

	mailer, err := mailer.New(cfg)

	feedparser := feedparser.New()

	env := &api.Env{
		Store:      store,
		Mailer:     mailer,
		Config:     cfg,
		Feedparser: feedparser,
	}

	handler, err := api.New(env)
	if err != nil {
		log.Fatalln(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler); err != nil {
		panic(err)
	}

}
