package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// CloudProvider defines the interface for interacting with
// AWS cloud services
type CloudProvider interface {
	SSM() ssmiface.SSMAPI
	EC2() ec2iface.EC2API
	EKS() eksiface.EKSAPI
	CloudFormation() cloudformationiface.CloudFormationAPI
	Region() string
	PrincipalARN() string
}
