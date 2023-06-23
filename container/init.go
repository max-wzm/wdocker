package container

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var defaultMountFlags = syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

func RunContainerInitProcess(cmd string, args []string) error {
	log.Infof("init with command %s", cmd)
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{cmd}
	err := syscall.Exec(cmd, argv, os.Environ())
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
