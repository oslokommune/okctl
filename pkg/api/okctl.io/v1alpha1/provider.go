package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// CloudProvider defines the interface for interacting with
// AWS cloud services
type CloudProvider interface {
	SSM() ssmiface.SSMAPI
	EC2() ec2iface.EC2API
	EKS() eksiface.EKSAPI
	ServiceQuotas() servicequotasiface.ServiceQuotasAPI
	Route53() route53iface.Route53API
	CloudFormation() cloudformationiface.CloudFormationAPI
	Region() string
	PrincipalARN() string
}
