//go:build darwin

package main

import (
	"fmt"
	"os/exec"

	"github.com/songgao/water"
)

type darwinTunSetup struct{}

func newTunSetup() TunSetup {
	return &darwinTunSetup{}
}

func (d *darwinTunSetup) Setup(ifaceName, localIP, peerIP, tunCIDR string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: ifaceName,
		},
	}
	iface, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("creating TUN interface: %v", err)
	}

	// assign IP address and peer
	cmd := exec.Command("ifconfig", iface.Name(), "inet", localIP, peerIP)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("assigning IP address: %v, output: %s", err, out)
	}

	// bring up
	cmd = exec.Command("ifconfig", iface.Name(), "up")
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("setting interface up: %v, output: %s", err, out)
	}

	// add route for tun network
	cmd = exec.Command("route", "add", tunCIDR, "-interface", iface.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("adding route: %v, output: %s", err, out)
	}

	return iface, nil
}
