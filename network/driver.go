package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/vishvananda/netlink"
)

type NetworkDriver interface {
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(nw *Network) error
	Connect(ep *Endpoint, nw *Network) error
	Disconnect(ep *Endpoint, nw *Network) error
}

var drivers = map[string]NetworkDriver{}

type BridgeNetworkDriver struct {
}

func (d *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (d *BridgeNetworkDriver) Create(subnet string, name string) (*Network, error) {
	ip, ipRange, _ := net.ParseCIDR(subnet)
	ipRange.IP = ip
	nw := &Network{
		Name:    name,
		IpRange: ipRange,
		Driver:  d.Name(),
	}
	err := d.init(nw)
	if err != nil {
		return nil, err
	}
	return nw, err
}

func (d *BridgeNetworkDriver) init(nw *Network) error {
	bName := nw.Name
	err := createBridgeInterface(bName)
	if err != nil {
		return fmt.Errorf("error init bridge %s: %v", bName, err)
	}

	gatewayIP := nw.IpRange
	gatewayIP.IP = nw.IpRange.IP
	err = setInterfaceIP(bName, gatewayIP.String())
	if err != nil {
		return fmt.Errorf("set interfacce ip %v err: %v", gatewayIP, err)
	}

	err = setUpBridge(bName)
	if err != nil {
		return fmt.Errorf("set up bridge %s err: %v", bName, err)
	}

	err = setUpIPTables(bName, nw.IpRange)
	if err != nil {
		return fmt.Errorf("set up ip tables err: %v", err)
	}
	return nil
}

func setUpIPTables(bName string, iPNet *net.IPNet) error {
	iptableCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", iPNet.String(), bName)
	cmd := exec.Command("iptables", strings.Split(iptableCmd, " ")...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("iptables cmd exec error with output %v: %v", output, err)
	}
	return nil
}

func setUpBridge(bName string) error {
	iface, err := netlink.LinkByName(bName)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(iface)
}

func setInterfaceIP(bName, ipNetStr string) error {
	iface, err := netlink.LinkByName(bName)
	if err != nil {
		return err
	}
	ipNet, err := netlink.ParseIPNet(ipNetStr)
	if err != nil {
		return err
	}
	addr := &netlink.Addr{
		IPNet:     ipNet,
		Peer:      ipNet,
		Label:     "",
		Flags:     0,
		Scope:     0,
		Broadcast: nil,
	}
	return netlink.AddrAdd(iface, addr)
}

func createBridgeInterface(bName string) error {
	_, err := net.InterfaceByName(bName)
	if err == nil || !strings.Contains(err.Error(), "no such network") {
		return fmt.Errorf("err in check interface %s: %v", bName, err)
	}
	la := netlink.NewLinkAttrs()
	la.Name = bName

	br := &netlink.Bridge{LinkAttrs: la}
	err = netlink.LinkAdd(br)
	if err != nil {
		return err
	}
	return nil
}

func (d *BridgeNetworkDriver) Delete(network *Network) error {
	bName := network.Name
	br, err := netlink.LinkByName(bName)
	if err != nil {
		return err
	}
	return netlink.LinkDel(br)
}

func (d *BridgeNetworkDriver) Connect(ep *Endpoint, nw *Network) error {
	bName := nw.Name
	br, err := netlink.LinkByName(bName)
	if err != nil {
		return err
	}

	la := netlink.NewLinkAttrs()
	la.Name = ep.ID[:5]
	la.MasterIndex = br.Attrs().Index

	ep.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  "cif-" + ep.ID[:5],
	}

	err = netlink.LinkAdd(&ep.Device)
	if err != nil {
		return fmt.Errorf("add veth ep err: %v", err)
	}
	err = netlink.LinkSetUp(&ep.Device)
	if err != nil {
		return fmt.Errorf("set up veth err: %v", err)
	}
	return nil
}

func (d *BridgeNetworkDriver) Disconnect(ep *Endpoint, nw *Network) error {
	return nil
}
