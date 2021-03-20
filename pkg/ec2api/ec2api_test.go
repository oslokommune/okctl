package ec2api_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/ec2api"
	"github.com/stretchr/testify/assert"
)

func TestEC2APIAuthorizePodToNodeGroupTraffic(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		expect    interface{}
		expectErr bool
	}{
		{
			name: "Should work",
			provider: mock.NewCloudProvider().
				DescribeSecurityGroupsResponse(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						{
							GroupId: aws.String("someGroupID"),
						},
					},
				}, nil).
				AuthorizeSecurityGroupIngressResponse(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil),
		},
		{
			name: "Fails describing security group",
			provider: mock.NewCloudProvider().
				DescribeSecurityGroupsResponse(nil, fmt.Errorf("something bad")),
			expect:    "getting security group for node: something bad",
			expectErr: true,
		},
		{
			name: "Fails authorizing",
			provider: mock.NewCloudProvider().
				DescribeSecurityGroupsResponse(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						{
							GroupId: aws.String("someGroupID"),
						},
					},
				}, nil).
				AuthorizeSecurityGroupIngressResponse(nil, fmt.Errorf("something bad")),
			expect:    "authorizing security group ingress: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := ec2api.New(tc.provider).AuthorizePodToNodeGroupTraffic("nodegroupName", "podSecurityGroup", "vpcID")

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEC2APIRevokePodToNodeGroupTraffic(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		expect    interface{}
		expectErr bool
	}{
		{
			name: "Should work",
			provider: mock.NewCloudProvider().
				DescribeSecurityGroupsResponse(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						{
							GroupId: aws.String("someGroupID"),
						},
					},
				}, nil).
				RevokeSecurityGroupIngressResponse(&ec2.RevokeSecurityGroupIngressOutput{}, nil),
		},
		{
			name: "Fails describing security group",
			provider: mock.NewCloudProvider().
				DescribeSecurityGroupsResponse(nil, fmt.Errorf("something bad")),
			expect:    "getting security group for node: something bad",
			expectErr: true,
		},
		{
			name: "Fails authorizing",
			provider: mock.NewCloudProvider().
				DescribeSecurityGroupsResponse(&ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						{
							GroupId: aws.String("someGroupID"),
						},
					},
				}, nil).
				RevokeSecurityGroupIngressResponse(nil, fmt.Errorf("something bad")),
			expect:    "revoking security group ingress: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := ec2api.New(tc.provider).RevokePodToNodeGroupTraffic("nodegroupName", "podSecurityGroup", "vpcID")

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
