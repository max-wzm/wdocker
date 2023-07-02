package container

import (
	"wdocker/cgroups/subsystems"
)

var (
	RUNNING    = "running"
	STOPPED    = "stopped"
	EXITED     = "exited"
	ConfigName = "config.json"
	LogName    = "container.log"
)

type RunningConfig struct {
	Tty    bool
	Remove bool
	Detach bool
	Volume string
}

type ContainerInfo struct {
	ID          string
	Name        string
	PID         string
	Status      string
	InitCmd     string
	CreatedTime string
}

type Container struct {
	ContainerInfo
	ImagePath      string
	URL            string
	ResourceConfig *subsystems.ResourceConfig
	RunningConfig  *RunningConfig
}
