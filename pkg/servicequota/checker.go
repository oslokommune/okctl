// Package servicequota check if you have enough resources in aws before cluster creation starts
package servicequota

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
)

// Checker defines what we need to know about a service quota
type Checker interface {
	CheckAvailability() (*Result, error)
}

// Result contains the data from an availability check
type Result struct {
	Required      int
	Available     int
	IsProvisioned bool
	HasCapacity   bool
	Description   string
}

// CheckQuotas will run through the checks and return an error if quotas are too small
func CheckQuotas(checks ...Checker) error {
	for _, check := range checks {
		r, err := check.CheckAvailability()
		if err != nil {
			return err
		}

		if r.IsProvisioned {
			continue
		}

		if !r.HasCapacity {
			return fmt.Errorf(constant.HasCapacityError, r.Description, r.Required, r.Available)
		}
	}

	return nil
}
