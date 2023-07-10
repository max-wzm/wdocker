package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"wdocker/container"
	"wdocker/log"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func Connect(nwName string, con *container.Container) error {
	nw, ok := networks[nwName]
	if !ok {
		return fmt.Errorf("no such network")
	}
	log.Info("nw iprange: %v", nw)
	ip, err := ipAllocator.Allocate(nw.IpRange)
	if err != nil {
		return fmt.Errorf("alloc err: %v", err)
	}
	log.Info("nw iprange: %v", nw)
	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", con.ID, nwName),
		IP:          ip,
		Network:     nw,
		PortMapping: con.PortMapping,
	}

	err = drivers[nw.Driver].Connect(ep, nw)
	if err != nil {
		return fmt.Errorf("connect err: %v", err)
	}
	err = configEndpointIPAddrAndRoute(ep, con)
	if err != nil {
		return fmt.Errorf("config ep err: %v", err)
	}
	return configPortMapping(ep, con)
}

func configPortMapping(ep *Endpoint, con *container.Container) error {
	for _, pm := range ep.PortMapping {
		portMapping := strings.Split(pm, ":")
		if len(portMapping) != 2 {
			log.Error("port mapping fmt err")
			continue
		}
		iptableCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s", portMapping[0], ep.IP.String(), portMapping[1])
		cmd := exec.Command("iptables", strings.Split(iptableCmd, " ")...)
		output, err := cmd.Output()
		if err != nil {
			log.Error("iptable ouput err: %v", output)
			continue
		}
	}
	return nil
}

func configEndpointIPAddrAndRoute(ep *Endpoint, con *container.Container) error {
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)
	if err != nil {
		return fmt.Errorf("find peer link err: %v", err)
	}
	defer enterConNetns(&peerLink, con)()

	ifaceIP := ep.Network.IpRange
	ifaceIP.IP = ep.IP
	log.Info("iface %v", ep.Network.IpRange)
	err = setInterfaceIP(ep.Device.PeerName, ifaceIP.String())
	if err != nil {
		log.Info("ep: %v", ep)
		return fmt.Errorf("set iface ip err: %v", err)
	}
	log.Info("a")

	err = setUpBridge(ep.Device.PeerName)
	if err != nil {
		return err
	}
	log.Info("b")

	err = setUpBridge("lo")
	if err != nil {
		return err
	}
	log.Info("c")

	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")
	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        ep.Network.IpRange.IP,
		Dst:       cidr,
	}
	err = netlink.RouteAdd(defaultRoute)
	if err != nil {
		log.Error("%v", err)
	}
	return nil
}

func enterConNetns(link *netlink.Link, con *container.Container) func() {
	log.Info("##########")
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", con.PID), os.O_RDONLY, 0)
	if err != nil {
		log.Error("enter ns err: %v", err)
		return nil
	}
	nsFD := f.Fd()
	runtime.LockOSThread()
	err = netlink.LinkSetNsFd(*link, int(nsFD))
	if err != nil {
		log.Error("set ns fd err: %v", err)
		return nil
	}
	originNS, err := netns.Get()
	if err != nil {
		log.Error("get cur net ns err: %v", err)
		return nil
	}
	err = netns.Set(netns.NsHandle(nsFD))
	if err != nil {
		return nil
	}
	return func() {
		netns.Set(originNS)
		originNS.Close()
		runtime.UnlockOSThread()
		f.Close()
	}
	
}
