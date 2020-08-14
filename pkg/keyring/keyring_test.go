package keyring_test

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/keyring"
	"github.com/stretchr/testify/assert"
	"testing"
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
			keytype: keyring.KeyTypeUserPassword,
			secret:  "foobar",
			expect:  true,
		},
		{
			name:      "Empty password should fail",
			keytype:   keyring.KeyTypeUserPassword,
			secret:    "",
			expectErr: true,
			expect:    fmt.Errorf("key of type userPassword cannot store empty value"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ring, err := keyring.New(keyring.NewInMemoryKeyring())
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
		ring	  keyring.Keyringer
		keytype   keyring.KeyType
		secret    string
		expect    string
		expectErr bool
	}{
		{
			name: "Fetched password should work",
			ring: func() keyring.Keyringer {
				r, err := keyring.New(keyring.NewInMemoryKeyring())
				assert.NoError(t, err)

				return r
			}(),
			keytype: keyring.KeyTypeUserPassword,
			secret : "s3cret",
			expect:  "s3cret",
		},
		{
			name: "Fetching non exisiting secret should return error",
			ring: func() keyring.Keyringer {
				r, err := keyring.New(keyring.NewInMemoryKeyring())
				assert.NoError(t, err)

				return r
			}(),
			keytype:   keyring.KeyTypeUserPassword,
			expectErr: true,
			expect:    "The specified item could not be found in the keyring",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.ring.Store(keyring.KeyTypeUserPassword, tc.expect)
			got, err := tc.ring.Fetch(tc.keytype)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, got)
		})
	}
}
