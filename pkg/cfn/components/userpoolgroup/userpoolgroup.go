// Package userpooldomain provides functionality for setting
// up a domain with a user pool
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html
package userpoolgroup

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// UserPoolGroup
type UserPoolGroup struct {
	StoredName  string
	Description string
	GroupName   string
	Precedence  int
	UserPool    cfn.NameReferencer
}

// Resource returns the cloud formation resource for a
// cognito user pool domain
// TODO fix hardcoded group name.
func (d *UserPoolGroup) Resource() cloudformation.Resource {
	return &cognito.UserPoolGroup{
		Description: "",
		GroupName:   "admins",
		Precedence:  10,
		UserPoolId:  d.UserPool.Ref(),
		AWSCloudFormationDependsOn: []string{
			d.UserPool.Name(),
		},
	}
}

// Name returns the logical id of the resource
func (d *UserPoolGroup) Name() string {
	return d.StoredName
}

// New returns an initialised cognito user pool domain
func New(groupName string, description string, userPool cfn.NameReferencer) *UserPoolGroup {
	return &UserPoolGroup{
		StoredName:  "UserPoolGroup",
		Description: description,
		GroupName:   groupName,
		Precedence:  0,
		UserPool:    userPool,
	}
}
