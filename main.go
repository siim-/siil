package main

import (
	"os"

	"github.com/siim-/siil/server"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "siil"
	app.Usage = "eID OAuth2 provider"

	app.Action = server.StartAPIServer

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8080,
			Usage: "Port the Siil API server will be listening on",
		},
	}

	app.Run(os.Args)
}
