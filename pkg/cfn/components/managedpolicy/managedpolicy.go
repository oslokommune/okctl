// Package managedpolicy knows how to create cloud formation
// for a managed IAM policy
package managedpolicy

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/iam"
)

// ManagedPolicy stores the state for a cloud formation managed
// IAM policy
type ManagedPolicy struct {
	StoredName  string
	PolicyName  string
	Description string
	Document    interface{}
}

// Resource returns the cloud formation resource for an IAM policy
func (p *ManagedPolicy) Resource() cloudformation.Resource {
	return &iam.ManagedPolicy{
		Description:       p.Description,
		ManagedPolicyName: p.PolicyName,
		PolicyDocument:    p.Document,
	}
}

// Name returns the name of the resource
func (p *ManagedPolicy) Name() string {
	return p.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (p *ManagedPolicy) Ref() string {
	return cloudformation.Ref(p.StoredName)
}

// New creates a managed IAM policy
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-managedpolicy.html
func New(resourceName, policyName, description string, document interface{}) *ManagedPolicy {
	return &ManagedPolicy{
		StoredName:  resourceName,
		PolicyName:  policyName,
		Description: description,
		Document:    document,
	}
}
