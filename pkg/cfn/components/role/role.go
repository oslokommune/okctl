// Package role knows how to create cloud formation
// for an IAM role
package role

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/iam"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// Role stores the state for a cloud formation iam role
type Role struct {
	RoleName                 string
	StoredName               string
	PermissionsBoundary      string
	ManagedPolicyARNs        []string
	AssumeRolePolicyDocument interface{}
	Policies                 map[string]interface{}
}

// NamedOutputs returns the resource outputs
func (r *Role) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(fmt.Sprintf("%sArn", r.Name()), cloudformation.GetAtt(r.Name(), "Arn")).NamedOutputs()
}

// Name returns the name of the resource
func (r *Role) Name() string {
	return r.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (r *Role) Ref() string {
	return cloudformation.Ref(r.Name())
}

// Resource returns the cloud formation resource for an IAM role
func (r *Role) Resource() cloudformation.Resource {
	role := &iam.Role{
		AssumeRolePolicyDocument: r.AssumeRolePolicyDocument,
		ManagedPolicyArns:        r.ManagedPolicyARNs,
		PermissionsBoundary:      r.PermissionsBoundary,
		RoleName:                 r.RoleName,
	}

	for policy, document := range r.Policies {
		role.Policies = append(role.Policies, iam.Role_Policy{
			PolicyDocument: document,
			PolicyName:     policy,
		})
	}

	return role
}

// New returns an initialised IAM role
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-role.html
func New(
	roleName, resourceName, permissionsBoundary string,
	managedPolicyARNs []string,
	assumeRolePolicyDocument interface{},
	policies map[string]interface{},
) *Role {
	return &Role{
		RoleName:                 roleName,
		StoredName:               resourceName,
		PermissionsBoundary:      permissionsBoundary,
		ManagedPolicyARNs:        managedPolicyARNs,
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		Policies:                 policies,
	}
}
