package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"wdocker/log"
)

var defaultMountFlags = syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

func RunContainerInitProcess() error {
	log.Info("start init proc")
	cmds := readCommands()
	if len(cmds) == 0 {
		return fmt.Errorf("failed to read commands, cmd array is nil")
	}
	log.Info("init with commands %v", cmds)

	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	path, err := exec.LookPath(cmds[0])
	if err != nil {
		log.Error("exec look path error: %v", err)
		return err
	}
	log.Info("find path %v", path)

	err = syscall.Exec(path, cmds, os.Environ())
	if err != nil {
		log.Error("exec cmd %s error: %v", cmds[0], err)
		return err
	}
	return nil
}

func readCommands() []string {
	rPipe := os.NewFile(uintptr(3), "rPipe")
	cmdsByte, err := io.ReadAll(rPipe)
	if err != nil {
		log.Error("read commands from pipe: %v", err)
		return nil
	}

	cmds := string(cmdsByte)
	log.Info("successfully read cmds: %s", cmds)

	return strings.Split(cmds, " ")
}
