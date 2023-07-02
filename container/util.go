package container

import (
	"encoding/json"
	"os"
	"path"
	"wdocker/log"
)

func RecordContainer(con *Container) {
	b, err := json.Marshal(con)
	if err != nil {
		log.Error("json marshal err: %v", err)
	}
	jsonStr := string(b)
	log.Info(jsonStr)
	configURL := path.Join(con.URL, ConfigName)
	f, err := os.Create(configURL)
	if err != nil {
		log.Error("create file %s err: %v", configURL, err)
	}
	defer f.Close()
	f.WriteString(jsonStr)
}
