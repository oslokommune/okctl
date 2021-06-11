package servicequota_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/servicequota"
)

// nolint: funlen
func TestCheckQuotas(t *testing.T) {
	testCases := []struct {
		name      string
		checks    []servicequota.Checker
		expect    interface{}
		expectErr bool
	}{
		{
			name: "EIP check valid",
			checks: []servicequota.Checker{
				servicequota.NewEipCheck(1, mock.NewGoodCloudProvider()),
			},
		},
		{
			name: "EIP check no more available",
			checks: []servicequota.Checker{
				servicequota.NewEipCheck(5, mock.NewGoodCloudProvider()),
			},
			expect:    "AWS VPC Elastic IPs: required 5, but only have 2 available",
			expectErr: true,
		},
		{
			name: "VPC check valid",
			checks: []servicequota.Checker{
				servicequota.NewVpcCheck(1, mock.NewGoodCloudProvider()),
			},
		},
		{
			name: "VPC check no more available",
			checks: []servicequota.Checker{
				servicequota.NewVpcCheck(5, mock.NewGoodCloudProvider()),
			},
			expect:    "AWS VPCs: required 5, but only have 2 available",
			expectErr: true,
		},
		{
			name: "IGW check valid",
			checks: []servicequota.Checker{
				servicequota.NewIgwCheck(1, mock.NewGoodCloudProvider()),
			},
		},
		{
			name: "IGW check no more available",
			checks: []servicequota.Checker{
				servicequota.NewIgwCheck(5, mock.NewGoodCloudProvider()),
			},
			expect:    "AWS VPC Internet Gateways: required 5, but only have 2 available",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := servicequota.CheckQuotas(tc.checks...)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
