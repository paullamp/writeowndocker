package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

// func RunContainerInitProcess(command string, args []string) error {
// 	logrus.Infof("command %s", command)
// 	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
// 	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
// 	argv := []string{command}
// 	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
// 		logrus.Errorf(err.Error())
// 	}
// 	return nil
// }

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdarray is nil")
	}

	// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	setUpMount()
	path, err := exec.LookPath(cmdArray[0])
	fmt.Println("The path found in cmdarray is : ", path)
	if err != nil {
		log.Errorf("exec loop path error %v", err)
		return err
	}
	log.Infof("Find Path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

/**
 INIT mount
**/

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("get current location field %v", err)
		return
	}

	log.Infof("Current location is :%s", pwd)
	pivoRoot(pwd)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivoRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("umount pivot_root dir %v", err)
	}
	return os.Remove(pivotDir)
}
