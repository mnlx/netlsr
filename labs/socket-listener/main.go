package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	// Create a packet capture handle on en0 (change if needed)
	handle, err := pcap.OpenLive("en0", 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal("Error creating packet capture:", err)
	}
	defer handle.Close()

	// Packet source
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		// Parse Ethernet layer
		if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
			eth := ethLayer.(*layers.Ethernet)
			fmt.Println("Ethernet:", eth.SrcMAC, "->", eth.DstMAC)
		}

		// Parse IP layer
		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ip := ipLayer.(*layers.IPv4)
			fmt.Println("IPv4:", ip.SrcIP, "->", ip.DstIP)
		}

		// Parse TCP layer
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)
			fmt.Printf("TCP: %d -> %d\n", tcp.SrcPort, tcp.DstPort)

			// Check for HTTP in the payload
			appLayer := packet.ApplicationLayer()
			if appLayer != nil {
				// payload := string(appLayer.Payload())
				fmt.Println("---- HTTP ----")
				fmt.Println(appLayer.LayerContents())
				fmt.Println("--------------")
				// // Rough check for HTTP content
				// if strings.HasPrefix(payload, "GET") || strings.HasPrefix(payload, "POST") ||
				// 	strings.HasPrefix(payload, "HTTP/") {

				// }
			}
		}

		// Parse UDP layer
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp := udpLayer.(*layers.UDP)
			fmt.Printf("UDP: %d -> %d Length: %d\n", udp.SrcPort, udp.DstPort, udp.Length)

			// Check UDP payload
			appLayer := packet.ApplicationLayer()
			if appLayer != nil {
				payload := string(appLayer.Payload())
				if len(payload) > 0 {
					fmt.Println("---- UDP Payload ----")
					fmt.Println(payload)
					fmt.Println("-------------------")
				}
			}
		}
	}
}
