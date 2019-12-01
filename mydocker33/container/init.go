package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func NewRunContainerInitProcess() error {
	fmt.Println("i am in runcontainer")
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("run container get user command error, command array is nil")
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("exec look path error: %v", err)
		return err
	}
	log.Infof("Find path is %s", path)
	log.Infof("cmdArray is : ", cmdArray)
	log.Infof("cmdArray[0:] is :", cmdArray[0:])
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error: %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
