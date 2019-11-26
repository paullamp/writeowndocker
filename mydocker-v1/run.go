package main

import (
	"mydocker-v1/container"
	"os"

	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Errorf(err.Error())
	}
	parent.Wait()
	os.Exit(-1)
}
