// Package iamapi provides some convenience functions
// for interacting with the AWS IAM API
package iamapi

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// IAMAPI contains the required state for interacting
// with the AWS IAM API
type IAMAPI struct {
	provider v1alpha1.CloudProvider
}

// New returns an initialised client
func New(provider v1alpha1.CloudProvider) *IAMAPI {
	return &IAMAPI{
		provider: provider,
	}
}

// AttachRolePolicy attaches the provided policy to the role
func (i *IAMAPI) AttachRolePolicy(policyARN, roleARN string) error {
	friendlyName, err := RoleFriendlyName(roleARN)
	if err != nil {
		return err
	}

	_, err = i.provider.IAM().AttachRolePolicy(&iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyARN),
		RoleName:  aws.String(friendlyName),
	})
	if err != nil {
		return fmt.Errorf("attaching policy to role: %w", err)
	}

	return nil
}

// DetachRolePolicy detaches the provided policy from the role
func (i *IAMAPI) DetachRolePolicy(policyARN, roleARN string) error {
	friendlyName, err := RoleFriendlyName(roleARN)
	if err != nil {
		return err
	}

	_, err = i.provider.IAM().DetachRolePolicy(&iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyARN),
		RoleName:  aws.String(friendlyName),
	})
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("NoSuchEntity: The role with name %s cannot be found.", friendlyName)) {
			return nil
		}

		if strings.Contains(err.Error(), fmt.Sprintf("NoSuchEntity: Policy %s was not found.", policyARN)) {
			return nil
		}

		return fmt.Errorf("detaching policy from role: %w", err)
	}

	return nil
}

// RoleFriendlyName returns the friendly name of the role
// by extracting it from the provided ARN
func RoleFriendlyName(roleARN string) (string, error) {
	a, err := arn.Parse(roleARN)
	if err != nil {
		return "", fmt.Errorf("getting role friendly name: %w", err)
	}

	return strings.TrimPrefix(a.Resource, "role/"), nil
}
