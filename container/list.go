package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"text/tabwriter"
	"wdocker/log"
)

func ListContainers() {
	dirEntries, _ := os.ReadDir("/wdocker")
	containers := make([]*Container, 0)
	for _, e := range dirEntries {
		fi, _ := e.Info()
		con, _ := GetContainerByName(fi.Name())
		if con == nil {
			continue
		}
		containers = append(containers, con)
	}
	log.Info("%v", containers)
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, con := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", con.ID, con.Name, con.PID, con.Status, con.InitCmd, con.CreatedTime)
	}
	w.Flush()
}

func GetContainerByName(name string) (*Container, error) {
	configURL := path.Join("/wdocker", name, ConfigName)
	b, err := os.ReadFile(configURL)
	if err != nil {
		return nil, err
	}
	var tmpCon Container
	json.Unmarshal(b, &tmpCon)
	return &tmpCon, nil
}
