package network

import (
	"encoding/json"
	"net"
	"os"
	"path"
	"wdocker/utils"

	"github.com/vishvananda/netlink"
)

type Network struct {
	Name    string
	IpRange *net.IPNet
	Driver  string
}

type Endpoint struct {
	ID          string
	Device      netlink.Veth
	IP          net.IP
	Mac         net.HardwareAddr
	PortMapping []string
	Network     *Network
}

var networks = map[string]*Network{}

func (nw *Network) dump(dpath string) error {
	ok, err := utils.PathExists(dpath)
	if err != nil {
		return err
	}
	if !ok {
		os.MkdirAll(dpath, 0777)
	}
	nwPath := path.Join(dpath, nw.Name)
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer nwFile.Close()

	nwJson, err := json.Marshal(nw)
	if err != nil {
		return err
	}
	nwFile.Write(nwJson)
	return nil
}

func (nw *Network) load(nwPath string) error {
	b, err := os.ReadFile(nwPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, nw)
	if err != nil {
		return err
	}
	return nil
}

func (nw *Network) remove(dpath string) error {
	nwPath := path.Join(dpath, nw.Name)
	ok, err := utils.PathExists(nwPath)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return os.Remove(nwPath)
}
