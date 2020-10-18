// Package userpoolgroup provides functionality for setting
// up a user pool group
package userpoolgroup

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn"
)

const (
	groupname  = "admins"
	precedence = 10
)

// UserPoolGroup output
type UserPoolGroup struct {
	StoredName  string
	Description string
	GroupName   string
	Precedence  int
	UserPool    cfn.NameReferencer
}

// Resource returns the cloud formation resource for a cognito user pool domain
func (d *UserPoolGroup) Resource() cloudformation.Resource {
	return &cognito.UserPoolGroup{
		Description: "",
		GroupName:   groupname,
		Precedence:  precedence,
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
