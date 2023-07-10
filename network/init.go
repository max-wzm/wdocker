package network

import (
	"os"
	"path"
	"path/filepath"
	"wdocker/log"
)

const defaultNetworkPath = "/wdocker/network/created"

func Init() error {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver

	if _, err := os.Stat(defaultNetworkPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(defaultNetworkPath, 0644)
		} else {
			return err
		}
	}

	filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, nwName := path.Split(nwPath)
		nw := &Network{
			Name: nwName,
		}

		if err := nw.load(nwPath); err != nil {
			log.Error("error load network: %s", err)
		}

		networks[nwName] = nw
		return nil
	})

	log.Info("networks: %v", networks)

	return nil
}
