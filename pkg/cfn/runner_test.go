package cfn_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

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
			expect:      "stack: something, failed events: ec2: something went wrong",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.runner.CreateIfNotExists("myCluster", mock.DefaultStackName, []byte{}, nil, 10)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
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

func TestGet(t *testing.T) {
	testCases := []struct {
		name        string
		runner      *cfn.Runner
		expectError bool
		expect      interface{}
	}{
		{
			name: "Should return a stack on existing cfn template",
			runner: cfn.NewRunner(
				mock.NewGoodCloudProvider().
					DescribeStacksResponse("OK"),
			),
			expect: cloudformation.Stack{
				StackId:           aws.String("myStack"),
				StackName:         aws.String("myStack"),
				StackStatus:       aws.String("OK"),
				StackStatusReason: aws.String("something"),
			},
		},
		{
			name: "Should return an error on missing stack",
			runner: cfn.NewRunner(
				mock.NewBadCloudProvider().
					DescribeStacksError(
						awserr.New(
							"ValidationError",
							"Stack with id stackname does not exist",
							nil,
						),
					),
			),
			expectError: true,
			expect:      cloudformation.Stack{},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			stack, err := tc.runner.Get("stackname")
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expect, stack)
		})
	}
}
