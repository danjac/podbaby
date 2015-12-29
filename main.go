package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/danjac/podbaby/commands"
	_ "github.com/lib/pq"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true

	var url, env string
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
				commands.Serve(url, port, env)
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
				commands.Fetch(url)
			},
		},
	}

	app.Run(os.Args)

}
