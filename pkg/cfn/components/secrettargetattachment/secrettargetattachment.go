// Package secrettargetattachment knows how to build a secret target attachment
package secrettargetattachment

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/secretsmanager"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// SecretTargetAttachment contains the state required for building
// a secret target attachment resource
type SecretTargetAttachment struct {
	StoredName string
	TargetType string
	Secret     cfn.Referencer
	Target     cfn.Referencer
}

// NamedOutputs returns the resource outputs
func (e *SecretTargetAttachment) NamedOutputs() map[string]cloudformation.Output {
	return nil
}

// Name returns the name of the cloud formation resource
func (e *SecretTargetAttachment) Name() string {
	return e.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (e *SecretTargetAttachment) Ref() string {
	return cloudformation.Ref(e.Name())
}

// Resource returns the cloud formation resource for a secret target attachment
func (e *SecretTargetAttachment) Resource() cloudformation.Resource {
	return &secretsmanager.SecretTargetAttachment{
		SecretId:   e.Secret.Ref(),
		TargetId:   e.Target.Ref(),
		TargetType: e.TargetType,
	}
}

// New returns an initialised secrets target attachment
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-secrettargetattachment.html
func New(resourceName, targetType string, secret, target cfn.Referencer) *SecretTargetAttachment {
	return &SecretTargetAttachment{
		StoredName: resourceName,
		TargetType: targetType,
		Secret:     secret,
		Target:     target,
	}
}

// NewRDSDBInstance returns an initialised secrets target attachment
// compatible with RDS DBInstance
func NewRDSDBInstance(resourceName string, secret, target cfn.Referencer) *SecretTargetAttachment {
	return New(resourceName, "AWS::RDS::DBInstance", secret, target)
}
