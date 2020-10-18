// Package userpoolusertogroupattachment provides functionality for setting
// up a user pool user to group attachment
package userpoolusertogroupattachment

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// UserPoolUserToGroupAttachment attaches a user pool user to a group
type UserPoolUserToGroupAttachment struct {
	StoredName string
	GroupName  string
	Username   string
	UserPoolID string
	UserNamer  cfn.Namer
}

// Resource of type UserPoolUserToGroupAttachment
func (a *UserPoolUserToGroupAttachment) Resource() cloudformation.Resource {
	return &cognito.UserPoolUserToGroupAttachment{
		GroupName:  a.GroupName,
		UserPoolId: a.UserPoolID,
		Username:   a.Username,
		AWSCloudFormationDependsOn: []string{
			a.UserNamer.Name(),
		},
	}
}

// Name returns the logical identifier
func (a *UserPoolUserToGroupAttachment) Name() string {
	return a.StoredName
}

// New constructor
func New(user cfn.Namer, username, groupname, userPoolid string) *UserPoolUserToGroupAttachment {
	return &UserPoolUserToGroupAttachment{
		StoredName: "UserPoolUserToGroupAttachment",
		GroupName:  groupname,
		Username:   username,
		UserPoolID: userPoolid,
		UserNamer:  user,
	}
}
