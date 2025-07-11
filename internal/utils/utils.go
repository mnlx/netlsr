package utils

import (
	"fmt"
	"log"
	"net"
)

func ExtractSubnetCIDR(cidr string) string {
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

func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
