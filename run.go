package main

import (
	"fmt"
	"os"
	"strings"

	"wdocker/cgroups"
	"wdocker/container"
	"wdocker/log"
)

func Run(con *container.Container, tty bool) error {
	parent, wPipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Error("new parent process error")
		return fmt.Errorf("new parent process error")
	}

	err := parent.Start()
	if err != nil {
		log.Error("par proc start error: %v", err)
		return err
	}

	res := con.ResourceConfig
	cmds := con.InitCmds

	cgManger := cgroups.NewCgoupManager(con.ID)
	defer cgManger.Destroy()
	cgManger.SetResourceConfig(res)
	cgManger.AddProc(parent.Process.Pid)

	sendInitCommand(cmds, wPipe)

	err = parent.Wait()
	if err != nil {
		log.Error("parent wait error: %v", err)
		return err
	}
	log.Info("run.go - quit")
	return nil
}

func sendInitCommand(cmds []string, wPipe *os.File) {
	log.Info("sending init cmd...")
	command := strings.Join(cmds, " ")
	wPipe.WriteString(command)
	wPipe.Close()
}
