package container

import (
	"os"
	"os/exec"
	"syscall"
	"wdocker/log"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	rPipe, wPipe, err := NewPipe()
	if err != nil {
		log.Error("new pipe error: %v", err)
		return nil, nil
	}
	proc, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Error("get init proc error: %v", err)
		return nil, nil
	}

	cmd := exec.Command(proc, "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.ExtraFiles = []*os.File{rPipe}
	cmd.Dir = "/root/busybox"

	return cmd, wPipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
