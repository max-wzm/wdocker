package network

import (
	"fmt"
	"net"
	"wdocker/log"
)

func CreateNetwork(driver, subnet, name string) error {
	_, cidr, err := net.ParseCIDR(subnet)
	if err != nil {
		return err
	}
	gatewayIP, err := ipAllocator.Allocate(cidr)
	if err != nil {
		return err
	}
	cidr.IP = gatewayIP

	nw, err := drivers[driver].Create(cidr.String(), name)
	log.Info("created nw %v", nw.IpRange)
	if err != nil {
		return err
	}
	return nw.dump(defaultNetworkPath)
}

func DeleteNetwork(nwName string) error {
	nw, ok := networks[nwName]
	if !ok {
		return fmt.Errorf("no such nw")
	}
	err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP)
	if err != nil {
		return err
	}
	err = drivers[nw.Driver].Delete(nw)
	if err != nil {
		return err
	}
	return nw.remove(defaultNetworkPath)
}
