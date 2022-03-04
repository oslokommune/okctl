// Package auth deals with credential types towards AWS & Github
package auth

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/oslokommune/okctl/pkg/okctl"
)

var (
	// AwsCredentialsType is the chosen AWS credential type
	AwsCredentialsType string //nolint:gochecknoglobals
	// GithubCredentialsType is the chosen Github credential type
	GithubCredentialsType string //nolint:gochecknoglobals
)

// EnableServiceUserAuthentication sets the aws & github credential type on `o`
func EnableServiceUserAuthentication(o *okctl.Okctl) {
	o.AWSCredentialsType = AwsCredentialsType
	o.GithubCredentialsType = GithubCredentialsType
}

// ValidateCredentialTypes validate that user has chosen valid credential types
func ValidateCredentialTypes() error {
	awsCredentialsTypes := GetAwsCredentialsTypes()
	githubCredentialsTypes := GetGithubCredentialsTypes()

	if !contains(awsCredentialsTypes, AwsCredentialsType) {
		return fmt.Errorf(
			"invalid AWS credentials type '%s'. Allowed values: %s. See %s for more information",
			AwsCredentialsType,
			strings.Join(awsCredentialsTypes, ","),
			constant.DefaultAwsAuthDocumentationURL,
		)
	}

	if !contains(githubCredentialsTypes, GithubCredentialsType) {
		return fmt.Errorf(
			"invalid Github credentials type '%s'. Allowed values: %s",
			GithubCredentialsType,
			strings.Join(githubCredentialsTypes, ","),
		)
	}

	return nil
}

// GetAwsCredentialsTypes gets all the valid AWS credential types we support
func GetAwsCredentialsTypes() []string {
	awsCredentialsTypes := []string{
		context.AWSCredentialsTypeSAML,
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

func contains(l []string, v string) bool {
	for _, el := range l {
		if v == el {
			return true
		}
	}

	return false
}
