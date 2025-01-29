package cmd

import (
	"github.com/urfave/cli/v2"
	"barbecue/server"
)

var Serve = &cli.Command {
	Name: "serve",
	Usage: "...",
	Action: server.Serve,
	Flags: []cli.Flag {
		&cli.StringFlag {
			Name:  "ip",
			Value: "0.0.0.0",
			Usage: "...",
		},
		&cli.IntFlag {
			Name:  "port",
			Value: 8080,
			Usage: "...",
		},
	},
}
