package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/bin/sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS,
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	// if err := cmd.Run(); err != nil {
	// 	fmt.Println(err)
	// }
	cmd.Start()
	newroot := "/opt/busybox"
	oldroot := "/opt/busybox/.back"
	root := "/"
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		fmt.Printf("Mount rootfs to it self failed:%v", err)
	}
	if err := syscall.PivotRoot(newroot, oldroot); err != nil {
		fmt.Println("i m in pivotroot error")
		fmt.Println(err)
	}
	cmd.Wait()

}
