package container

import (
	"fmt"
	"strconv"
	"syscall"
)

func StopContainer(name string) error {
	con, err := GetContainerByName(name)
	if err != nil {
		return err
	}
	pid, err := strconv.Atoi(con.PID)
	if err != nil {
		return fmt.Errorf("read pid err: %v", err)
	}
	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("kill error: %v", err)
	}
	con.Status = EXITED
	con.PID = " "
	RecordContainer(con)
	return nil
}
