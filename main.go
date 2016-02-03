package main

import (
	"os"

	"github.com/danjac/podbaby/Godeps/_workspace/src/github.com/codegangsta/cli"
	_ "github.com/danjac/podbaby/Godeps/_workspace/src/github.com/joho/godotenv/autoload"
	_ "github.com/danjac/podbaby/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/danjac/podbaby/commands"
	"github.com/danjac/podbaby/config"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true

	cfg := config.New()

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Run the server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "url",
					EnvVar:      "DB_URL",
					Usage:       "Database connection URL",
					Destination: &cfg.DatabaseURL,
				},
				cli.StringFlag{
					Name:        "secret",
					EnvVar:      "SECRET_KEY",
					Usage:       "Secret key",
					Destination: &cfg.SecretKey,
				},
				cli.StringFlag{
					Name:        "cookie-secret",
					EnvVar:      "SECURE_COOKIE_KEY",
					Usage:       "32-bit secure key",
					Destination: &cfg.SecureCookieKey,
				},
				cli.StringFlag{
					Name:        "google-analytics-id",
					EnvVar:      "GOOGLE_ANALYTICS_ID",
					Usage:       "Your Google Analytics ID (optional)",
					Destination: &cfg.GoogleAnalyticsID,
				},
				cli.IntFlag{
					Name:        "port",
					Value:       5000,
					EnvVar:      "PORT",
					Usage:       "Server port",
					Destination: &cfg.Port,
				},
				cli.IntFlag{
					Name:        "max-db-conns",
					Value:       99,
					EnvVar:      "MAx_DB_CONNECTIONS",
					Usage:       "Maximum database connections",
					Destination: &cfg.MaxDBConnections,
				},
				cli.StringFlag{
					Name:        "env",
					Value:       "prod",
					Usage:       "Environment",
					Destination: &cfg.Env,
				},
				cli.StringFlag{
					Name:        "mail-addr",
					EnvVar:      "MAIL_ADDR",
					Value:       "",
					Usage:       "Email address",
					Destination: &cfg.Mail.Addr,
				},
				cli.StringFlag{
					Name:        "mail-host",
					EnvVar:      "MAIL_HOST",
					Value:       "",
					Usage:       "Email host",
					Destination: &cfg.Mail.Host,
				},
				cli.StringFlag{
					Name:        "mail-user",
					EnvVar:      "MAIL_USER",
					Value:       "",
					Usage:       "Email user",
					Destination: &cfg.Mail.User,
				},
				cli.StringFlag{
					Name:        "mail-password",
					EnvVar:      "MAIL_PASSWORD",
					Value:       "",
					Usage:       "Email password",
					Destination: &cfg.Mail.Password,
				},
				cli.StringFlag{
					Name:        "mail-id",
					EnvVar:      "MAIL_ID",
					Value:       "",
					Usage:       "Email identity",
					Destination: &cfg.Mail.ID,
				},
			},
			Action: func(c *cli.Context) {
				cfg.MustValidate()
				commands.Serve(cfg)
			},
		},
		{
			Name:  "fetch",
			Usage: "Fetch new podcasts",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "url",
					EnvVar:      "DB_URL",
					Usage:       "Database connection URL",
					Destination: &cfg.DatabaseURL,
				},
				cli.IntFlag{
					Name:        "max-db-conns",
					Value:       99,
					EnvVar:      "MAx_DB_CONNECTIONS",
					Usage:       "Maximum database connections",
					Destination: &cfg.MaxDBConnections,
				},
			},
			Action: func(c *cli.Context) {
				cfg.MustValidate()
				commands.Fetch(cfg)
			},
		},
	}

	app.Run(os.Args)

}
