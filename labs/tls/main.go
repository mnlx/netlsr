package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
)

var (
	dhP           = big.NewInt(23) // Small prime for demo only!
	dhG           = big.NewInt(5)
	serverRSAPriv *rsa.PrivateKey
	serverRSAPub  *rsa.PublicKey
)

func init() {
	var err error
	serverRSAPriv, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	serverRSAPub = &serverRSAPriv.PublicKey
}

func randomSecret() *big.Int {
	secret, _ := rand.Int(rand.Reader, dhP)
	return secret
}

func dhPublicKey(secret *big.Int) *big.Int {
	pk := new(big.Int).Exp(dhG, secret, dhP)
	return pk
}

func sharedSecret(otherPub, mySecret *big.Int) *big.Int {
	s := new(big.Int).Exp(otherPub, mySecret, dhP)
	return s
}

func handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// 1. Send server's RSA public key (PEM-encoded, single-line for easy parsing)
	pubASN1, _ := x509.MarshalPKIXPublicKey(serverRSAPub)
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	// Send public key (single-line, base64 to avoid multi-line PEM issues)
	pemB64 := base64.StdEncoding.EncodeToString(pubPEM)
	conn.Write([]byte(pemB64 + "\n"))

	// 2. DH key exchange
	serverDHPriv := randomSecret()
	serverDHPub := dhPublicKey(serverDHPriv)

	// Receive client DH public key
	clientDHPubStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read client DH pubkey:", err)
		return
	}
	clientDHPub := new(big.Int)
	clientDHPub.SetString(clientDHPubStr[:len(clientDHPubStr)-1], 10)

	// 3. Send server DH public key
	conn.Write([]byte(serverDHPub.String() + "\n"))

	// 4. Sign DH public value and send signature
	dhBytes := serverDHPub.Bytes()
	dhHash := sha256.Sum256(dhBytes)
	sig, err := rsa.SignPSS(rand.Reader, serverRSAPriv, crypto.SHA256, dhHash[:], nil)
	if err != nil {
		fmt.Println("Failed to sign:", err)
		return
	}
	sigB64 := base64.StdEncoding.EncodeToString(sig)
	conn.Write([]byte(sigB64 + "\n"))

	// 5. Compute shared secret (normally would use this as a session key)
	secret := sharedSecret(clientDHPub, serverDHPriv)
	fmt.Printf("Shared secret (server): %s\n", secret.String())

	// 6. Basic HTTP (unencrypted, for simplicity in this example)
	httpLine, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Println("HTTP read error:", err)
		return
	}
	fmt.Printf("Received HTTP: %s", httpLine)

	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, Authenticated World!\n"
	conn.Write([]byte(response))
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on :8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handle(conn)
	}
}
