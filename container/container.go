package container

import "wdocker/cgroups/subsystems"

type RunningConfig struct {
	Tty bool
	Remove bool
	Detach bool
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