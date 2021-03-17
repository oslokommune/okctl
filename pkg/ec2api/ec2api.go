// Package ec2api provides some convenience functions
// for
package ec2api

import (
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
		return "", err
	}

	return *sgs.SecurityGroups[0].GroupId, nil
}

// AuthorizePodToNodeGroupTraffic adds ingress rules that allows the pod
// to communicate with the node
func (a *EC2API) AuthorizePodToNodeGroupTraffic(nodegroupName, podSecurityGroup, vpcID string) error {
	nodegroupSecurityGroup, err := a.securityGroupForNodeGroup(nodegroupName, vpcID)
	if err != nil {
		return err
	}

	for _, protocol := range []string{"tcp", "udp"} {
		_, err = a.provider.EC2().AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			FromPort:                aws.Int64(dnsPort),
			GroupId:                 aws.String(nodegroupSecurityGroup),
			IpProtocol:              aws.String(protocol),
			SourceSecurityGroupName: aws.String(podSecurityGroup),
			ToPort:                  aws.Int64(dnsPort),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// RevokePodToNodeGroupTraffic removes communications
func (a *EC2API) RevokePodToNodeGroupTraffic(nodegroupName, podSecurityGroup, vpcID string) error {
	nodegroupSecurityGroup, err := a.securityGroupForNodeGroup(nodegroupName, vpcID)
	if err != nil {
		return err
	}

	for _, protocol := range []string{"tcp", "udp"} {
		_, err = a.provider.EC2().RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
			FromPort:                aws.Int64(dnsPort),
			GroupId:                 aws.String(nodegroupSecurityGroup),
			IpProtocol:              aws.String(protocol),
			SourceSecurityGroupName: aws.String(podSecurityGroup),
			ToPort:                  aws.Int64(dnsPort),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
