package main

import (
	"os"
	"strings"

	"wdocker/cgroups"
	"wdocker/cgroups/subsystems"
	"wdocker/container"
	"wdocker/log"
)

func Run(tty bool, cmds []string, res *subsystems.ResourceConfig) {
	parent, wPipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Error("new parent process error")
		return
	}

	err := parent.Start()
	if err != nil {
		log.Error("par proc start error: %v", err)
	}

	cgManger := cgroups.NewCgoupManager("wdocker-cgroup")
	defer cgManger.Destroy()

	cgManger.SetResourceConfig(res)
	cgManger.AddProc(parent.Process.Pid)

	sendInitCommand(cmds, wPipe)

	err = parent.Wait()
	if err != nil {
		log.Error("parent wait error: %v", err)
	}
	log.Info("run.go - quit")
}

func sendInitCommand(cmds []string, wPipe *os.File) {
	log.Info("sending init cmd...")
	command := strings.Join(cmds, " ")
	wPipe.WriteString(command)
	wPipe.Close()
}
