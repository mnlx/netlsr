package main

import "github.com/songgao/water"

// TunSetup defines the interface for platform-specific TUN setup
type TunSetup interface {
	Setup(ifaceName, localIP, peerIP, tunCIDR string) (*water.Interface, error)
}

// Common configuration for TUN interface
type tunConfig struct {
	ifaceName string
	localIP   string
	peerIP    string
	tunCIDR   string
}
