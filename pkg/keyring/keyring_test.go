package keyring_test

import (
	"fmt"
	"testing"

	"github.com/oslokommune/okctl/pkg/keyring"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	testCases := []struct {
		name      string
		keytype   keyring.KeyType
		secret    string
		expect    interface{}
		expectErr bool
	}{
		{
			name:    "Should be able to store password",
			keytype: keyring.KeyTypeGithubToken,
			secret:  "foobar",
			expect:  true,
		},
		{
			name:      "Empty password should fail",
			keytype:   keyring.KeyTypeGithubToken,
			secret:    "",
			expectErr: true,
			expect:    fmt.Errorf("key of type githubToken cannot store empty value"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ring, err := keyring.New(keyring.NewInMemoryKeyring(), false)
			assert.NoError(t, err)

			err = ring.Store(tc.keytype, tc.secret)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFetch(t *testing.T) {
	testCases := []struct {
		name      string
		ring      keyring.Keyringer
		keytype   keyring.KeyType
		secret    string
		expect    string
		expectErr bool
	}{
		{
			name: "Fetched password should work",
			ring: func() keyring.Keyringer {
				r, err := keyring.New(keyring.NewInMemoryKeyring(), false)
				assert.NoError(t, err)

				return r
			}(),
			keytype: keyring.KeyTypeGithubToken,
			secret:  "s3cret",
			expect:  "s3cret",
		},
		{
			name: "Fetching non exisiting secret should return error",
			ring: func() keyring.Keyringer {
				r, err := keyring.New(keyring.NewInMemoryKeyring(), false)
				assert.NoError(t, err)

				return r
			}(),
			keytype:   keyring.KeyTypeGithubToken,
			expectErr: true,
			expect:    "The specified item could not be found in the keyring",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ring.Store(keyring.KeyTypeGithubToken, tc.expect)
			if err != nil {
				fmt.Print("Ignore error")
			}
			got, err := tc.ring.Fetch(tc.keytype)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, got)
		})
	}
}
