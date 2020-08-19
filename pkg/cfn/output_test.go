package cfn_test

import (
	"encoding/base64"
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/stretchr/testify/assert"
)

const (
	// DefaultSubnetID is a valid, but fake, AWS subnet id
	DefaultSubnetID = "subnet-0bb1c79de3EXAMPLE"
	// DefaultSubnetCIDR is a valid CIDR
	DefaultSubnetCIDR = "192.168.0.0/24"
	// DefaultAvailabilityZone is a valid AWS availability zone
	DefaultAvailabilityZone = "eu-west-1a"
)

func TestSubnets(t *testing.T) {
	testCases := []struct {
		name        string
		provider    v1alpha1.CloudProvider
		value       string
		expect      interface{}
		expectError bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			value:    DefaultSubnetID,
			expect: []api.VpcSubnet{
				{
					ID:               DefaultSubnetID,
					Cidr:             DefaultSubnetCIDR,
					AvailabilityZone: DefaultAvailabilityZone,
				},
			},
		},
		{
			name:        "Should fail",
			provider:    mock.NewBadCloudProvider(),
			value:       DefaultSubnetID,
			expect:      "failed to describe subnet outputs: something bad",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var got []api.VpcSubnet

			err := cfn.Subnets(tc.provider, &got)(tc.value)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name  string
		value string
		input string
	}{
		{
			name:  "Should work",
			input: "hi there",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := cfn.String(&tc.value)(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, tc.value)
		})
	}
}

func TestOutput(t *testing.T) {
	testCases := []struct {
		name     string
		outputer cfn.StackOutputer
		expect   map[string]map[string]interface{}
	}{
		{
			name:     "Joined",
			outputer: cfn.NewJoined("JoinedTest").Add("value"),
			expect: map[string]map[string]interface{}{
				"JoinedTest": {
					"Value": base64.StdEncoding.EncodeToString([]byte("{ \"Fn::Join\": [ \",\", [ \"value\" ] ] }")),
				},
			},
		},
		{
			name:     "Value",
			outputer: cfn.NewValue("ValueTest", "value"),
			expect: map[string]map[string]interface{}{
				"ValueTest": {
					"Value": "value",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, tc.outputer.NamedOutputs())
		})
	}
}
