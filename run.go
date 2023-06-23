package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"wdocker/container"
)

func Run(tty bool, cmd string) {
	parent := container.NewParentProcess(tty, cmd)
	err := parent.Start()
	if err != nil {
		log.Error(err)
	}
	parent.Wait()
	log.Infof("run.go - quit")
	os.Exit(-1)
}
