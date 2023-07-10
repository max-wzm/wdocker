package network

import (
	"encoding/json"
	"net"
	"os"
	"path"
	"strings"
	"wdocker/log"
	"wdocker/utils"
)

// network ip addr allocator

const ipamDefaultAllocatorPath = "/wdocker/network/ipam/subnet.json"

var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
	Subnets:             make(map[string]string),
}

type IPAM struct {
	SubnetAllocatorPath string
	Subnets             map[string]string
}

func (ipam *IPAM) load() error {
	exists, err := utils.PathExists(ipam.SubnetAllocatorPath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	b, err := os.ReadFile(ipam.SubnetAllocatorPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &ipam.Subnets)
	if err != nil {
		return err
	}
	return nil
}

func (ipam *IPAM) dump() error {
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	exists, err := utils.PathExists(ipamConfigFileDir)
	if err != nil {
		return err
	}
	if !exists {
		os.MkdirAll(ipamConfigFileDir, 0777)
	}

	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer subnetConfigFile.Close()

	configJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		return err
	}

	_, err = subnetConfigFile.Write(configJson)
	return err
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (net.IP, error) {
	subnet.IP = subnet.IP.To4()
	ipam.Subnets = map[string]string{}
	err := ipam.load()
	if err != nil {
		return nil, err
	}
	maskSize, size := subnet.Mask.Size()
	_, ok := ipam.Subnets[subnet.String()]
	if !ok {
		ipam.Subnets[subnet.String()] = strings.Repeat("0", 1<<uint8(size-maskSize))
	}
	ip := subnet.IP
	//
	for c := range (ipam.Subnets)[subnet.String()] {
		if (ipam.Subnets)[subnet.String()][c] == '0' {
			ipalloc := []byte((ipam.Subnets)[subnet.String()])
			ipalloc[c] = '1'
			(ipam.Subnets)[subnet.String()] = string(ipalloc)
			for t := uint(4); t > 0; t -= 1 {
				ip[4-t] += uint8(c >> ((t - 1) * 8))
			}
			ip[3] += 1
			break
		}
	}
	log.Info("ip %v", ip)
	ipam.dump()
	return ip, nil
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = map[string]string{}

	_, subnet, _ = net.ParseCIDR(subnet.String())

	err := ipam.load()
	if err != nil {
		log.Error("Error dump allocation info, %v", err)
	}

	c := 0
	releaseIP := ipaddr.To4()
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t -= 1 {
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}

	ipalloc := []byte((ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(ipam.Subnets)[subnet.String()] = string(ipalloc)

	ipam.dump()
	return nil
}
