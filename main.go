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
		cli.StringFlag{
			Name: "mysql, m",
			Value: "development_user:devboxpw@tcp(127.0.0.1:3306)/siil",
			Usage: "URL for mysql connection",
		},
	}

	app.Run(os.Args)
}
