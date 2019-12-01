package main

import (
	"fmt"
	"mydocker33/container"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, cmdArray []string) {
	parent, writePipe := container.NewParentProcess2(tty)
	if parent == nil {
		log.Errorf("new parent process error")
		return
	}
	fmt.Println("Run function called and parent have enabled")
	if err := parent.Start(); err != nil {
		log.Errorf(err.Error())
	}
	SendInitCommand(cmdArray, writePipe)
	parent.Wait()
	os.Exit(0)
}

func SendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("SendInitCommand ---- command all is : %s ", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
