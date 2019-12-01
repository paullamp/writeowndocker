package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = "Mydocker is a simple container runtime implement"

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []*cli.Command{
		&initCommand,
		&runCommand,
		&commitCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
