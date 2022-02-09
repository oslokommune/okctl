// Package keypair provides cryptographic key pairs using the ed25519 format
package keypair

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"github.com/ScaleFT/sshkeys"
	"golang.org/x/crypto/ssh"
	"io"
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

	privateKey, err := privateKeyToOpenSSHFormat(privateKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("converting private key to PEM format: %w", err)
	}

	publicKey, err := toSSHKeyFormat(publicKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("converting public key to PEM format: %w", err)
	}

	return &Keypair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func toSSHKeyFormat(publicKey ed25519.PublicKey) ([]byte, error) {
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create ssh-ed25519 public key: %w", err)
	}

	return ssh.MarshalAuthorizedKey(sshPublicKey), nil
}

func privateKeyToOpenSSHFormat(privateKey ed25519.PrivateKey) ([]byte, error) {
	// When Golang adds support for ED25519 keys in OpenSSH format, replace the sshkeys library dependency.
	// See: https://github.com/golang/go/issues/37132
	// Currently using: https://github.com/ScaleFT/sshkeys
	privateKeyOpenSSHFormat, err := sshkeys.Marshal(privateKey, &sshkeys.MarshalOptions{
		Format: sshkeys.FormatOpenSSHv1,
	})
	if err != nil {
		return nil, fmt.Errorf("marshalling private key: %w", err)
	}

	return privateKeyOpenSSHFormat, nil
}
