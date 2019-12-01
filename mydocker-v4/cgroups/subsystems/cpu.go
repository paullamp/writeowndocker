package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct{}

func (s *CpuSubSystem) Name() string {
	return "cpu"
}

func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath,
		true); err == nil {
		if res.CpuSet != "" {
			if err := ioutil.WriteFile(path.Join(subsystemCgroupPath,
				"cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
				return fmt.Errorf("set cgroup cpu share fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	if subsystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath,
		false); err == nil {
		return os.RemoveAll(subsystemCgroupPath)
	} else {
		return err
	}
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	if subsystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath,
		false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsystemCgroupPath, "tasks"),
			[]byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
	} else {
		return err
	}
	return nil
}
