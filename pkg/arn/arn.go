// Package arn knows how to validate arns
package arn

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/config/constant"

	arnpkg "github.com/aws/aws-sdk-go/aws/arn"
)

// ServiceType enumerates known AWS services
type ServiceType string

// String returns the string
func (s ServiceType) String() string {
	return string(s)
}

// ResourceType enumerates known AWS resources
type ResourceType string

// String returns the string
func (t ResourceType) String() string {
	return string(t)
}

// nolint: golint
const (
	ServiceElasticLoadBalancing ServiceType = "elasticloadbalancing"

	ResourceLoadBalancer ResourceType = "loadbalancer/app/"
)

// Parse and validate the provided ARN
func Parse(arn string) (*arnpkg.ARN, error) {
	if !arnpkg.IsARN(arn) {
		return nil, fmt.Errorf(constant.InvalidArnError, arn)
	}

	a, err := arnpkg.Parse(arn)
	if err != nil {
		return nil, fmt.Errorf(constant.ParseArnError, err)
	}

	return &a, nil
}

// Is returns true if the arn if of expected service and resource type
// nolint: interfacer
func Is(service ServiceType, resource ResourceType, arn string) (bool, error) {
	a, err := Parse(arn)
	if err != nil {
		return false, err
	}

	if a.Service == service.String() && strings.HasPrefix(a.Resource, resource.String()) {
		return true, nil
	}

	return false, nil
}

// IsLoadBalancer returns true if the provided ARN
// is an ARN and has correct service and resource type
func IsLoadBalancer(arn string) (bool, error) {
	return Is(ServiceElasticLoadBalancing, ResourceLoadBalancer, arn)
}
