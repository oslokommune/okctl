package servicequota

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// IgwCheck is used to check if you have enough Internet Gateways
type IgwCheck struct {
	provider      v1alpha1.CloudProvider
	required      int
	isProvisioned bool
}

// NewIgwCheck makes a new instance of check for Internet Gateways
func NewIgwCheck(isProvisioned bool, required int, provider v1alpha1.CloudProvider) *IgwCheck {
	return &IgwCheck{
		provider:      provider,
		required:      required,
		isProvisioned: isProvisioned,
	}
}

// CheckAvailability determines if you will be able to make required Internet Gateway(s)
func (i *IgwCheck) CheckAvailability() (*Result, error) {
	if i.isProvisioned {
		return &Result{
			IsProvisioned: true,
		}, nil
	}

	quotas, err := i.provider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-A4707A72"),
		ServiceCode: aws.String("vpc"),
	})
	if err != nil {
		return nil, fmt.Errorf("getting igw quota: %w", err)
	}

	igws, err := i.provider.EC2().DescribeInternetGateways(nil)
	if err != nil {
		return nil, fmt.Errorf("getting current igw count: %w", err)
	}

	quota := int(*quotas.Quota.Value)
	count := len(igws.InternetGateways)
	available := quota - count

	return &Result{
		Required:    i.required,
		Available:   available,
		HasCapacity: i.required <= available,
		Description: "AWS VPC Internet Gateways",
	}, nil
}
