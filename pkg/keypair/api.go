// Package keypair provides cryptographic key pairs using the ed25519 format
package keypair

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
)

// Keypair is the key pair and its metadata
type Keypair struct {
	PublicKey  []byte
	PrivateKey []byte
}

// Generate generates a ed25519 keypair
func Generate() (*Keypair, error) {
	return generate(rand.Reader)
}

// GenerateFromReader generates a ed25519 keypair with the given reader
func GenerateFromReader(reader io.Reader) (*Keypair, error) {
	return generate(reader)
}

func generate(reader io.Reader) (*Keypair, error) {
	publicKeyRaw, privateKeyRaw, err := ed25519.GenerateKey(reader)
	if err != nil {
		return nil, fmt.Errorf("generating keypair: %w", err)
	}

	privateKey, err := privateKeyToPEMFormat(privateKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("converting private key to PEM format: %w", err)
	}

	publicKey, err := publicKeyToPEMFormat(publicKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("converting public key to PEM format: %w", err)
	}

	return &Keypair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func publicKeyToPEMFormat(publicKey ed25519.PublicKey) ([]byte, error) {
	// PKIX format explained: https://stackoverflow.com/a/49878687/915441
	privateKeyPKCS8, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("marshalling private key: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: privateKeyPKCS8,
	}

	return toPEMFormat(pemBlock)
}

func privateKeyToPEMFormat(privateKey ed25519.PrivateKey) ([]byte, error) {
	// PKCS8 format explained: https://stackoverflow.com/a/48960291/915441
	privateKeyPKCS8, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("marshalling private key: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyPKCS8,
	}

	return toPEMFormat(pemBlock)
}

func toPEMFormat(pemBlock *pem.Block) ([]byte, error) {
	var blockAsBytes bytes.Buffer

	err := pem.Encode(&blockAsBytes, pemBlock)
	if err != nil {
		return nil, fmt.Errorf("encoding to PEM format: %w", err)
	}

	b, err := ioutil.ReadAll(&blockAsBytes)
	if err != nil {
		return nil, fmt.Errorf("reading bytes: %w", err)
	}

	return b, nil
}
