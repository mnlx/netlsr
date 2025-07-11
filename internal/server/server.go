package server

import (
	"log"
	"net"
	"os/exec"

	"github.com/mnlx/netlsr/internal/config"
	"github.com/mnlx/netlsr/internal/tun"
	"github.com/mnlx/netlsr/internal/utils"
)

func Server(config *config.ServerConfig) {
	iface, err := tun.SetupTun(config.TunName, config.LocalIP, config.PeerIP, config.TunCIDR)
	utils.CheckError(err, "setupTun")

	setupNAT(config.ExtIface)

	serverAddr := net.UDPAddr{Port: config.Port}
	conn, err := net.ListenUDP("udp", &serverAddr)

	utils.CheckError(err, "listening UDP")

	log.Printf("listening for client on %s", serverAddr.String())

	buf := make([]byte, 1500)
	n, clientAddr, err := conn.ReadFromUDP(buf)
	utils.CheckError(err, "reading initial packet")
	log.Printf("client address: %s", clientAddr.String())
	_, err = conn.WriteTo(buf[:n], clientAddr)
	utils.CheckError(err, "writing initial packet")

	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := iface.Read(buf)
			if err != nil {
				log.Printf("iface read: %v", err)
				return
			}
			if config.Debug {
				log.Printf("receiving packet from %s", clientAddr.String())
			}
			_, err = conn.WriteTo(buf[:n], clientAddr)
			if err != nil {
				log.Printf("conn write: %v", err)
			}
		}
	}()

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("conn read: %v", err)
			return
		}

		if config.Debug {
			log.Printf("sending packet to %s", addr.String())
		}

		if addr.String() != clientAddr.String() {
			continue
		}
		_, err = iface.Write(buf[:n])
		if err != nil {
			log.Printf("iface write: %v", err)
		}
	}
}

func setupNAT(extIface string) {
	// enable IP forwarding
	cmd := exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("warning: enabling IP forwarding: %v, output: %s", err, out)
	}

	// configure NAT
	cmd = exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-s", "10.177.0.0/24", "-o", extIface, "-j", "MASQUERADE")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("warning: configuring NAT: %v, output: %s", err, out)
	}
}
