package cfn

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

// Joined stores state for creating an intrinsic join
type Joined struct {
	StoredName string
	Values     []string
}

// NamedOutputs returns the named outputs
func (j *Joined) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		j.Name(): j.Outputs(),
	}
}

// Outputs returns the outputs only
func (j *Joined) Outputs() map[string]interface{} {
	return map[string]interface{}{
		"Value": cloudformation.Join(",", j.Values),
	}
}

// Name returns the name of the output
func (j *Joined) Name() string {
	return j.StoredName
}

// Add a value to the outputs that should be joined
func (j *Joined) Add(v ...string) *Joined {
	j.Values = append(j.Values, v...)

	return j
}

// NewJoined is a helper for creating cloud formation Joined data
// in the output of a cloud formation stack
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-join.html
func NewJoined(name string) *Joined {
	return &Joined{
		StoredName: name,
	}
}

// Ensure that Joined implements the StackOutputer interface
var _ StackOutputer = &Joined{}

// Value stores the state for creating an output
type Value struct {
	StoredName string
	Value      string
}

// NamedOutputs returns the named cloud formation outputs
func (v *Value) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		v.Name(): v.Outputs(),
	}
}

// Outputs returns only the cloud formation outputs
func (v *Value) Outputs() map[string]interface{} {
	return map[string]interface{}{
		"Value": v.Value,
	}
}

// Name returns the name given to the output
func (v *Value) Name() string {
	return v.StoredName
}

// NewValue is a helper for creating Value outputs in
// a cloud formation stack
func NewValue(name, v string) *Value {
	return &Value{
		StoredName: name,
		Value:      v,
	}
}

// Ensure that Value implements the StackOutputer interface
var _ StackOutputer = &Value{}

// Subnets knows how to process the output from a subnet creation
func Subnets(p v1alpha1.CloudProvider, to *[]api.VpcSubnet) ProcessOutputFn {
	return func(v string) error {
		got, err := p.EC2().DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: aws.StringSlice(strings.Split(v, ",")),
		})
		if err != nil {
			return fmt.Errorf("failed to describe subnet outputs: %w", err)
		}

		for _, s := range got.Subnets {
			*to = append(*to, api.VpcSubnet{
				ID:               *s.SubnetId,
				Cidr:             *s.CidrBlock,
				AvailabilityZone: *s.AvailabilityZone,
			})
		}

		return nil
	}
}

// String knows how to process the output from a value
func String(to *string) ProcessOutputFn {
	return func(v string) error {
		*to = v

		return nil
	}
}
