package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = "mydocker33 is a simple container runtime implementation"

func main() {
	app := cli.NewApp()
	app.Name = "mydocker33"
	app.Usage = usage

	app.Commands = []*cli.Command{
		&initCommand,
		&runCommand,
	}
	app.Before = func(c *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
