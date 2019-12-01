package main

import (
	"fmt"
	"mydocker33/container"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "init a new container",
	Action: func(c *cli.Context) error {
		log.Infof("init come on")
		err := container.NewRunContainerInitProcess()
		return err
	},
}
var runCommand = cli.Command{
	Name:  "run",
	Usage: "run command in container",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return fmt.Errorf("Missing container command")
		}

		var cmdArray []string
		for _, arg := range c.Args().Slice() {
			cmdArray = append(cmdArray, arg)
		}
		fmt.Println("main-command file , Action; cmdArray is : ", cmdArray)
		tty := c.Bool("ti")
		Run(tty, cmdArray)
		return nil
	},
}
