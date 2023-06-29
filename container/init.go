package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
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

	err := setUpMount()
	if err != nil {
		log.Error("set up mount err: %v", err)
		return err
	}

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

func setUpMount() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Info("current working dir is %s", dir)

	err = pivotRoot(dir)
	if err != nil {
		log.Error("pivot_root: %v", err)
		return err
	}
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
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

func pivotRoot(newRoot string) error {
	err := syscall.Mount(newRoot, newRoot, "bind", syscall.MS_BIND, "")
	if err != nil {
		return fmt.Errorf("mount bind %s error: %v", newRoot, err)
	}

	putOld := path.Join(newRoot, ".pivot_root")
	err = os.Mkdir(putOld, 0777)
	if err != nil {
		return err
	}
	err = syscall.PivotRoot(newRoot, putOld)
	if err != nil {
		return err
	}
	err = syscall.Unmount("/.pivot_root", syscall.MNT_DETACH)
	if err != nil {
		return err
	}
	err = os.Remove("/.pivot_root")
	if err != nil {
		return err
	}

	syscall.Chdir("/")
	return nil
}
