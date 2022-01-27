package keypair_test

import (
	"math/rand"
	"testing"

	"github.com/oslokommune/okctl/pkg/keypair"
	"github.com/sebdah/goldie/v2"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Run("should create expected keypair", func(t *testing.T) {
		// Given
		r := rand.New(rand.NewSource(234149819819191918)) //nolint:gosec

		// When
		keyPair, err := keypair.GenerateFromReader(r)

		// Then
		assert.NoError(t, err)

		g := goldie.New(t)
		g.Assert(t, "publicKey", keyPair.PublicKey)
		g.Assert(t, "privateKey", keyPair.PrivateKey)
	})
}
