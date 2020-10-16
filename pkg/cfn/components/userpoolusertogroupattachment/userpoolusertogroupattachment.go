// Package userpooldomain provides functionality for setting
// up a domain with a user pool
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html
package userpoolusertogroupattachment

import (
	"fmt"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// UserPoolUser
type UserPoolUserToGroupAttachment struct {
	StoredName string
	GroupName  string
	Username   string
	UserPoolId string
	UserNamer  cfn.Namer
}

// TODO fix hardcoded group name. HMMM - wha tis nil?
// TODO make linter happy
func (attachment *UserPoolUserToGroupAttachment) Resource() cloudformation.Resource {
	fmt.Println("ATTACHMENT ***")
	fmt.Println(attachment)
	fmt.Println(attachment.UserNamer.Name())
	return &cognito.UserPoolUserToGroupAttachment{
		GroupName:                            attachment.GroupName,
		UserPoolId:                           attachment.UserPoolId,
		Username:                             attachment.Username,
		AWSCloudFormationDependsOn: []string{
			attachment.UserNamer.Name(),
		},
	}
}

// TODO make linter happy
func (d *UserPoolUserToGroupAttachment) Name() string {
	return d.StoredName
}

// TODO make linter happy
func New(user cfn.Namer, username, groupname, userPoolid string) *UserPoolUserToGroupAttachment {
	return &UserPoolUserToGroupAttachment{
		StoredName: "UserPoolUserToGroupAttachment",
		GroupName:  groupname,
		Username:   username,
		UserPoolId: userPoolid,
		UserNamer: user,
	}
}
