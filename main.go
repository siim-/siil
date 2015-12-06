package main

import (
	"log"
	"os"

	"github.com/siim-/siil/server"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "siil"
	app.Usage = "eID authentication provider"

	app.Action = server.StartAPIServer

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8080,
			Usage: "Port the Siil API server will be listening on",
		},
		cli.StringFlag{
			Name:  "mysql, m",
			Value: "development_user:devboxpw@tcp(127.0.0.1:3306)/siil?parseTime=true",
			Usage: "URL for mysql connection",
		},
		cli.StringFlag{
			Name:  "working-directory, wd",
			Value: wd,
			Usage: "Working directory for the application",
		},
		cli.StringFlag{
			Name:  "site-id, sid",
			Value: "a1s2d34",
			Usage: "Site ID for the Siil entry",
		},
	}

	app.Run(os.Args)
}
