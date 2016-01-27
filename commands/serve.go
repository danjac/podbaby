package commands

import (
	"github.com/danjac/podbaby/api"
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/store"
	"log"
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

	if err := api.Run(env); err != nil {
		log.Fatalln(err)
	}

}
