package commands

import (
	"fmt"
	"net/http"

	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/server"
)

// Serve runs the webserver
func Serve(cfg *config.Config) {

	log := configureLogger()

	db := database.MustConnect(cfg)
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
