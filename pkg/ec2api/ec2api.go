// Package ec2api provides some convenience functions
// for interacting with the AWS EC2 API
package ec2api

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const (
	dnsPort = 53
)

// EC2API contains the state required for interacting
// with the AWS EC2 API
type EC2API struct {
	provider v1alpha1.CloudProvider
}

// New returns an initialised AWS EC2 API client
func New(provider v1alpha1.CloudProvider) *EC2API {
	return &EC2API{
		provider: provider,
	}
}

// CLI version
// aws ec2 describe-security-groups \
// --filter Name=tag:alpha.eksctl.io/nodegroup-name,Values=ng-generic-1-20-1a
// --filter Name=vpc-id,Values=vpc-06a32841508c1721b \
// | jq '.SecurityGroups[0].GroupId'
func (a *EC2API) securityGroupForNodeGroup(name, vpcID string) (string, error) {
	sgs, err := a.provider.EC2().DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:alpha.eksctl.io/nodegroup-name"),
				Values: aws.StringSlice([]string{name}),
			},
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{vpcID}),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("calling EC2 API: %w", err)
	}

	if len(sgs.SecurityGroups) == 0 {
		return "", ErrNotFound
	}

	return *sgs.SecurityGroups[0].GroupId, nil
}

func (a *EC2API) permissionsForPodToNodeGroup(podSecurityGroup, vpcID string) []*ec2.IpPermission {
	protocols := []string{"tcp", "udp"}

	permissions := make([]*ec2.IpPermission, len(protocols))

	for i, protocol := range protocols {
		permissions[i] = &ec2.IpPermission{
			FromPort:   aws.Int64(dnsPort),
			IpProtocol: aws.String(protocol),
			ToPort:     aws.Int64(dnsPort),
			UserIdGroupPairs: []*ec2.UserIdGroupPair{
				{
					GroupId: aws.String(podSecurityGroup),
					VpcId:   aws.String(vpcID),
				},
			},
		}
	}

	return permissions
}

// AuthorizePodToNodeGroupTraffic adds ingress rules that allows the pod
// to communicate with the node
func (a *EC2API) AuthorizePodToNodeGroupTraffic(nodegroupName, podSecurityGroup, vpcID string) error {
	nodegroupSecurityGroup, err := a.securityGroupForNodeGroup(nodegroupName, vpcID)
	if err != nil {
		return fmt.Errorf("finding security group for node group: %w", err)
	}

	_, err = a.provider.EC2().AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       aws.String(nodegroupSecurityGroup),
		IpPermissions: a.permissionsForPodToNodeGroup(podSecurityGroup, vpcID),
	})
	if err != nil {
		if strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
			return nil
		}

		return fmt.Errorf("authorizing security group ingress: %w", err)
	}

	return nil
}

// RevokePodToNodeGroupTraffic removes communications
func (a *EC2API) RevokePodToNodeGroupTraffic(nodegroupName, podSecurityGroup, vpcID string) error {
	nodegroupSecurityGroup, err := a.securityGroupForNodeGroup(nodegroupName, vpcID)
	if err != nil {
		return fmt.Errorf("finding security group for node group: %w", err)
	}

	_, err = a.provider.EC2().RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
		GroupId:       aws.String(nodegroupSecurityGroup),
		IpPermissions: a.permissionsForPodToNodeGroup(podSecurityGroup, vpcID),
	})
	if err != nil {
		if strings.Contains(err.Error(), "InvalidPermission.NotFound") {
			return nil
		}

		return fmt.Errorf("revoking security group ingress: %w", err)
	}

	return nil
}
