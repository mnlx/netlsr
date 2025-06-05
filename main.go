package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"

	"github.com/songgao/water"
)

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func setupTun(ifaceName, localIP, peerIP, tunCIDR string) (*water.Interface, error) {
	tunSetup := newTunSetup()
	return tunSetup.Setup(ifaceName, localIP, peerIP, tunCIDR)
}

func clientMode(ifaceName, localIP, peerIP, tunCIDR, serverAddr string, port int) {
	iface, err := setupTun(ifaceName, localIP, peerIP, tunCIDR)
	checkError(err, "setupTun")

	server := net.JoinHostPort(serverAddr, strconv.Itoa(port))
	conn, err := net.Dial("udp", server)
	checkError(err, "dialing server")

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

func serverMode(ifaceName, localIP, peerIP, tunCIDR string, port int, extIface string, debug bool) {
	iface, err := setupTun(ifaceName, localIP, peerIP, tunCIDR)
	checkError(err, "setupTun")

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

	addr := net.UDPAddr{Port: port}
	conn, err := net.ListenUDP("udp", &addr)

	checkError(err, "listening UDP")

	log.Printf("listening for client on %s", addr.String())

	buf := make([]byte, 1500)
	n, clientAddr, err := conn.ReadFromUDP(buf)
	checkError(err, "reading initial packet")
	log.Printf("client address: %s", clientAddr.String())
	_, err = conn.WriteTo(buf[:n], clientAddr)
	checkError(err, "writing initial packet")

	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := iface.Read(buf)
			if err != nil {
				log.Printf("iface read: %v", err)
				return
			}
			if debug {
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

		if debug {
			log.Printf("sending packet to %s", addr.String())
		}

		if addr.String() != clientAddr.String() {
			continue
		}
		_, err = iface.Write(buf[:n])
		if err != nil {
			log.Printf("iface write: %v", err)
		}

		// conn, err := net.Dial("ip4:icmp", "10.100.0.1")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// conn.Write([]byte{ /* raw IP packet */ })
	}
}

func main() {
	mode := flag.String("mode", "client", "mode: client or server")
	serverAddr := flag.String("remote", "", "server address for client mode")
	tunName := flag.String("ifname", "utun99", "TUN interface name")
	localIP := flag.String("local-ip", "", "local TUN IP, e.g. 10.100.0.1/16")
	peerIP := flag.String("peer-ip", "", "peer TUN IP, e.g. 10.100.0.2")
	port := flag.Int("port", 5000, "UDP port")
	extIface := flag.String("ext-iface", "eth1", "external interface for NAT (server mode)")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	if *mode == "client" && *serverAddr == "" {
		log.Fatal("remote server address required in client mode")
	}

	tunCIDR := extractSubnetCIDR(*localIP)
	// if *localIP == "" || *peerIP == "" {
	// 	log.Fatal("local-ip and peer-ip are required")
	// }

	if *mode == "client" {
		clientMode(*tunName, *localIP, *peerIP, tunCIDR, *serverAddr, *port)
	} else {
		serverMode(*tunName, *localIP, *peerIP, tunCIDR, *port, *extIface, *debug)
	}
}

func extractSubnetCIDR(cidr string) string {
	ip, ipNet, err := net.ParseCIDR(cidr)
	fmt.Println("ip", ip)
	fmt.Println("ipNet", ipNet)
	if err != nil {
		panic(err)
	}

	// Apply the network mask to get the subnet base (e.g. 10.177.0.0)
	networkIP := ip.Mask(ipNet.Mask)
	maskSize, _ := ipNet.Mask.Size()

	return fmt.Sprintf("%s/%d", networkIP.String(), maskSize)
}
