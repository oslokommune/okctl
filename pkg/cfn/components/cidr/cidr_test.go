package cidr_test

import (
	"net"
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/cidr"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name      string
		cidr      string
		expect    interface{}
		expectErr bool
	}{
		{
			name:      "Invalid format",
			cidr:      "192.168.0.0",
			expect:    "invalid CIDR address: 192.168.0.0",
			expectErr: true,
		},
		{
			name:      "Incorrect IP version",
			cidr:      "2001:db8::/32",
			expect:    "cidr (2001:db8::/32) is not of type IPv4",
			expectErr: true,
		},
		{
			name:      "Too small address space",
			cidr:      "192.168.0.0/28",
			expect:    "address space of cidr (192.168.0.0/28) is less than required: 16 < 4096",
			expectErr: true,
		},
		{
			name:      "Not a private range",
			cidr:      "120.120.0.0/16",
			expect:    "provided cidr (120.120.0.0/16) is not in the legal ranges: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16",
			expectErr: true,
		},
		{
			name: "Validate range",
			cidr: "192.168.0.0/20",
			expect: &cidr.Cidr{
				Block: func() *net.IPNet {
					_, n, err := net.ParseCIDR("192.168.0.0/20")
					assert.NoError(t, err)
					return n
				}(),
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := cidr.New(tc.cidr, cidr.RequiredHosts(9, 24), cidr.PrivateCidrRanges())
			if tc.expectErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
