package domain_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestNotTaken(t *testing.T) {
	testCases := []struct {
		name        string
		fqdn        string
		expectError bool
		expect      interface{}
	}{
		{
			name: "Available domain",
			fqdn: "nosuchsubdomain.oslo.systems",
		},
		{
			name:        "Taken domain",
			fqdn:        "test.oslo.systems",
			expectError: true,
			expect:      "domain 'test.oslo.systems' already in use, found DNS records",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := domain.NotTaken(tc.fqdn)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name        string
		fqdn        string
		expectError bool
		expect      interface{}
	}{
		{
			name: "Validate domain",
			fqdn: "test.oslo.systems",
		},
		{
			name:        "Invalid domain",
			fqdn:        "not a domain.oslo.systems",
			expectError: true,
			expect:      "'not a domain.oslo.systems' is not a valid domain",
		},
		{
			name:        "Validate domain, doesn't end with oslo.systems",
			fqdn:        "some.other.domain.com",
			expectError: true,
			expect:      "'some.other.domain.com' must end with .oslo.systems",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := domain.Validate(tc.fqdn)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
