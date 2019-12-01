package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSetSubSystem struct{}

func (s *CpuSetSubSystem) Name() string {
	return "cpuset"
}

func (s *CpuSetSubSystem) Set(cgrouppath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgrouppath, true); err == nil {
		if res.CpuSet != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"),
				[]byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("set cgroup cpuset fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *CpuSetSubSystem) Remove(cgrouppath string) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgrouppath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

func (s *CpuSetSubSystem) Apply(cgrouppath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgrouppath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"),
			[]byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgrouppath, err)
	}
}
