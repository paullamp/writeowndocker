package container

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) error {
	fmt.Println("container-init: command is : ", command)
	fmt.Println("container-init: args is : ", args)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		fmt.Println("i am error")
		log.Fatal(err)
	}
	return nil
}
