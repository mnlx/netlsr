package main

import (
	"flag"
	"log"

	"github.com/mnlx/netlsr/internal/client"
	"github.com/mnlx/netlsr/internal/config"
	"github.com/mnlx/netlsr/internal/utils"
)

func main() {
	// mode := flag.String("mode", "client", "mode: client or server")
	// serverAddr := flag.String("remote", "", "server address for client mode")
	// tunName := flag.String("ifname", "utun99", "TUN interface name")
	// localIP := flag.String("local-ip", "", "local TUN IP, e.g. 10.100.0.1/16")
	// peerIP := flag.String("peer-ip", "", "peer TUN IP, e.g. 10.100.0.2")
	// port := flag.Int("port", 5000, "UDP port")
	// extIface := flag.String("ext-iface", "eth1", "external interface for NAT (server mode)")
	// debug := flag.Bool("debug", false, "debug mode")
	configPath := flag.String("config", "client.yaml", "config file path")
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	flag.Parse()

	// if *mode == "client" && *serverAddr == "" {
	// 	log.Fatal("remote server address required in client mode")
	// }

	config.Client.TunCIDR = utils.ExtractSubnetCIDR(config.Client.LocalIP)
	// if *localIP == "" || *peerIP == "" {
	// 	log.Fatal("local-ip and peer-ip are required")
	// }

	client.Client(&config.Client)
}
