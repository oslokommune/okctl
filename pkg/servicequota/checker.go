// Package servicequota check if you have enough resources in aws before cluster creation starts
package servicequota

// Checker defines what we need to know about a service quota
type Checker interface {
	CheckAvailability() error
}

// CheckQuotas will run through the checks and return an error if quotas are too small
func CheckQuotas(checkers ...Checker) error {
	for i := range checkers {
		checker := checkers[i]

		err := checker.CheckAvailability()
		if err != nil {
			return err
		}
	}

	return nil
}
