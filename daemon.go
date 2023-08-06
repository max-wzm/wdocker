package main

import (
	"wdocker/container"
	"wdocker/log"

	"github.com/max-wzm/geerpc"
)

type Daemon struct {
}

func (d *Daemon) RunContainer(con *container.Container, err *int) error {
	log.Info("eee")
	Run(con)
	return nil
}

func (d *Daemon) ListContainers(p geerpc.Args, err *error) error {
	container.ListContainers()
	return nil
}
