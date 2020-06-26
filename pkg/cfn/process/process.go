// Package process knows how to process cloud formation outputs
package process

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/manager"
	"github.com/pkg/errors"
)

// Subnets knows how to process the output from a subnet creation
func Subnets(p v1alpha1.CloudProvider, to map[string]v1alpha1.ClusterNetwork) manager.ProcessOutputFn {
	return func(v string) error {
		got, err := p.EC2().DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: aws.StringSlice(strings.Split(v, ",")),
		})
		if err != nil {
			return errors.Wrap(err, "failed to describe subnet outputs")
		}

		for _, s := range got.Subnets {
			to[*s.AvailabilityZone] = v1alpha1.ClusterNetwork{
				ID:   *s.SubnetId,
				CIDR: *s.CidrBlock,
			}
		}

		return nil
	}
}

// String knows how to process the output from a value
func String(to *string) manager.ProcessOutputFn {
	return func(v string) error {
		to = &v

		return nil
	}
}
