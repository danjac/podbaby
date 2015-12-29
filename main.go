package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/danjac/podbaby/commands"
	_ "github.com/lib/pq"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true

	var url, env, secretKey string
	var port int

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Run the server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "url",
					EnvVar:      "DB_URL",
					Usage:       "Database connection URL",
					Destination: &url,
				},
				cli.StringFlag{
					Name:        "secret",
					EnvVar:      "SECRET_KEY",
					Usage:       "Secret key",
					Destination: &secretKey,
				},
				cli.IntFlag{
					Name:        "port",
					Value:       5000,
					EnvVar:      "PORT",
					Usage:       "Server port",
					Destination: &port,
				},
				cli.StringFlag{
					Name:        "env",
					Value:       "prod",
					Usage:       "Environment",
					Destination: &env,
				},
			},
			Action: func(c *cli.Context) {
				if url == "" {
					log.Fatal("url is required")
				}
				if secretKey == "" {
					log.Fatal("secret is required")
				}
				commands.Serve(url, port, secretKey, env)
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
					Destination: &url,
				},
			},
			Action: func(c *cli.Context) {
				if url == "" {
					log.Fatal("url is required")
				}
				commands.Fetch(url)
			},
		},
	}

	app.Run(os.Args)

}
