// Package rotationschedule knows how to build a secrets manager secret
// rotation schedule
package rotationschedule

import (
	"fmt"
	"strings"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/secretsmanager"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// RotationSchedule contains the required state for building
// the cloud formation resource
type RotationSchedule struct {
	StoredName     string
	Secret         cfn.Referencer
	Attachment     cfn.Namer
	RotationType   string
	SecurityGroup  cfn.Namer
	VPCDBSubnetIDs []string
}

// NamedOutputs returns the resource outputs
func (e *RotationSchedule) NamedOutputs() map[string]cloudformation.Output {
	return nil
}

// Name returns the name of the cloud formation resource
func (e *RotationSchedule) Name() string {
	return e.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (e *RotationSchedule) Ref() string {
	return cloudformation.Ref(e.Name())
}

const (
	rotateAfterDays = 30
)

// Resource returns the cloud formation resource for a secret target attachment
func (e *RotationSchedule) Resource() cloudformation.Resource {
	return &secretsmanager.RotationSchedule{
		HostedRotationLambda: &secretsmanager.RotationSchedule_HostedRotationLambda{
			RotationLambdaName:  fmt.Sprintf("%s-SecretsManagerRotation", e.Name()),
			RotationType:        e.RotationType,
			VpcSecurityGroupIds: cloudformation.GetAtt(e.SecurityGroup.Name(), "GroupId"),
			VpcSubnetIds:        strings.Join(e.VPCDBSubnetIDs, ","),
		},
		RotationRules: &secretsmanager.RotationSchedule_RotationRules{
			AutomaticallyAfterDays: rotateAfterDays,
		},
		SecretId: e.Secret.Ref(),
		AWSCloudFormationDependsOn: []string{
			e.Attachment.Name(),
		},
	}
}

// New returns an initialised secrets manager secret rotation schedule
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-rotationschedule.html
func New(
	resourceName, rotationType string,
	secret cfn.Referencer,
	attachment cfn.Namer,
	vpcDBSubnetIDs []string,
	securityGroup cfn.Namer,
) *RotationSchedule {
	return &RotationSchedule{
		StoredName:     resourceName,
		Secret:         secret,
		Attachment:     attachment,
		RotationType:   rotationType,
		SecurityGroup:  securityGroup,
		VPCDBSubnetIDs: vpcDBSubnetIDs,
	}
}

// NewPostgres returns an initialised secrets manager secret rotation schedule
// compatible with postgres
func NewPostgres(
	resourceName string,
	secret cfn.Referencer,
	attachment cfn.Namer,
	vpcDBSubnetIDs []string,
	securityGroup cfn.Namer,
) *RotationSchedule {
	return New(
		resourceName,
		"PostgreSQLSingleUser",
		secret,
		attachment,
		vpcDBSubnetIDs,
		securityGroup,
	)
}
