package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net"
)

var (
	// Tiny numbers for demo purposes only!
	dhP = big.NewInt(23) // prime modulus
	dhG = big.NewInt(5)  // generator
)

func secretToAESKey(secret *big.Int) []byte {
	b := secret.Bytes()
	key := make([]byte, 16)
	copy(key, b)
	return key
}

func randomSecret() *big.Int {
	secret, _ := rand.Int(rand.Reader, dhP)
	return secret
}

func publicKey(secret *big.Int) *big.Int {
	return new(big.Int).Exp(dhG, secret, dhP)
}

func sharedSecret(otherPub, mySecret *big.Int) *big.Int {
	return new(big.Int).Exp(otherPub, mySecret, dhP)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	pubPEM, _ := reader.ReadString('\n')
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		panic("failed to decode PEM block containing public key")
	}
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	serverPub := pubAny.(*rsa.PublicKey)

	// After receiving server's DH public value:
	serverDHPubStr, _ := reader.ReadString('\n')
	serverDHPub := new(big.Int)
	serverDHPub.SetString(serverDHPubStr[:len(serverDHPubStr)-1], 10)

	// Receive signature
	sig, _ := reader.ReadString('\n')
	sigBytes, _ := base64.StdEncoding.DecodeString(sig[:len(sig)-1])

	// Verify the signature over the DH public value
	dhBytes := serverDHPub.Bytes()
	err = rsa.VerifyPSS(serverPub, crypto.SHA256, sha256.Sum256(dhBytes)[:], sigBytes, nil)
	if err != nil {
		panic("Signature verification failed! You are talking to an impostor!")
	} else {
		fmt.Println("Server identity verified by signature!")
	}

	// Diffie-Hellman as before...
	clientSecret := randomSecret()
	clientPub := publicKey(clientSecret)
	fmt.Fprintf(conn, "%s\n", clientPub.String())
	serverPubStr, _ := reader.ReadString('\n')
	serverPub = new(big.Int)
	serverPub.SetString(serverPubStr[:len(serverPubStr)-1], 10)
	secret := sharedSecret(serverPub, clientSecret)
	fmt.Printf("Shared secret (client): %s\n", secret.String())

	// === AES-CTR key/iv setup ===
	key := secretToAESKey(secret)

	// Receive IV from server
	iv := make([]byte, aes.BlockSize)
	io.ReadFull(reader, iv)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	streamReader := &cipher.StreamReader{S: stream, R: reader}
	streamWriter := &cipher.StreamWriter{S: stream, W: conn}

	// Write encrypted HTTP request
	streamWriter.Write([]byte("GET / HTTP/1.1\r\n\r\n"))

	// Read encrypted HTTP response
	resp, _ := bufio.NewReader(streamReader).ReadSlice('8')
	fmt.Printf("Decrypted response: %s", resp)
}
