package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const usage = `mydocker is a simple container runtime implement
			  The purpose of this project is to learn how  docker works and how to wirte  docker by ourselves
			  Enjoy it ,just for fun.
	`

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage
	app.Commands = []*cli.Command{
		&initCommand,
		&runCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
