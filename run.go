package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"wdocker/cgroups"
	"wdocker/cgroups/subsystems"
	"wdocker/container"
)

func Run(tty bool, cmds []string, res *subsystems.ResourceConfig) {
	parent, wPipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("new parent process error")
		return
	}

	err := parent.Start()
	if err != nil {
		log.Error(err)
	}

	cgManger := cgroups.NewCgoupManager("wdocker-cgroup")
	defer cgManger.Destroy()

	cgManger.SetResourceConfig(res)
	cgManger.AddProc(parent.Process.Pid)

	sendInitCommand(cmds, wPipe)

	parent.Wait()
	log.Infof("run.go - quit")
}

func sendInitCommand(cmds []string, wPipe *os.File){
	log.Infof("sending init cmd...")
	command := strings.Join(cmds, " ")
	wPipe.WriteString(command)
	wPipe.Close()
}
