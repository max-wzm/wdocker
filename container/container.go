package container

import "wdocker/cgroups/subsystems"

type RunningConfig struct {
	Remove bool
	Volume string
}

type Container struct {
	ID string
	Name string
	ImagePath string
	URL string
	ResourceConfig *subsystems.ResourceConfig
	RunningConfig *RunningConfig
	InitCmds []string
}