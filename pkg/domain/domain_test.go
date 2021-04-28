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

// nolint: funlen
func TestHasNameServers(t *testing.T) {
	testCases := []struct {
		name      string
		domain    string
		ns        []string
		expect    interface{}
		expectErr bool
	}{
		{
			name:   "Should work with equal sets",
			domain: "test.oslo.systems",
			ns: []string{
				"ns-612.awsdns-12.net.",
				"ns-1322.awsdns-37.org.",
				"ns-1706.awsdns-21.co.uk.",
				"ns-327.awsdns-40.com.",
			},
		},
		{
			name:   "Nameservers that aren't fqdns should work",
			domain: "test.oslo.systems",
			ns: []string{
				"ns-1322.awsdns-37.org",
				"ns-612.awsdns-12.net",
				"ns-327.awsdns-40.com",
				"ns-1706.awsdns-21.co.uk",
			},
		},
		{
			name:   "Should work with partial matches",
			domain: "test.oslo.systems",
			ns: []string{
				"ns-612.awsdns-12.net.",
				"ns-1322.awsdns-37.org.",
				"fake.awsdns-21.co.uk.",
			},
		},
		{
			name:      "Should fail, with no nameservers",
			domain:    "test.oslo.systems",
			ns:        []string{},
			expectErr: true,
			expect:    "nameservers do not match, expected: [], but got: [ns-1322.awsdns-37.org. ns-1706.awsdns-21.co.uk. ns-327.awsdns-40.com. ns-612.awsdns-12.net.]",
		},
		{
			name:   "Should fail, with no matches",
			domain: "test.oslo.systems",
			ns: []string{
				"a",
				"b",
				"c",
				"d",
			},
			expectErr: true,
			expect:    "nameservers do not match, expected: [a. b. c. d.], but got: [ns-1322.awsdns-37.org. ns-1706.awsdns-21.co.uk. ns-327.awsdns-40.com. ns-612.awsdns-12.net.]",
		},
		{
			name:      "Should fail",
			domain:    "test-does-not-exist.oslo.systems",
			expect:    "unable to get NS records for domain 'test-does-not-exist.oslo.systems', does not appear to be delegated yet",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := domain.ShouldHaveNameServers(tc.domain, tc.ns)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
