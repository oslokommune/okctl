// Package cleaner knows how to clean things up
package cleaner

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/oslokommune/okctl/pkg/acmapi"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	arnpkg "github.com/oslokommune/okctl/pkg/arn"
	"github.com/oslokommune/okctl/pkg/elbv2api"
)

const (
	elbv2KubernetesClusterTag = "tag:elbv2.k8s.aws/cluster"
)

// Cleaner contains state for cleaning things up
type Cleaner struct {
	elbv2  *elbv2api.ELBv2API
	acmapi *acmapi.ACMAPI

	provider v1alpha1.CloudProvider
}

// New returns an initialised cleaner
func New(provider v1alpha1.CloudProvider) *Cleaner {
	return &Cleaner{
		elbv2:    elbv2api.New(provider),
		acmapi:   acmapi.New(provider),
		provider: provider,
	}
}

// RemoveThingsUsingCertForDomain removes things using certificate after finding cert
// for domain
func (c *Cleaner) RemoveThingsUsingCertForDomain(domain string) error {
	certificateARN, err := c.acmapi.CertificateARNForDomain(domain)
	if err != nil {
		if errors.Is(err, acmapi.ErrNotFound) {
			return nil
		}

		return err
	}

	return c.RemoveThingsThatAreUsingCertificate(certificateARN)
}

// RemoveThingsThatAreUsingCertificate removes usages of a certificate
func (c *Cleaner) RemoveThingsThatAreUsingCertificate(certificateARN string) error {
	arns, err := c.acmapi.InUseBy(certificateARN)
	if err != nil {
		return err
	}

	for _, arn := range arns {
		isLoadBalancer, err := arnpkg.IsLoadBalancer(arn)
		if err != nil {
			return err
		}

		if isLoadBalancer {
			listeners, err := c.elbv2.GetListenersForLoadBalancer(arn)
			if err != nil {
				return err
			}

			err = c.elbv2.DeleteListenersWithCertificate(certificateARN, listeners)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteDanglingALBs will delete any remaining ALBs in a vpc
func (c *Cleaner) DeleteDanglingALBs(vpcID string) error {
	balancers, err := c.provider.ELBV2().DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return fmt.Errorf("listing load balancers: %w", err)
	}

	loadBalancers := balancers.LoadBalancers

	for _, balancer := range loadBalancers {
		if vpcID == *balancer.VpcId {
			arn := *balancer.LoadBalancerArn

			_, err := c.provider.ELBV2().DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
				LoadBalancerArn: &arn,
			})
			if err != nil {
				return fmt.Errorf("deleting load balancer: %w", err)
			}
		}
	}

	return nil
}

// DeleteDanglingSecurityGroups will remove any remaining security groups in a vpc
func (c *Cleaner) DeleteDanglingSecurityGroups(vpcID string) error {
	groups, err := c.provider.EC2().DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcID),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("listing security groups for vpc: %w", err)
	}

	for _, group := range groups.SecurityGroups {
		if *group.GroupName != "default" {
			_, err = c.provider.EC2().DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
				GroupId: group.GroupId,
			})
			if err != nil {
				return fmt.Errorf("deleting security group: %w", err)
			}
		}
	}

	return nil
}

// DeleteDanglingTargetGroups deletes dangling target groups in vpc
func (c *Cleaner) DeleteDanglingTargetGroups(clusterName string) error {
	var marker *string = nil

	var all []*elbv2.TargetGroup

	for {
		groups, err := c.provider.ELBV2().DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
			Marker: marker,
		})
		if err != nil {
			return fmt.Errorf("listing target groups: %w", err)
		}

		all = append(all, groups.TargetGroups...)

		if groups.NextMarker == nil {
			break
		}

		marker = groups.NextMarker
	}

	remove := map[string]struct{}{}

NextTargetGroup:
	for _, tg := range all {
		tags, err := c.provider.ELBV2().DescribeTags(&elbv2.DescribeTagsInput{
			ResourceArns: []*string{
				tg.TargetGroupArn,
			},
		})
		if err != nil {
			return fmt.Errorf("describing tags for target group: %w", err)
		}

		for _, desc := range tags.TagDescriptions {
			for _, tag := range desc.Tags {
				if *tag.Key == elbv2KubernetesClusterTag && *tag.Value == clusterName {
					remove[*tg.TargetGroupArn] = struct{}{}
					continue NextTargetGroup
				}
			}
		}
	}

	for targetGroupARN := range remove {
		_, err := c.provider.ELBV2().DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
			TargetGroupArn: aws.String(targetGroupARN),
		})
		if err != nil {
			return fmt.Errorf("removing target group: %w", err)
		}
	}

	return nil
}
