package subnet_test

import (
	"net"
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/subnet"
	"github.com/stretchr/testify/assert"
)

func CidrFromString(t *testing.T, block string) *net.IPNet {
	_, network, err := net.ParseCIDR(block)
	assert.NoError(t, err)

	return network
}

func TestSubnets(t *testing.T) {
	testCases := []struct {
		name      string
		cidr      *net.IPNet
		num       int
		prefix    int
		creator   subnet.CreatorFn
		expect    interface{}
		expectErr bool
	}{
		{
			name:    "Validate range",
			cidr:    CidrFromString(t, "192.168.0.0/20"),
			num:     3,
			prefix:  24,
			creator: subnet.NoopCreator(),
			expect:  3,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := subnet.New(tc.num, tc.prefix, tc.cidr, tc.creator)
			if tc.expectErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
