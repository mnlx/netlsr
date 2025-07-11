package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"crypto/aes"
	"crypto/cipher"
)

var (
	// Tiny numbers for demo purposes only!
	dhP        = big.NewInt(23) // prime modulus
	dhG        = big.NewInt(5)  // generator
	serverPriv *rsa.PrivateKey
	serverPub  *rsa.PublicKey
)

func secretToAESKey(secret *big.Int) []byte {
	b := secret.Bytes()
	key := make([]byte, 16)
	copy(key, b)
	return key
}

// Generate a random secret < p
func randomSecret() *big.Int {
	secret, _ := rand.Int(rand.Reader, dhP)
	fmt.Println(dhP, dhG, secret)
	return secret
}

// Compute (g^a mod p)
func publicKey(secret *big.Int) *big.Int {
	pk := new(big.Int).Exp(dhG, secret, dhP)
	fmt.Println(pk)
	return pk
}

// Compute (otherPub^mySecret mod p)
func sharedSecret(otherPub, mySecret *big.Int) *big.Int {
	s := new(big.Int).Exp(otherPub, mySecret, dhP)
	fmt.Printf("Shared key: %s %s %s", otherPub, mySecret, dhP)
	return s
}

func init() {
	var err error
	// Generate a 2048-bit RSA keypair
	serverPriv, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	serverPub = &serverPriv.PublicKey
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
			fmt.Println("accept error:", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Send the server's public key to the client (PEM-encoded)
	pubASN1, _ := x509.MarshalPKIXPublicKey(serverPub)
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	conn.Write(pubPEM)
	conn.Write([]byte("\n")) // Delimiter

	// 1. Server chooses secret, computes public key
	serverSecret := randomSecret()
	serverDHPub := dhPublicKey(serverSecret)

	fmt.Printf("Server Pub String: %s", serverPub)

	// Sign the DH public value with server's RSA key
	dhBytes := serverPub.Bytes()
	sig, err := rsa.SignPSS(rand.Reader, serverPriv, crypto.SHA256, sha256.Sum256(dhBytes)[:], nil)
	if err != nil {
		panic(err)
	}

	// After sending the DH public value, also send the signature (base64 for clarity)
	fmt.Fprintf(conn, "%s\n", serverPub.String())
	conn.Write(sig)
	conn.Write([]byte("\n")) // Delimiter

	// 2. Receive client public key
	clientPubStr, err := reader.ReadString('\n')
	fmt.Printf("Client Pub String: %s", clientPubStr)
	if err != nil {
		fmt.Println("failed to read client pubkey:", err)
		return
	}
	clientPub := new(big.Int)
	clientPub.SetString(clientPubStr[:len(clientPubStr)-1], 10)

	// 3. Send our public key to client
	fmt.Fprintf(conn, "%s\n", serverPub.String())

	// 4. Compute shared secret
	secret := sharedSecret(clientPub, serverSecret)
	fmt.Printf("Shared secret (server): %s\n", secret.String())

	// === AES-CTR key/iv setup ===
	key := secretToAESKey(secret)

	// Server creates random IV and sends to client
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}

	conn.Write(iv)

	// Set up decrypt/encrypt streams
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	streamReader := &cipher.StreamReader{S: stream, R: reader}
	streamWriter := &cipher.StreamWriter{S: stream, W: conn}

	// Read encrypted HTTP request
	httpLine, err := bufio.NewReader(streamReader).ReadString('\n')
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	fmt.Printf("Decrypted HTTP request: %s", httpLine)

	// Write encrypted HTTP response
	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, Encrypted World!\n"
	streamWriter.Write([]byte(response))
}
