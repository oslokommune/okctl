// Package userpool implements the AWS Cognito UserPool
// - https://docs.amazonaws.cn/en_us/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpool.html
package userpool

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
)

const (
	temporaryPasswordValidityDays = 7
	minPasswordLength             = 8
)

// UserPool stores the state for a cloud formation
// cognito user pool
type UserPool struct {
	StoredName  string
	PoolName    string
	ClusterName string
}

// NamedOutputs returns the named outputs
func (p *UserPool) NamedOutputs() map[string]cloudformation.Output {
	return map[string]cloudformation.Output{
		p.Name(): {
			Value: p.Ref(),
			Export: &cloudformation.Export{
				Name: cloudformation.Sub(fmt.Sprintf("${AWS::StackName}-%s", p.Name())),
			},
		},
		"Arn": {
			Value: cloudformation.GetAtt(p.Name(), "Arn"),
			Export: &cloudformation.Export{
				Name: cloudformation.Sub(fmt.Sprintf("${AWS::StackName}-%s", "Arn")),
			},
		},
		"ProviderName": {
			Value: cloudformation.GetAtt(p.Name(), "ProviderName"),
			Export: &cloudformation.Export{
				Name: cloudformation.Sub(fmt.Sprintf("${AWS::StackName}-%s", "ProvideName")),
			},
		},
		"ProviderURL": {
			Value: cloudformation.GetAtt(p.Name(), "ProviderURL"),
			Export: &cloudformation.Export{
				Name: cloudformation.Sub(fmt.Sprintf("${AWS::StackName}-%s", "ProviderURL")),
			},
		},
	}
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
					Name:     "admin_only",
					Priority: 1,
				},
			},
		},
		AdminCreateUserConfig: &cognito.UserPool_AdminCreateUserConfig{
			AllowAdminCreateUserOnly: true,
			InviteMessageTemplate: &cognito.UserPool_InviteMessageTemplate{
				EmailMessage: "Your username is {username} and temporary password is {####}",
				EmailSubject: fmt.Sprintf("Your temporary password for %s", p.ClusterName),
			},
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
		EmailVerificationSubject: fmt.Sprintf("Your verification code for %s", p.ClusterName),
		EmailVerificationMessage: "Your verification code is {####}.",
		Policies: &cognito.UserPool_Policies{
			PasswordPolicy: &cognito.UserPool_PasswordPolicy{
				MinimumLength:                 minPasswordLength,
				RequireLowercase:              true,
				RequireNumbers:                true,
				RequireSymbols:                false,
				RequireUppercase:              true,
				TemporaryPasswordValidityDays: temporaryPasswordValidityDays,
			},
		},
		Schema: []cognito.UserPool_SchemaAttribute{
			{
				Mutable:  false,
				Name:     "email",
				Required: true,
			},
			{
				AttributeDataType: "String",
				Mutable:           true,
				Name:              "name",
				Required:          true,
			},
		},
		UserPoolAddOns: &cognito.UserPool_UserPoolAddOns{
			AdvancedSecurityMode: "AUDIT",
		},
		UserPoolName: p.PoolName,
		UsernameAttributes: []string{
			"email",
		},
		VerificationMessageTemplate: &cognito.UserPool_VerificationMessageTemplate{
			DefaultEmailOption: "CONFIRM_WITH_LINK",
			EmailMessage:       "Your verification code is {####}.",
			EmailSubject:       fmt.Sprintf("Your verification code for %s", p.ClusterName),
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
func New(clusterName string) *UserPool {
	return &UserPool{
		StoredName:  "UserPool",
		PoolName:    fmt.Sprintf("okctl-%s-userpool", clusterName),
		ClusterName: clusterName,
	}
}
