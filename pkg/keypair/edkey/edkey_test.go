package edkey_test

import (
	"bytes"
	"crypto/ed25519"
	"github.com/oslokommune/okctl/pkg/keypair/edkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

// TestEdkeyContainsOriginalKey verifies that package edkey wraps a ed25519 key
func TestEdkeyContainsOriginalKey(t *testing.T) {
	t.Run("edkey should produce a OpenSSH key that contains a ED25519 key", func(t *testing.T) {
		r := rand.New(rand.NewSource(13656478913)) //nolint:gosec

		_, priv, err := ed25519.GenerateKey(r)
		require.NoError(t, err)

		privBytes := []byte(priv)

		privOpenSSH := edkey.MarshalED25519PrivateKey(priv)
		assert.True(t, bytes.Contains(privOpenSSH, privBytes))
	})
}
