package servicequota

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// VpcCheck is used to check if you have enough vpcs
type VpcCheck struct {
	Out           io.Writer
	CloudProvider v1alpha1.CloudProvider
	Required      int
}

// NewVpcCheck make a new instance of check for VPCs
func NewVpcCheck(out io.Writer, provider v1alpha1.CloudProvider, required int) *VpcCheck {
	return &VpcCheck{Out: out, CloudProvider: provider, Required: required}
}

// CheckAvailability determines if you will be able to make required Virtual Private Cloud(s)
func (v *VpcCheck) CheckAvailability() error {
	quota, err := v.CloudProvider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-F678F1CE"),
		ServiceCode: aws.String("vpc"),
	})
	if err != nil {
		return fmt.Errorf("failed to get vpc quota: %w", err)
	}

	quotaValue := int(*quota.Quota.Value)

	vpcs, err := v.CloudProvider.EC2().DescribeVpcs(nil)
	if err != nil {
		return fmt.Errorf("failed to get vpc count: %w", err)
	}

	currentCount := len(vpcs.Vpcs)

	avail := quotaValue - currentCount

	fmt.Fprintf(v.Out,
		"You have %d VPCs, and the quota for %s is %d,\n%d are available and you need %d",
		currentCount, v.CloudProvider.Region(), quotaValue, avail, v.Required)

	if avail < v.Required {
		fmt.Println(aurora.Red(" ❌"))
		fmt.Println("You need a VPC to put your cluster in")

		return fmt.Errorf("not enough VPCs available")
	}

	fmt.Fprintln(v.Out, aurora.Green(" ✔"))

	return nil
}
