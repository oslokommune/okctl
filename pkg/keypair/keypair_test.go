package keypair_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/oslokommune/okctl/pkg/keypair"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name       string
		bitSize    int
		randReader io.Reader
		expect     *keypair.Keypair
	}{
		{
			name:       "Should work",
			bitSize:    128,
			randReader: bytes.NewReader(nil),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			k := keypair.New(tc.randReader, tc.bitSize)
			got, err := k.Generate()
			assert.NoError(t, err)
			assert.NotNil(t, got.PublicKey)
			assert.NotNil(t, got.PrivateKey)
		})
	}
}
