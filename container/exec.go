package container

import (
	"os"
	"os/exec"
	"strings"
	"wdocker/log"
)

var (
	ENV_EXEC_PID = "wdocker_pid"
	ENV_EXEC_CMD = "wdocker_cmd"
)

func ExecContainer(containerName string, cmds []string){
	con, err := GetContainerByName(containerName)
	if err != nil {
		log.Error("get container by name err: %v", err)
		return
	}
	pid := con.PID
	cmdStr := strings.Join(cmds, " ")
	log.Info("container id = %s, cmd = %s", pid, cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)

	err = cmd.Run()
	if err != nil {
		log.Error("exec container %s error %v", containerName, err)
	}
}