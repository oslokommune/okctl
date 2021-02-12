package servicequota

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// VpcCheck is used to check if you have enough vpcs
type VpcCheck struct {
	provider      v1alpha1.CloudProvider
	required      int
	isProvisioned bool
}

// NewVpcCheck make a new instance of check for VPCs
func NewVpcCheck(isProvisioned bool, required int, provider v1alpha1.CloudProvider) *VpcCheck {
	return &VpcCheck{
		provider:      provider,
		required:      required,
		isProvisioned: isProvisioned,
	}
}

// CheckAvailability determines if you will be able to make required Virtual Private Cloud(s)
func (v *VpcCheck) CheckAvailability() (*Result, error) {
	if v.isProvisioned {
		return &Result{
			IsProvisioned: true,
		}, nil
	}

	q, err := v.provider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-F678F1CE"),
		ServiceCode: aws.String("vpc"),
	})
	if err != nil {
		return nil, fmt.Errorf("getting vpc quota: %w", err)
	}

	vpcs, err := v.provider.EC2().DescribeVpcs(nil)
	if err != nil {
		return nil, fmt.Errorf("getting current vpc count: %w", err)
	}

	quota := int(*q.Quota.Value)
	count := len(vpcs.Vpcs)
	available := quota - count

	return &Result{
		Required:    v.required,
		Available:   available,
		HasCapacity: v.required <= available,
		Description: "AWS VPCs",
	}, nil
}
