package servicequota

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

type vpcusage struct {
	// TODO don't need to export cloudprovider?
	CloudProvider v1alpha1.CloudProvider
}

func (v vpcusage) Count() (int, error) {
	vpcs, err := v.CloudProvider.EC2().DescribeVpcs(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get vpc count: %w", err)
	}
	return len(vpcs.Vpcs), nil
}

func (v vpcusage) Quota() (int, error) {
	quotas, err := v.CloudProvider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-F678F1CE"),
		ServiceCode: aws.String("vpc"),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get vpc quota: %w", err)
	}

	return int(*quotas.Quota.Value), nil
}

func (vpcusage) Required() int {
	return config.DefaultRequiredVpcs
}

func (vpcusage) Description() string {
	return "You need a VPC to put your cluster in"
}

func (v vpcusage) CheckAvailability() error {
	vpcQuota, err := v.Quota()
	if err != nil {
		return err
	}

	vpcCount, err := v.Count()
	if err != nil {
		return err
	}

	vpcAvail := vpcQuota - vpcCount
	fmt.Printf("You have %d VPCs, and the quota for %s is %d,\n%d are available and you need %d", vpcCount, v.CloudProvider.Region(), vpcQuota, vpcAvail, v.Required())

	if vpcAvail < v.Required() {
		fmt.Println(aurora.Red(" ❌"))
		fmt.Println(v.Description())

		return fmt.Errorf("not enough VPCs available")
	}

	fmt.Println(aurora.Green(" ✔"))

	return nil
}
