package servicequota

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// EipCheck is used to check if you have enough Elastic Ips
type EipCheck struct {
	Out           io.Writer
	CloudProvider v1alpha1.CloudProvider
	Required      int
}

// NewEipCheck makes a new instance of check for Elastic Ips
func NewEipCheck(out io.Writer, provider v1alpha1.CloudProvider, required int) *EipCheck {
	return &EipCheck{Out: out, CloudProvider: provider, Required: required}
}

// CheckAvailability determines if you will be able to make required Elastic Ip(s)
func (e *EipCheck) CheckAvailability() error {
	quota, err := e.CloudProvider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-0263D0A3"),
		ServiceCode: aws.String("ec2"),
	})
	if err != nil {
		return fmt.Errorf("failed to get eip quota: %w", err)
	}

	quotaValue := int(*quota.Quota.Value)

	eips, err := e.CloudProvider.EC2().DescribeAddresses(nil)
	if err != nil {
		return fmt.Errorf("failed to get eip count: %w", err)
	}

	currentCount := len(eips.Addresses)

	avail := quotaValue - currentCount

	fmt.Fprintf(e.Out,
		"You have %d EIPs, and the quota for %s is %d,\n%d are available and you need %d",
		currentCount, e.CloudProvider.Region(), quotaValue, avail, e.Required)

	if avail < e.Required {
		fmt.Fprintln(e.Out, aurora.Red(" ❌"))
		fmt.Fprintln(e.Out, "You need Elastic Ips that do not change for Egress traffic in nat gateways")

		return fmt.Errorf("not enough EIPs available")
	}

	fmt.Fprintln(e.Out, aurora.Green(" ✔"))

	return nil
}
