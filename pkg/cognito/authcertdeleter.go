package cognito

import (
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// AuthCertDeleter is used just to delete the auth. cert in us-east-1
type AuthCertDeleter struct {
	usprovider v1alpha1.CloudProvider
}

// DeleteAuthCertOpts options needed to delete an auth certificate
type DeleteAuthCertOpts struct {
	Domain string
}

// Validate the inputs to delete an auth certificate
func (o *DeleteAuthCertOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Domain,
			validation.Required,
			validation.Match(regexp.MustCompile("^auth")).Error(
				"Auth cert domain must start with 'auth'. It was "+o.Domain)),
	)
}

// DeleteAuthCert delete the auth certificate that was made in us-east-1 for identitypool
func (c *AuthCertDeleter) DeleteAuthCert(opts DeleteAuthCertOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	stacklist, err := c.usprovider.CloudFormation().ListStacks(&cloudformation.ListStacksInput{
		NextToken:         nil,
		StackStatusFilter: nil,
	})
	if err != nil {
		return err
	}

	for _, stack := range stacklist.StackSummaries {
		if strings.Contains(*stack.StackId, strings.ReplaceAll(opts.Domain, ".", "-")) {
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
