package main

import (
	"flag"

	"github.com/mnlx/netlsr/internal/server"
	"github.com/mnlx/netlsr/internal/utils"
)

func main() {
	// mode := flag.String("mode", "client", "mode: client or server")
	// serverAddr := flag.String("remote", "", "server address for client mode")
	tunName := flag.String("ifname", "utun99", "TUN interface name")
	localIP := flag.String("local-ip", "", "local TUN IP, e.g. 10.100.0.1/16")
	peerIP := flag.String("peer-ip", "", "peer TUN IP, e.g. 10.100.0.2")
	port := flag.Int("port", 5000, "UDP port")
	extIface := flag.String("ext-iface", "eth1", "external interface for NAT (server mode)")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	// if *mode == "client" && *serverAddr == "" {
	// 	log.Fatal("remote server address required in client mode")
	// }

	tunCIDR := utils.ExtractSubnetCIDR(*localIP)
	// if *localIP == "" || *peerIP == "" {
	// 	log.Fatal("local-ip and peer-ip are required")
	// }

	server.Server(*tunName, *localIP, *peerIP, tunCIDR, *port, *extIface, *debug)
}
