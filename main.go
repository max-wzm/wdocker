package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	usage = `mydocker is a simple container runtime implementation.
	The purpose of this project is to learn how docker works and how to write a docker by ourselves
	Enjoy it, just for fun.`
)

func main() {
	app := cli.NewApp()
	app.Name = "wdocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		listCommand,
	}

	app.Before = func(ctx *cli.Context) error {
		log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
		})
		log.SetOutput(os.Stdout)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
