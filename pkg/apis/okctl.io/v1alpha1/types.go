package v1alpha1

import "fmt"

const (
	RegionEuWest1 = "eu-west-1"

	OkPrincipalARNPattern = "arn:aws:iam::%s:saml-provider/keycloak"
	OkRoleARNPattern      = "arn:aws:iam::%s:role/oslokommune/iamadmin-SAML"

	OkSamlURL = "https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws"
)

func SupportedRegions() []string {
	return []string{
		RegionEuWest1,
	}
}

func IsSupportedRegion(region string) bool {
	for _, r := range SupportedRegions() {
		if region == r {
			return true
		}
	}

	return false
}

func SupportedAvailabilityZones(region string) ([]string, error) {
	if !IsSupportedRegion(region) {
		return nil, fmt.Errorf("region: %s is not supported", region)
	}

	var azs []string

	// nolint
	switch region {
	case RegionEuWest1:
		azs = []string{
			"eu-west-1a",
			"eu-west-1b",
			"eu-west-1c",
		}
	}

	return azs, nil
}

func PrincipalARN(awsAccountID string) string {
	return fmt.Sprintf(OkPrincipalARNPattern, awsAccountID)
}

func RoleARN(awsAccountID string) string {
	return fmt.Sprintf(OkRoleARNPattern, awsAccountID)
}
