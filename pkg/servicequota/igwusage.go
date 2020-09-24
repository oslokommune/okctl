package servicequota

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

type igwusage struct {
	CloudProvider v1alpha1.CloudProvider
}

func (i igwusage) Count() (int, error) {
	igws, err := i.CloudProvider.EC2().DescribeInternetGateways(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get igw count: %w", err)
	}

	return getLengthOf(getStringMapOf(igws.String()), "InternetGateways")
}

func (i igwusage) Quota() (int, error) {
	quotas, err := i.CloudProvider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-A4707A72"),
		ServiceCode: aws.String("vpc"),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get igw quota: %w", err)
	}

	return int(*quotas.Quota.Value), nil
}

func (igwusage) Required() int {
	return config.DefaultRequiredIgws
}

func (igwusage) Description() string {
	return "You need a internet gateway to enable traffic from outside AWS."
}

func (i igwusage) CheckAvailability() error {
	igwCount, err := i.Count()
	if err != nil {
		return err
	}

	igwQuota, err := i.Quota()
	if err != nil {
		return err
	}

	igwAvail := igwQuota - igwCount
	fmt.Printf("You have %d IGWs, and the quota for %s is %d,\n%d are available and you need %d", igwCount, i.CloudProvider.Region(), igwQuota, igwAvail, i.Required())

	if igwAvail < i.Required() {
		fmt.Println(aurora.Red(" ❌"))
		fmt.Println(i.Description())

		return fmt.Errorf("not enough IGWs available")
	}

	fmt.Println(aurora.Green(" ✔"))

	return nil
}
