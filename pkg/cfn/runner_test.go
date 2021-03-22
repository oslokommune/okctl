package cfn_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name        string
		runner      *cfn.Runner
		expect      interface{}
		expectError bool
	}{
		{
			name: "Should work",
			runner: cfn.NewRunner(
				mock.NewCloudProvider().
					DescribeStacksEmpty().
					CreateStackSuccess().
					DescribeStacksResponse(cloudformation.StackStatusCreateComplete),
			),
		},
		{
			name: "Should fail",
			runner: cfn.NewRunner(
				mock.NewCloudProvider().
					DescribeStacksEmpty().
					CreateStackSuccess().
					DescribeStacksResponse(cloudformation.StackStatusCreateFailed).
					DescribeStackEventsSuccess(),
			),
			expect:      `{"detail":{"failed events":["ec2: something went wrong"]},"type":"internal error","error":"stack failed","code":7}`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.runner.CreateIfNotExists(mock.DefaultStackName, []byte{}, nil, 10)
			if tc.expectError {
				result, _ := json.Marshal(err)

				assert.Error(t, err)
				assert.Equal(t, tc.expect, string(result))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name        string
		runner      *cfn.Runner
		expect      interface{}
		expectError bool
	}{
		{
			name: "Should work",
			runner: cfn.NewRunner(
				mock.NewCloudProvider().
					DeleteStackSuccess().
					DescribeStacksResponse(cloudformation.StackStatusDeleteComplete),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.runner.Delete(mock.DefaultStackName)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
