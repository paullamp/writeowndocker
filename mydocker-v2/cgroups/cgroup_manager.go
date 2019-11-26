package cgroups

import (
	"mydocker-v2/cgroups/subsystems"

	"github.com/Sirupsen/logrus"
)

type CgroupManager struct {
	Path    string
	Resouce *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// set cgropu limit
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// add pid to cgroup

func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsystemIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// release cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsystemIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
