package tun

import "github.com/songgao/water"

func SetupTun(ifaceName, localIP, peerIP, tunCIDR string) (*water.Interface, error) {
	tunSetup := newTunSetup()
	return tunSetup.Setup(ifaceName, localIP, peerIP, tunCIDR)
}
