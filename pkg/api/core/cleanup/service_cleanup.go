// Package cleanup provides functionality to clean up AWS resources not managed by us
// For example when deleting VPC
package cleanup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// DeleteDanglingALBs will delete any remaining ALBs in a vpc
func DeleteDanglingALBs(provider v1alpha1.CloudProvider, vpcID string) error {
	balancers, err := provider.ELBV2().DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return err
	}

	loadBalancers := balancers.LoadBalancers

	for _, balancer := range loadBalancers {
		if vpcID == *balancer.VpcId {
			arn := *balancer.LoadBalancerArn

			_, err := provider.ELBV2().DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
				LoadBalancerArn: &arn,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteDanglingSecurityGroups will remove any remaining security groups in a vpc
func DeleteDanglingSecurityGroups(provider v1alpha1.CloudProvider, vpcID string) error {
	groups, err := provider.EC2().DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcID),
				},
			},
			{
				Name: aws.String(fmt.Sprintf("tag:%s", v1alpha1.OkctlManagedTag)),
				Values: []*string{
					aws.String("true"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	for _, group := range groups.SecurityGroups {
		if *group.GroupName != "default" {
			_, err = provider.EC2().DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
				GroupId: group.GroupId,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
