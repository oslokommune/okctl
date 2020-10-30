// Package v1alpha1 defines the currently active version of the okctl api
package v1alpha1

import (
	"fmt"
)

const (
	// RegionEuWest1 defines the AWS region
	RegionEuWest1 = "eu-west-1"
	// OkPrincipalARNPattern defines what the Oslo kommune principal ARN for KeyCloak looks like
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/intro-structure.html#intro-structure-principal
	OkPrincipalARNPattern = "arn:aws:iam::%s:saml-provider/keycloak"
	// OkRoleARNPattern defines the standard role that we should select, i.e., the one with most privileges
	// This might be made configurable in the future, but for now it only makes sense to do it this way
	// as the other roles are pretty much useless in their current state
	OkRoleARNPattern = "arn:aws:iam::%s:role/oslokommune/iamadmin-SAML"
	// OkSamlURL is the starting point for authenticating via KeyCloak towards AWS
	OkSamlURL = "https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws"
)

// SupportedRegions returns the supported regions on AWS
func SupportedRegions() []string {
	return []string{
		RegionEuWest1,
	}
}

// IsSupportedRegion determines if a given region is supported
func IsSupportedRegion(region string) bool {
	for _, r := range SupportedRegions() {
		if region == r {
			return true
		}
	}

	return false
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
	default:
		return nil, fmt.Errorf("region: %s is not supported", region)
	}
}

// PrincipalARN returns the Ok principal ARN with account number set
func PrincipalARN(awsAccountID string) string {
	return fmt.Sprintf(OkPrincipalARNPattern, awsAccountID)
}

// RoleARN return the Ok role ARN with account number set
func RoleARN(awsAccountID string) string {
	return fmt.Sprintf(OkRoleARNPattern, awsAccountID)
}

// PermissionsBoundaryARN return the Ok permissions boundary ARN
func PermissionsBoundaryARN(awsAccountID string) string {
	return fmt.Sprintf("arn:aws:iam::%s:policy/oslokommune/oslokommune-boundary", awsAccountID)
}
