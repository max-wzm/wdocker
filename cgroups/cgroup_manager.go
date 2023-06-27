package cgroups

import (
	"wdocker/cgroups/subsystems"

	"github.com/sirupsen/logrus"
)

type CgroupManager struct {
	Path     string
	Resource *subsystems.ResourceConfig
}

func NewCgoupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) AddProc(pid int) error {
	for _, subSys := range subsystems.Subsystems {
		subSys.AddProc(c.Path, pid)
	}
	return nil
}

func (c *CgroupManager) SetResourceConfig(res *subsystems.ResourceConfig) error {
	for _, subSys := range subsystems.Subsystems {
		subSys.SetResourceConfig(c.Path, res)
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, subSys := range subsystems.Subsystems {
		err := subSys.Remove(c.Path)
		if err != nil {
			logrus.Warnf("remove cgroup fail: %v", err)
		}
	}
	return nil
}
