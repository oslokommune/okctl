// Package servicequota check if you have enough resources in aws before cluster creation starts
package servicequota

import (
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

// Checker defines what we need to know about a service quota
type Checker interface {
	Count() (int, error)
	Quota() (int, error)
	Required() int
	CheckAvailability() error
	Description() string
}

// CheckQuotas check if we have enough resources for given cloud provider
func CheckQuotas(provider v1alpha1.CloudProvider) error {
	vpcs := vpcusage{}
	vpcs.CloudProvider = provider

	err := vpcs.CheckAvailability()
	if err != nil {
		return err
	}

	eips := eipusage{}
	eips.CloudProvider = provider

	err = eips.CheckAvailability()
	if err != nil {
		return err
	}

	igws := igwusage{}
	igws.CloudProvider = provider

	err = igws.CheckAvailability()
	if err != nil {
		return err
	}

	return nil
}
