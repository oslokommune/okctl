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
		r := rand.New(rand.NewSource(215823169816)) //nolint:gosec

		// When
		keyPair, err := keypair.GenerateFromReader(r)

		// Then
		assert.NoError(t, err)

		g := goldie.New(t)
		g.Assert(t, "publicKey", keyPair.PublicKey)
		g.Assert(t, "privateKey", keyPair.PrivateKey)

		// Another way to validate these keys:
		// ssh-add privateKey.golden
		// ssh-keygen -l -f publicKey.golden
	})
}
