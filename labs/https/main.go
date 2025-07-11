package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		// You can add more config here for debugging, mTLS, etc
	}

	// Listen for raw TCP connections
	ln, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on :8443")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept:", err)
			continue
		}

		network := ln.Addr().Network()
		address := ln.Addr().String()
		fmt.Println("network:", network)
		fmt.Println("address:", address)

		go handleConn(conn, tlsConfig)
	}
}

func handleConn(rawConn net.Conn, tlsConfig *tls.Config) {
	defer rawConn.Close()

	// Manually wrap in tls.Conn but don't handshake yet
	tlsConn := tls.Server(rawConn, tlsConfig)

	// Here, you could set deadlines, wrap the handshake in tracing, etc
	tlsConn.SetDeadline(time.Now().Add(60 * time.Second))
	err := tlsConn.Handshake()
	if err != nil {
		fmt.Println("TLS handshake failed:", err)
		return
	}
	// After handshake, you can inspect state:
	state := tlsConn.ConnectionState()
	fmt.Printf("TLS Version: %x, CipherSuite: %x\n", state.Version, state.CipherSuite)
	fmt.Printf("HandshakeComplete: %t\n", state.HandshakeComplete)
	fmt.Printf("NegotiatedProtocol: %s\n", state.NegotiatedProtocol)
	fmt.Printf("ServerName: %s\n", state.ServerName)
	fmt.Printf("PeerCertificates: %v\n", state.PeerCertificates)
	fmt.Printf("VerifiedChains: %v\n", state.VerifiedChains)
	fmt.Printf("OCSPResponse: %v\n", state.OCSPResponse)

	// Now use tlsConn like any net.Conn
	buf := make([]byte, 4096)
	n, err := tlsConn.Read(buf)
	if err != nil {
		fmt.Println("read:", err)
		return
	}
	request := string(buf[:n])
	fmt.Println("Got request:\n", request)

	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 16\r\n\r\nHello TLS world!\n"
	tlsConn.Write([]byte(response))
}
