package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// Generate RSA private key dengan panjang 2048 bit
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Encode private key menjadi PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key menjadi PEM (SPKI)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Base64 encode agar bisa disimpan di environment variable
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyPEM)

	fmt.Println("JWT_PRIVATE_KEY:")
	fmt.Println(privateKeyBase64)

	fmt.Println("\nJWT_PUBLIC_KEY:")
	fmt.Println(publicKeyBase64)

	// Simpan juga ke file PEM untuk backup lokal
	if err := os.WriteFile("private.pem", privateKeyPEM, 0600); err != nil {
		panic(err)
	}

	if err := os.WriteFile("public.pem", publicKeyPEM, 0644); err != nil {
		panic(err)
	}

	fmt.Println("\nRSA 2048-bit key berhasil dibuat.")
}
