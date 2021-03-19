package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// CloudProvider defines the interface for interacting with
// AWS cloud services
type CloudProvider interface {
	SecretsManager() secretsmanageriface.SecretsManagerAPI
	S3() s3iface.S3API
	IAM() iamiface.IAMAPI
	SSM() ssmiface.SSMAPI
	EC2() ec2iface.EC2API
	EKS() eksiface.EKSAPI
	ELBV2() elbv2iface.ELBV2API
	ServiceQuotas() servicequotasiface.ServiceQuotasAPI
	Route53() route53iface.Route53API
	CloudFront() cloudfrontiface.CloudFrontAPI
	CognitoIdentityProvider() cognitoidentityprovideriface.CognitoIdentityProviderAPI
	CloudFormation() cloudformationiface.CloudFormationAPI
	Region() string
	PrincipalARN() string
}
