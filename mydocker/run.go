package main

import (
	// "mydocker/cgroups"
	// "mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// func Run(tty bool, command string) {
// 	parent := container.NewParentProcess(tty, command)
// 	if err := parent.Start(); err != nil {
// 		log.Errorf(err.Error())
// 	}
// 	parent.Wait()
// 	os.Exit(-1)
// }

// func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
// 	parent, writePipe := container.NewParentProcess(tty)
// 	if parent == nil {
// 		log.Errorf("new  parent process error")
// 		return
// 	}
// 	if err := parent.Start(); err != nil {
// 		log.Error(err)
// 	}

// 	//use mydocker-cgroup as cgroup name
// 	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
// 	defer cgroupManager.Destroy()
// 	cgroupManager.Set(res)
// 	cgroupManager.Apply(parent.Process.Pid)
// 	sendInitCommand(comArray, writePipe)
// 	parent.Wait()
// }

// make Run func for volumes
func Run(tty bool, comArray []string, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("new parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	mntURL := "/root/mnt"
	rootURL := "/root"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
