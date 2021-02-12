// Package servicequota check if you have enough resources in aws before cluster creation starts
package servicequota

import (
	"fmt"
)

// Checker defines what we need to know about a service quota
type Checker interface {
	CheckAvailability() (*Result, error)
}

// Result contains the data from an availability check
type Result struct {
	Required    int
	Available   int
	HasCapacity bool
	Description string
}

// CheckQuotas will run through the checks and return an error if quotas are too small
func CheckQuotas(checks ...Checker) error {
	for _, check := range checks {
		r, err := check.CheckAvailability()
		if err != nil {
			return err
		}

		if !r.HasCapacity {
			return fmt.Errorf("%s: required %d, but only have %d available", r.Description, r.Required, r.Available)
		}
	}

	return nil
}
