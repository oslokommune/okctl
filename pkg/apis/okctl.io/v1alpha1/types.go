// Package v1alpha1 defines the currently active version of the okctl api
package v1alpha1

import (
	"fmt"
)

const (
	// RegionEuWest1 defines the AWS region
	RegionEuWest1 = "eu-west-1"
	// RegionEuCentral1 defines the AWS region
	RegionEuCentral1 = "eu-central-1"
	// RegionEuNorth1 defines the AWS region
	RegionEuNorth1 = "eu-north-1"

	// OkctlVersionTag defines the version of okctl used to provision the given resources
	OkctlVersionTag = "alpha.okctl.io/okctl-version"
	// OkctlCommitTag defines the git commit hash used to provision the given resources
	OkctlCommitTag = "alpha.okctl.io/okctl-commit"
	// OkctlManagedTag defines if this resource is managed by okctl
	OkctlManagedTag = "alpha.okctl.io/managed"
	// OkctlClusterNameTag defines the name of the cluster
	OkctlClusterNameTag = "alpha.okctl.io/cluster-name"
)

// SupportedRegions returns the supported regions on AWS
func SupportedRegions() []string {
	return []string{
		RegionEuWest1,
		RegionEuCentral1,
	}
}

// SupportedAvailabilityZones returns the availability zones for a
// region, given we know of it
func SupportedAvailabilityZones(region string) ([]string, error) {
	switch region {
	case RegionEuWest1:
		return []string{
			"eu-west-1a",
			"eu-west-1b",
			"eu-west-1c",
		}, nil
	case RegionEuCentral1:
		return []string{
			"eu-central-1a",
			"eu-central-1b",
			"eu-central-1c",
		}, nil
	case RegionEuNorth1:
		return []string{
			"eu-north-1a",
			"eu-north-1b",
			"eu-north-1c",
		}, nil
	default:
		return nil, fmt.Errorf("region: %s is not supported", region)
	}
}

// PermissionsBoundaryARN return the Ok permissions boundary ARN
func PermissionsBoundaryARN(awsAccountID string) string {
	return fmt.Sprintf("arn:aws:iam::%s:policy/oslokommune/oslokommune-boundary", awsAccountID)
}
