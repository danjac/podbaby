package commands

import (
	"fmt"
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/store"
	"golang.org/x/net/context"
	"net/http"
)

// Serve runs the webserver
func Serve(cfg *config.Config) {

	log := configureLogger()

	db, err := store.New(cfg)
	defer db.Close()

	mailer, err := mailer.New(cfg)

	if err != nil {
		panic(err)
	}

	handler := server.New(db, mailer, log, cfg).Configure()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler); err != nil {
		panic(err)
	}

}
