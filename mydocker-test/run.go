package main

import (
	"log"
	"mydocker-v1/container"
	"os"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Fatal(err)
	}
	parent.Wait()
	os.Exit(-1)
}
