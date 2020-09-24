// Package servicequota check if you have enough resources in aws before cluster creation starts
package servicequota

// Checker defines what we need to know about a service quota
type Checker interface {
	CheckAvailability() error
}
