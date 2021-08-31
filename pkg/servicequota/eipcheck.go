package servicequota

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// EipCheck is used to check if you have enough Elastic Ips
type EipCheck struct {
	provider v1alpha1.CloudProvider
	required int
}

// NewEipCheck makes a new instance of check for Elastic Ips
func NewEipCheck(required int, provider v1alpha1.CloudProvider) *EipCheck {
	return &EipCheck{
		provider: provider,
		required: required,
	}
}

// CheckAvailability determines if you will be able to make required Elastic Ip(s)
func (e *EipCheck) CheckAvailability() (*Result, error) {
	q, err := e.provider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-0263D0A3"),
		ServiceCode: aws.String("ec2"),
	})
	if err != nil {
		return nil, fmt.Errorf(constant.GetEIPQuotasError, err)
	}

	eips, err := e.provider.EC2().DescribeAddresses(nil)
	if err != nil {
		return nil, fmt.Errorf(constant.GetEIPCountError, err)
	}

	quota := int(*q.Quota.Value)
	count := len(eips.Addresses)
	available := quota - count

	return &Result{
		Required:    e.required,
		Available:   available,
		HasCapacity: e.required <= available,
		Description: "AWS VPC Elastic IPs",
	}, nil
}
