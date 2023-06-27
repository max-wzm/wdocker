package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var defaultMountFlags = syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

func RunContainerInitProcess() error {
	log.Infof("start init proc")
	cmds := readCommands()
	if len(cmds) == 0 {
		return fmt.Errorf("failed to read commands, cmd array is nil")
	}
	log.Infof("init with commands %v", cmds)
	
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	
	path, err := exec.LookPath(cmds[0])
	if err != nil {
		log.Errorf("exec look path error: %v", err)
		return err
	}
	log.Infof("find path %v", path)

	err = syscall.Exec(path, cmds, os.Environ())
	if err != nil {
		log.Errorf("exec cmd %s error: %v", cmds[0], err)
		return err
	}
	return nil
}

func readCommands() []string {
	rPipe := os.NewFile(uintptr(3), "rPipe")
	cmdsByte, err := io.ReadAll(rPipe)
	if err != nil {
		log.Errorf("read commands from pipe: %v", err)
		return nil
	}

	cmds := string(cmdsByte)
	log.Infof("successfully read cmds: %s", cmds)

	return strings.Split(cmds, " ")
}
