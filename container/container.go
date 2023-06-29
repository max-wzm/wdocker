package container

import "wdocker/cgroups/subsystems"

type Container struct {
	ID string
	Name string
	ImagePath string
	ResourceConfig *subsystems.ResourceConfig
	InitCmds []string
}