package servicequota

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

type eipusage struct {
	CloudProvider v1alpha1.CloudProvider
}

func (e eipusage) Count() (int, error) {
	eips, err := e.CloudProvider.EC2().DescribeAddresses(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get eip count: %w", err)
	}
	return len(eips.Addresses), nil
}

func (e eipusage) Quota() (int, error) {
	quotas, err := e.CloudProvider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-0263D0A3"),
		ServiceCode: aws.String("ec2"),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get eip quota: %w", err)
	}

	return int(*quotas.Quota.Value), nil
}

func (eipusage) Required() int {
	return config.DefaultRequiredEpis
}

func (eipusage) Description() string {
	return "You need Elastic Ips that do not change for Egress traffic in nat gateways"
}

func (e eipusage) CheckAvailability() error {
	eipCount, err := e.Count()
	if err != nil {
		return err
	}

	eipQuota, err := e.Quota()
	if err != nil {
		return err
	}

	eipAvail := eipQuota - eipCount
	fmt.Printf("You have %d EIPs, and the quota for %s is %d,\n%d are available and you need %d", eipCount, e.CloudProvider.Region(), eipQuota, eipAvail, e.Required())

	if eipAvail < e.Required() {
		fmt.Println(aurora.Red(" ❌"))
		fmt.Println(e.Description())

		return fmt.Errorf("not enough EIPs available")
	}

	fmt.Println(aurora.Green(" ✔"))

	return nil
}
