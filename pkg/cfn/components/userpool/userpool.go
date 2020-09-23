// Package userpool implements the AWS Cognito UserPool
// - https://docs.amazonaws.cn/en_us/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpool.html
package userpool

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
)

const (
	unusedAccountValidityDays     = 7
	temporaryPasswordValidityDays = 7
	minPasswordLength             = 8
)

// UserPool stores the state for a cloud formation
// cognito user pool
type UserPool struct {
	StoredName  string
	PoolName    string
	Environment string
	Repository  string
}

// NamedOutputs returns the named outputs
func (p *UserPool) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{}
}

// Resource returns the cloud formation resource for a
// cognito user pool
// nolint: funlen
func (p *UserPool) Resource() cloudformation.Resource {
	return &cognito.UserPool{
		AccountRecoverySetting: &cognito.UserPool_AccountRecoverySetting{
			// Only allow an administrator to reset access
			RecoveryMechanisms: []cognito.UserPool_RecoveryOption{
				{
					Name: "admin_only",
				},
			},
		},
		AdminCreateUserConfig: &cognito.UserPool_AdminCreateUserConfig{
			AllowAdminCreateUserOnly: true,
			InviteMessageTemplate: &cognito.UserPool_InviteMessageTemplate{
				EmailMessage: "Your username is {username} and temporary password is {####}.",
				EmailSubject: fmt.Sprintf("Your temporary password for %s (%s)", p.Repository, p.Environment),
			},
			UnusedAccountValidityDays: unusedAccountValidityDays,
		},
		AliasAttributes: []string{
			"email",
		},
		AutoVerifiedAttributes: []string{
			"email",
		},
		DeviceConfiguration: &cognito.UserPool_DeviceConfiguration{
			ChallengeRequiredOnNewDevice:     true,
			DeviceOnlyRememberedOnUserPrompt: true,
		},
		// Consider verifying a domain, so we can send from our own
		// email address, using something like:
		// - https://github.com/binxio/cfn-ses-provider
		EmailConfiguration: &cognito.UserPool_EmailConfiguration{
			EmailSendingAccount: "COGNITO_DEFAULT",
		},
		EmailVerificationMessage: fmt.Sprintf("Your verification code for %s (%s)", p.Repository, p.Environment),
		EmailVerificationSubject: "Your verification code is {####}.",
		EnabledMfas: []string{
			"SOFTWARE_TOKEN_MFA",
		},
		MfaConfiguration: "ON",
		Policies: &cognito.UserPool_Policies{
			PasswordPolicy: &cognito.UserPool_PasswordPolicy{
				MinimumLength:                 minPasswordLength,
				RequireLowercase:              true,
				RequireNumbers:                true,
				RequireSymbols:                true,
				RequireUppercase:              true,
				TemporaryPasswordValidityDays: temporaryPasswordValidityDays,
			},
		},
		Schema: []cognito.UserPool_SchemaAttribute{
			{
				AttributeDataType:      "String",
				DeveloperOnlyAttribute: true,
				Mutable:                false,
				Name:                   "email",
				Required:               true,
			},
		},
		UserPoolAddOns: &cognito.UserPool_UserPoolAddOns{
			AdvancedSecurityMode: "AUDIT",
		},
		UserPoolName: p.PoolName,
		UsernameAttributes: []string{
			"email",
		},
		UsernameConfiguration: &cognito.UserPool_UsernameConfiguration{
			CaseSensitive: false,
		},
		VerificationMessageTemplate: &cognito.UserPool_VerificationMessageTemplate{
			DefaultEmailOption: "CONFIRM_WITH_LINK",
			EmailMessage:       "Your verification code is {####}.",
			EmailSubject:       fmt.Sprintf("Your verification code for %s (%s)", p.Repository, p.Environment),
		},
	}
}

// Name returns the logical identifier of the resource
func (p *UserPool) Name() string {
	return p.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (p *UserPool) Ref() string {
	return cloudformation.Ref(p.Name())
}

// New returns an initialised cognito user pool
func New(env, repository string) *UserPool {
	return &UserPool{
		StoredName:  "UserPool",
		PoolName:    fmt.Sprintf("okctl-%s-%s-userpool", env, repository),
		Environment: env,
		Repository:  repository,
	}
}
