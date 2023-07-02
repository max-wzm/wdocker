package container

import (
	"fmt"
	"os"
	"path"
	"strings"
	"wdocker/utils"
)

func RemoveContainer(name string) error {
	con, err := GetContainerByName(name)
	if err != nil {
		return fmt.Errorf("get container err: %v ", err)
	}
	if con.Status != EXITED {
		return fmt.Errorf("status err: remove only for exited container")
	}
	if con.RunningConfig.Volume != "" {
		volumeURLs := strings.Split(con.RunningConfig.Volume, ":")
		containerVolURL := path.Join(con.URL, "mnt", volumeURLs[1])
		utils.CmdRunStd("umount", containerVolURL)
	}
	mntURL := path.Join(con.URL, "mnt")
	utils.CmdRunStd("umount", mntURL)
	err = os.RemoveAll(con.URL)
	return err
}
