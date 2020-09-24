// Package servicequota check if you have enough resources in aws before cluster creation starts
package servicequota

import (
	"github.com/hjson/hjson-go"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

// Usage defines what we need to know about a service quota
type Usage interface {
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

func getStringMapOf(humanJSON string) map[string]interface{} {
	var result map[string]interface{}

	err := hjson.Unmarshal([]byte(humanJSON), &result)
	if err != nil {
		return nil
	}

	return result
}

func getLengthOf(result map[string]interface{}, path string) (int, error) {
	if result[path] == nil {
		return 0, nil
	}

	return len(result[path].([]interface{})), nil
}
