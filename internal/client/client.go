package client

import (
	"log"
	"net"
	"strconv"

	"github.com/mnlx/netlsr/internal/config"
	"github.com/mnlx/netlsr/internal/tun"
	"github.com/mnlx/netlsr/internal/utils"
)

func Client(config *config.ClientConfig) {
	iface, err := tun.SetupTun(config.TunName, config.LocalIP, config.PeerIP, config.TunCIDR)
	utils.CheckError(err, "SetupTun")

	server := net.JoinHostPort(config.ServerAddr, strconv.Itoa(config.Port))
	conn, err := net.Dial("udp", server)
	utils.CheckError(err, "dialing server")

	log.Printf("connected to server %s", server)

	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := iface.Read(buf)
			if err != nil {
				log.Printf("iface read: %v", err)
				return
			}
			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Printf("conn write: %v", err)
			}
		}
	}()

	buf := make([]byte, 1500)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read: %v", err)
			return
		}
		_, err = iface.Write(buf[:n])
		if err != nil {
			log.Printf("iface write: %v", err)
		}
	}
}
