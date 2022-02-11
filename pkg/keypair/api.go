// Package keypair provides cryptographic key pairs using the ed25519 format
package keypair

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/keypair/edkey"
	"golang.org/x/crypto/ssh"
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
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return nil, fmt.Errorf("generating ed25519 keys: %w", err)
	}

	pubSSH, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("creating SSH public key: %w", err)
	}

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(priv),
	}
	privateKey := pem.EncodeToMemory(pemKey)
	authorizedKey := ssh.MarshalAuthorizedKey(pubSSH)

	return &Keypair{
		PublicKey:  authorizedKey,
		PrivateKey: privateKey,
	}, nil
}
