package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"wdocker/cgroups"
	"wdocker/container"
	"wdocker/log"
)

func Run(con *container.Container, tty bool) error {
	initCmd, wPipe := container.NewInitCommand(con, tty)
	if initCmd == nil {
		log.Error("new parent process error")
		return fmt.Errorf("new parent process error")
	}

	newWorkSpace(con)
	defer deleteWorkspace(con)
	// workspace must be created and mounted before initCmd starts
	// since inidCmd starts with Mount Namespace, if u create workspace after initCmd, mount in newWorkspace() by parentProc does not work for initProc.
	// as a result, initProc only sees empty root content.
	initCmd.Dir = con.Root

	err := initCmd.Start()
	if err != nil {
		log.Error("par proc start error: %v", err)
		return err
	}

	cgManger := cgroups.NewCgoupManager(con.ID)
	defer cgManger.Destroy()
	cgManger.SetResourceConfig(con.ResourceConfig)
	cgManger.AddProc(initCmd.Process.Pid)

	// after sending, init.go starts working.
	sendInitCommand(con.InitCmds, wPipe)

	err = initCmd.Wait()

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

func newWorkSpace(con *container.Container) {
	containerURL := path.Join("/wdocker", con.ID)
	os.Mkdir("/wdocker", 0777)
	os.Mkdir(containerURL, 0777)
	readURL := newReadLayer(containerURL, con.ImagePath)
	writeURL := newWriteLayer(containerURL)
	con.Root = createMountPoint(containerURL, readURL, writeURL)
	con.URL = containerURL
}

func newWriteLayer(containerURL string) string {
	writeURL := path.Join(containerURL, "write_layer")
	os.Mkdir(writeURL, 0777)
	return writeURL
}

func newReadLayer(containerURL, imagePath string) string {
	readURL := path.Join(containerURL, "read_layer")
	os.Mkdir(readURL, 0777)
	_, err := exec.Command("tar", "-xvf", imagePath, "-C", readURL).CombinedOutput()
	if err != nil {
		log.Error("untar %s error: %v", imagePath, err)
	}
	return readURL
}

func createMountPoint(containerURL, readURL, writeURL string) string {
	mntURL := path.Join(containerURL, "mnt")
	os.Mkdir(mntURL, 0777)
	dirOpt := "br:" + writeURL + ":" + readURL + "=ro"
	mntCmd := exec.Command("mount", "-t", "aufs", "-o", dirOpt, "none", mntURL)
	mntCmd.Stdout = os.Stdout
	mntCmd.Stderr = os.Stderr
	err := mntCmd.Run()
	if err != nil {
		log.Error("run mount aufs with dirOpt: %s error: %v", dirOpt, err)
	}
	return mntURL
}

func deleteWorkspace(con *container.Container) {
	c := exec.Command("umount", con.Root)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Run()
	err := os.RemoveAll(con.URL)
	if err != nil {
		log.Error("remove workspace err: %v", err)
	}
	log.Info("removed worspace: %s", con.URL)
}
