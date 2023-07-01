package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"wdocker/utils"
)

// find the dir of the root cg, to which a subsystem is attached.
func FindCgMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")
		opts := strings.Split(fields[len(fields)-1], ",")
		for _, opt := range opts {
			if opt == subsystem {
				return fields[4]
			}
		}
	}

	return ""
}

func GetAbsCgPath(subsystem string, cgPath string, autoCreate bool) (string, error) {
	cgRoot := FindCgMountPoint(subsystem)
	absCgPath := path.Join(cgRoot, cgPath)

	exists, err := utils.PathExists(absCgPath)
	if err != nil {
		return "", fmt.Errorf("cg path error %v", err)
	}

	if exists {
		return absCgPath, nil
	}

	if autoCreate {
		err := os.Mkdir(absCgPath, 0755)
		if err != nil {
			return "", fmt.Errorf("error create cg %v", err)
		}
		return absCgPath, nil
	}

	return "", fmt.Errorf("cg path error %v", err)
}
