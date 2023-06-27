package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

type MemorySubsystem struct {
}

func (sys *MemorySubsystem) Name() string {
	return "memory"
}
func (sys *MemorySubsystem) SetResourceConfig(cgPath string, res *ResourceConfig) error {
	subsysCgPath, err := GetAbsCgPath(sys.Name(), cgPath, true)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(subsysCgPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
	if err != nil {
		return fmt.Errorf("set cg memory fail: %v", err)
	}

	return nil
}
func (sys *MemorySubsystem) AddProc(cgPath string, pid int) error {
	subsysCgPath, err := GetAbsCgPath(sys.Name(), cgPath, false)
	if err != nil {
		return fmt.Errorf("get cg %s error %v", cgPath, err)
	}
	err = os.WriteFile(path.Join(subsysCgPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("set cg proc fail: %v", err)
	}

	return nil
}

func (sys *MemorySubsystem) Remove(cgPath string) error {
	subsysCgPath, err := GetAbsCgPath(sys.Name(), cgPath, false)
	if err != nil {
		return err
	}
	logrus.Infof("susSysCgPath is %s.", subsysCgPath)
	return os.Remove(subsysCgPath)
}
