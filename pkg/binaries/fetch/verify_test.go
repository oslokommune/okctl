package fetch_test

import (
	"github.com/oslokommune/okctl/pkg/binaries/digest"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestVerifer(t *testing.T) {
	testCases := []struct {
		name         string
		content      string
		expectedHash string
	}{
		{
			name:         "Should fail if invalid hash",
			content:      "hi there",
			expectedHash: "aaaa",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expectedDigests := make(map[digest.Type]string)
			expectedDigests[digest.TypeSHA256] = tc.expectedHash
			v := fetch.NewVerifier(expectedDigests)

			reader := strings.NewReader(tc.content)

			// When
			err := v.Verify(reader)

			// Then
			assert.Error(t, err)
			assert.Equal(t,
				"verification failed, hash mismatch, "+
					"got: 9b96a1fe1d548cbbc960cc6a0286668fd74a763667b06366fb2324269fcabaa4, expected: aaaa",
				err.Error(),
			)
		})
	}
}
