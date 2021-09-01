// Package keypair implements functionality for generating an RSA
// keypair. This is based on:
// - https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673
package keypair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"io"

	"golang.org/x/crypto/ssh"
)

// DefaultBitSize sets the d
const DefaultBitSize = 4096

// DefaultRandReader returns the default random generator
func DefaultRandReader() io.Reader {
	return rand.Reader
}

// Keypair contains the state for generating the keypair
type Keypair struct {
	BitSize    int
	RandReader io.Reader
	PublicKey  []byte
	PrivateKey []byte
}

// New returns an initialise keypair generator
func New(randReader io.Reader, bitSize int) *Keypair {
	return &Keypair{
		RandReader: randReader,
		BitSize:    bitSize,
	}
}

// Generate the keypair
func (k *Keypair) Generate() (*Keypair, error) {
	if k.PublicKey != nil && k.PrivateKey != nil {
		return k, nil
	}

	privateKey, err := k.generatePrivateKey()
	if err != nil {
		return nil, fmt.Errorf(constant.GeneratePrivateKeyError, err)
	}

	publicKeyBytes, err := k.generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf(constant.GeneratePublicKeyError, err)
	}

	privateKeyBytes := k.encodePrivateKeyAsPEM(privateKey)

	k.PrivateKey = privateKeyBytes
	k.PublicKey = publicKeyBytes

	return k, nil
}

// generatePrivateKey creates a RSA Private Key of specified size
func (k *Keypair) generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, k.BitSize)
	if err != nil {
		return nil, fmt.Errorf(constant.GeneratePrivateKeyError, err)
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, fmt.Errorf(constant.ValidatePrivateKeyError, err)
	}

	return privateKey, nil
}

// encodePrivateKeyAsPEM encode private key in PEM format
func (k *Keypair) encodePrivateKeyAsPEM(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func (k *Keypair) generatePublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf(constant.CreateSshRsaPublicKeyError, err)
	}

	return ssh.MarshalAuthorizedKey(publicRsaKey), nil
}
