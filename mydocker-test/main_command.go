package main

import (
	"fmt"
	"log"
	"mydocker-test/container"

	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user'process in container. Donot call it outside",
	Action: func(c *cli.Context) error {
		log.Println("init come on")
		cmd := c.Args().Get(0)
		fmt.Println("command is :", cmd)
		if err := container.RunContainerInitProcess(cmd, nil); err != nil {
			log.Fatal("init continer failed")
			return err
		}
		return nil
	},
}

var runCommand = cli.Command{
	Name: "run",
	Usage: `create a continer with namespace and cgroup limit
			mydocker-test run -ti [command]
		`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},

	Action: func(c *cli.Context) error {
		fmt.Println(c.NArg())
		fmt.Println(c.Args())
		if c.NArg() < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := c.Args().Get(0)
		fmt.Println("cmd is :", cmd)
		tty := c.Bool("ti")
		fmt.Println("tty is :", tty)
		Run(tty, cmd)
		return nil
	},
}
