// Package auth deals with setting up correct credential types from the user towards AWS & Github
// The values are passed on to the `okctl.Okctl` struct where they are used as plain strings
package auth

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/oslokommune/okctl/pkg/okctl"
)

// AwsCredentialsType is the chosen AWS credential type
// Should only be accessed by main.addAuthenticationFlags()
var AwsCredentialsType awsCredential //nolint: gochecknoglobals

type awsCredential string

func (c awsCredential) Validate() error {
	return validation.Validate(
		c.String(),
		validation.By(
			validateInList(GetAwsCredentialsTypes(), "validating AWS credentials"),
		),
	)
}

func (c awsCredential) String() string {
	return string(c)
}

// GithubCredentialsType is the chosen Github credential type
// Should only be accessed by main.addAuthenticationFlags()
var GithubCredentialsType githubCredential //nolint: gochecknoglobals

type githubCredential string

func (c githubCredential) Validate() error {
	return validation.Validate(
		c.String(),
		validation.By(
			validateInList(GetGithubCredentialsTypes(), "validating Github credentials"),
		),
	)
}

func (c githubCredential) String() string {
	return string(c)
}

// EnableServiceUserAuthentication sets the aws & github credential type on `o`
// The single point of write to the credential types on `okctl.Okctl`
func EnableServiceUserAuthentication(o *okctl.Okctl) {
	o.AWSCredentialsType = AwsCredentialsType.String()
	o.GithubCredentialsType = GithubCredentialsType.String()
}

// ValidateCredentialTypes validate that user has chosen valid credential types
func ValidateCredentialTypes() error {
	awsCredentialsTypes := GetAwsCredentialsTypes()
	githubCredentialsTypes := GetGithubCredentialsTypes()

	var err error

	err = AwsCredentialsType.Validate()
	if err != nil {
		return fmt.Errorf(
			"%s: type '%s' is not valid. Allowed values: %s. See %s for more information",
			err,
			AwsCredentialsType,
			strings.Join(awsCredentialsTypes, ","),
			constant.DefaultAwsAuthDocumentationURL,
		)
	}

	err = GithubCredentialsType.Validate()
	if err != nil {
		return fmt.Errorf(
			"%s: type '%s' is not valid. Allowed values: %s. See %s for more information",
			err,
			GithubCredentialsType,
			strings.Join(githubCredentialsTypes, ","),
			constant.DefaultAwsAuthDocumentationURL,
		)
	}

	return nil
}

// GetAwsCredentialsTypes gets all the valid AWS credential types we support
func GetAwsCredentialsTypes() []string {
	awsCredentialsTypes := []string{
		context.AWSCredentialsTypeAccessKey,
		context.AWSCredentialsTypeAwsProfile,
	}

	return awsCredentialsTypes
}

// GetGithubCredentialsTypes gets all the valid Github credential types we support
func GetGithubCredentialsTypes() []string {
	githubCredentialsTypes := []string{
		context.GithubCredentialsTypeDeviceAuthentication,
		context.GithubCredentialsTypeToken,
	}

	return githubCredentialsTypes
}

func validateInList(matchList []string, errorString string) validation.RuleFunc {
	return func(value interface{}) error {
		for _, s := range matchList {
			if value == s {
				return nil
			}
		}

		return errors.New(errorString)
	}
}
