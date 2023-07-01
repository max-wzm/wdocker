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

func Run(con *container.Container) error {
	initCmd, wPipe := container.NewInitCommand(con)
	if initCmd == nil {
		log.Error("new parent process error")
		return fmt.Errorf("new parent process error")
	}

	newWorkSpace(con)
	// workspace must be created and mounted before initCmd starts
	// since inidCmd starts with Mount Namespace, if u create workspace after initCmd, mount in newWorkspace() by parentProc does not work for initProc.
	// as a result, initProc only sees empty root content.
	initCmd.Dir = path.Join(con.URL, "mnt")

	err := initCmd.Start()
	if err != nil {
		log.Error("par proc start error: %v", err)
		return err
	}

	cgManger := cgroups.NewCgoupManager(con.ID)
	cgManger.SetResourceConfig(con.ResourceConfig)
	cgManger.AddProc(initCmd.Process.Pid)

	// after sending, init.go starts working.
	sendInitCommand(con.InitCmds, wPipe)

	if !con.RunningConfig.Detach {
		defer deleteWorkspace(con)
		defer cgManger.Destroy()

		err = initCmd.Wait()
		if err != nil {
			log.Error("parent wait error: %v", err)
			return err
		}
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
	con.URL = containerURL
	os.MkdirAll(con.URL, 0777)
	readURL := newReadLayer(containerURL, con.ImagePath)
	writeURL := newWriteLayer(containerURL)
	createMountPoint(containerURL, readURL, writeURL)
	MountVolume(con)
}

func MountVolume(con *container.Container) {
	volume := con.RunningConfig.Volume
	if volume == "" {
		return
	}
	volumnURLs := strings.Split(volume, ":")
	if len(volumnURLs) == 2 && volumnURLs[0] != "" && volumnURLs[1] != "" {
		parentVolURL := volumnURLs[0]
		os.MkdirAll(parentVolURL, 0777)
		containerVolURL := path.Join(con.URL, "mnt", volumnURLs[1])
		os.MkdirAll(containerVolURL, 0777)
		err := cmdRunStd("mount", "-t", "aufs", "-o", "dirs="+parentVolURL, "none", containerVolURL)
		if err != nil {
			log.Error("mount volumn err: %v", err)
		} else {
			log.Info("mount volume success: %q", volumnURLs)
		}
	} else {
		log.Error("extra volume mapping error: wrong format")
	}

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
	dirOpt := "dirs=" + writeURL + ":" + readURL + "=ro"
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
	if con.RunningConfig.Volume != "" {
		volumeURLs := strings.Split(con.RunningConfig.Volume, ":")
		containerVolURL := path.Join(con.URL, "mnt", volumeURLs[1])
		cmdRunStd("umount", containerVolURL)
	}

	mntURL := path.Join(con.URL, "mnt")
	writeLayerURL := path.Join(con.URL, "write_layer")
	cmdRunStd("umount", mntURL)
	os.RemoveAll(mntURL)
	os.RemoveAll(writeLayerURL)
	log.Info("removed worspace: %s & %s", mntURL, writeLayerURL)

	if con.RunningConfig.Remove {
		os.RemoveAll(con.URL)
	}
}

func cmdRunStd(name string, arg ...string) error {
	c := exec.Command(name, arg...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	return err
}
