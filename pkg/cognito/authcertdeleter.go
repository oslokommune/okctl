package cognito

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// AuthCertDeleter is used just to delete the auth. cert in us-east-1
type AuthCertDeleter struct {
	usprovider v1alpha1.CloudProvider
}

// DeleteAuthCert delete the auth certificate that was made in us-east-1 for identitypool
func (c *AuthCertDeleter) DeleteAuthCert(domain string) error {
	stacklist, err := c.usprovider.CloudFormation().ListStacks(&cloudformation.ListStacksInput{
		NextToken:         nil,
		StackStatusFilter: nil,
	})
	if err != nil {
		return err
	}

	for _, stack := range stacklist.StackSummaries {
		if strings.Contains(*stack.StackId, strings.ReplaceAll(domain, ".", "-")) {
			_, err := c.usprovider.CloudFormation().DeleteStack(&cloudformation.DeleteStackInput{
				StackName: stack.StackId,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NewCertDeleter returns an initialised AuthCertDeleter
func NewCertDeleter(usprovider v1alpha1.CloudProvider) *AuthCertDeleter {
	return &AuthCertDeleter{
		usprovider: usprovider,
	}
}
