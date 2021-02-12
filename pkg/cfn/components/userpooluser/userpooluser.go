// Package userpooluser provides functionality for setting
// up a user pool user
package userpooluser

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// UserPoolUser output
type UserPoolUser struct {
	StoredName  string
	Email       string
	Description string
	UserPoolID  string
}

// NamedOutputs stored names of users
func (u *UserPoolUser) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(fmt.Sprintf("%sID", u.StoredName), u.Ref()).NamedOutputs()
}

// Resource for type UserPoolUser
func (u *UserPoolUser) Resource() cloudformation.Resource {
	return &cognito.UserPoolUser{
		ClientMetadata:         nil,
		DesiredDeliveryMediums: []string{"EMAIL"},
		ForceAliasCreation:     true,
		UserAttributes: []cognito.UserPoolUser_AttributeType{
			{
				Name:  "email",
				Value: u.Email,
			},
			{
				Name:  "name",
				Value: u.Email,
			},
			{
				Name:  "email_verified",
				Value: "True",
			},
		},
		Username:                             u.Email,
		UserPoolId:                           u.UserPoolID,
		AWSCloudFormationDeletionPolicy:      "",
		AWSCloudFormationUpdateReplacePolicy: "",
		AWSCloudFormationDependsOn:           nil,
		AWSCloudFormationMetadata:            nil,
		AWSCloudFormationCondition:           "",
	}
}

// Name stored name for user pool user
func (u *UserPoolUser) Name() string {
	return u.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (u *UserPoolUser) Ref() string {
	return cloudformation.Ref(u.Name())
}

// New creates a new UserPoolUser
func New(email string, description string, userPoolid string) *UserPoolUser {
	return &UserPoolUser{
		StoredName:  "UserPoolUser",
		Email:       email,
		Description: description,
		UserPoolID:  userPoolid,
	}
}
