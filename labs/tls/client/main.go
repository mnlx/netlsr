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
	"math/big"
	"net"
)

var (
	dhP = big.NewInt(23)
	dhG = big.NewInt(5)
)

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

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// 1. Receive server's RSA public key (as base64'd PEM)
	pemB64, _ := reader.ReadString('\n')
	pemB64 = pemB64[:len(pemB64)-1]
	pemBytes, _ := base64.StdEncoding.DecodeString(pemB64)
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		panic("Failed to decode PEM block")
	}
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	serverRSAPub := pubAny.(*rsa.PublicKey)

	// 2. DH key exchange
	clientDHPriv := randomSecret()
	clientDHPub := dhPublicKey(clientDHPriv)
	// Send to server
	fmt.Fprintf(conn, "%s\n", clientDHPub.String())

	// Receive server DH public key
	serverDHPubStr, _ := reader.ReadString('\n')
	serverDHPub := new(big.Int)
	serverDHPub.SetString(serverDHPubStr[:len(serverDHPubStr)-1], 10)

	// Receive signature
	sigB64, _ := reader.ReadString('\n')
	sigB64 = sigB64[:len(sigB64)-1]
	sig, _ := base64.StdEncoding.DecodeString(sigB64)

	// Verify the signature over serverDHPub
	dhBytes := serverDHPub.Bytes()
	dhHash := sha256.Sum256(dhBytes)
	err = rsa.VerifyPSS(serverRSAPub, crypto.SHA256, dhHash[:], sig, nil)
	if err != nil {
		panic("Signature verification failed! You are talking to an impostor!")
	} else {
		fmt.Println("Server identity verified by signature!")
	}

	// Compute shared secret
	secret := sharedSecret(serverDHPub, clientDHPriv)
	fmt.Printf("Shared secret (client): %s\n", secret.String())

	// HTTP (unencrypted for this example)
	fmt.Fprintf(conn, "GET / HTTP/1.1\r\n\r\n")
	resp, _ := reader.ReadString('\n')
	fmt.Printf("HTTP response: %s", resp)
}
