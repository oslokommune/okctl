package saml_test

import (
	"io/ioutil"
	"testing"

	"github.com/oslokommune/okctl/pkg/credentials/aws/saml"
	"github.com/stretchr/testify/assert"
)

// nolint: lll
func TestVerifyAssertion(t *testing.T) {
	testCases := []struct {
		name      string
		role      string
		data      []byte
		expect    interface{}
		expectErr bool
	}{
		{
			name: "Should work",
			role: "arn:aws:iam::000000000000:role/oslokommune/iamadmin-SAML",
			data: func() []byte {
				data, err := ioutil.ReadFile("testdata/samlAssertion.blob")
				assert.NoError(t, err)

				return data
			}(),
			expect:    nil,
			expectErr: false,
		},
		{
			name:      "Should fail decoding",
			role:      "",
			data:      []byte("garbage"),
			expect:    "an error occurred when processing the SAML response from AWS:\nbase64 decoding SAML assertion: illegal base64 data at input byte 4\n",
			expectErr: true,
		},
		{
			name: "Should fail finding role",
			role: " arn:aws:iam::1234567890:role/oslokommune/iamadmin-SAM",
			data: func() []byte {
				data, err := ioutil.ReadFile("testdata/samlAssertion.blob")
				assert.NoError(t, err)

				return data
			}(),
			expect:    "an error occurred when processing the SAML response from AWS:\nyou do not have permission to use the role:  arn:aws:iam::1234567890:role/oslokommune/iamadmin-SAM, ask for help in #kjørermiljø-support on slack\n",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := saml.VerifyAssertion(tc.role, tc.data)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
