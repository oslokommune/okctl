package servicequota

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// IgwCheck is used to check if you have enough Internet Gatweays
type IgwCheck struct {
	Out           io.Writer
	CloudProvider v1alpha1.CloudProvider
	Required      int
}

// NewIgwCheck makes a new instance of check for Internet Gateways
func NewIgwCheck(out io.Writer, provider v1alpha1.CloudProvider, required int) *IgwCheck {
	return &IgwCheck{Out: out, CloudProvider: provider, Required: required}
}

// CheckAvailability determines if you will be able to make required Internet Gatweay(s)
func (i *IgwCheck) CheckAvailability() error {
	quotas, err := i.CloudProvider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-A4707A72"),
		ServiceCode: aws.String("vpc"),
	})
	if err != nil {
		return fmt.Errorf("failed to get igw quota: %w", err)
	}

	quotaValue := int(*quotas.Quota.Value)

	igws, err := i.CloudProvider.EC2().DescribeInternetGateways(nil)
	if err != nil {
		return fmt.Errorf("failed to get igw count: %w", err)
	}

	currentCount := len(igws.InternetGateways)

	avail := quotaValue - currentCount

	fmt.Fprintf(i.Out,
		"You have %d IGWs, and the quota for %s is %d,\n%d are available and you need %d",
		currentCount, i.CloudProvider.Region(), quotaValue, avail, i.Required)

	if avail < i.Required {
		fmt.Fprintln(i.Out, aurora.Red(" ❌"))
		fmt.Fprintln(i.Out, "You need a internet gateway to enable traffic from outside AWS.")

		return fmt.Errorf("not enough IGWs available")
	}

	fmt.Fprintln(i.Out, aurora.Green(" ✔"))

	return nil
}
