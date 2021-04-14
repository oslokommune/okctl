// Package cloud provides access to AWS APIs
package cloud

import (
	"fmt"

	credentialspkg "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"

	"github.com/aws/aws-sdk-go/service/cloudfront"

	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	"github.com/aws/aws-sdk-go/aws"
	awsCreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	awsauth "github.com/oslokommune/okctl/pkg/credentials/aws"
)

// Provider stores state required for interacting with the AWS API
type Provider struct {
	Provider v1alpha1.CloudProvider
}

// New returns a new AWS API provider and builds a session from
// the provided authenticator
func New(region string, a awsauth.Authenticator) (*Provider, error) {
	sess, creds, err := NewSession(region, a)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with aws: %w", err)
	}

	return NewFromSession(region, creds.PrincipalARN, sess)
}

// NewFromSession returns a new AWS API provider and builds a session
// from the provided authenticator
func NewFromSession(region, principalARN string, sess *session.Session) (*Provider, error) {
	services := &Services{
		region:       region,
		principalARN: principalARN,
	}
	p := &Provider{
		Provider: services,
	}

	services.sm = secretsmanager.New(sess)
	services.s3 = s3.New(sess)
	services.iam = iam.New(sess)
	services.cfn = cloudformation.New(sess)
	services.ec2 = ec2.New(sess)
	services.elbv2 = elbv2.New(sess)
	services.eks = eks.New(sess)
	services.ssm = ssm.New(sess)
	services.sq = servicequotas.New(sess)
	services.r53 = route53.New(sess)
	services.cip = cognitoidentityprovider.New(sess)
	services.cf = cloudfront.New(sess)

	return p, nil
}

// NewSession returns an AWS session using the provided authenticator
func NewSession(region string, auth awsauth.Authenticator) (*session.Session, *awsauth.Credentials, error) {
	creds, err := auth.Raw()
	if err != nil {
		return nil, nil, err
	}

	config := aws.NewConfig().
		WithRegion(region).
		WithCredentials(
			awsCreds.NewStaticCredentials(
				creds.AccessKeyID,
				creds.SecretAccessKey,
				creds.SessionToken,
			),
		)

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, nil, err
	}

	return sess, creds, err
}

// NewSessionFromEnv returns an initialised session
func NewSessionFromEnv(region string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentialspkg.NewEnvCredentials(),
		Region:      aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving credentials from env: %w", err)
	}

	return sess, nil
}

// Services stores access to the various AWS APIs
type Services struct {
	sm    secretsmanageriface.SecretsManagerAPI
	s3    s3iface.S3API
	iam   iamiface.IAMAPI
	cfn   cloudformationiface.CloudFormationAPI
	ec2   ec2iface.EC2API
	elbv2 elbv2iface.ELBV2API
	eks   eksiface.EKSAPI
	ssm   ssmiface.SSMAPI
	sq    servicequotasiface.ServiceQuotasAPI
	r53   route53iface.Route53API
	cip   cognitoidentityprovideriface.CognitoIdentityProviderAPI
	cf    cloudfrontiface.CloudFrontAPI

	region       string
	principalARN string
}

// SecretsManager returns an interface to the SecretsManager API
func (s *Services) SecretsManager() secretsmanageriface.SecretsManagerAPI {
	return s.sm
}

// S3 returns an interface to the S3 API
func (s *Services) S3() s3iface.S3API {
	return s.s3
}

// IAM returns an interface to the IAM API
func (s *Services) IAM() iamiface.IAMAPI {
	return s.iam
}

// CloudFront returns an interface to the AWS CloudFront API
func (s *Services) CloudFront() cloudfrontiface.CloudFrontAPI {
	return s.cf
}

// CognitoIdentityProvider returns an interface to the AWS Cognito API
func (s *Services) CognitoIdentityProvider() cognitoidentityprovideriface.CognitoIdentityProviderAPI {
	return s.cip
}

// ServiceQuotas returns an interface to AWS ServiceQuota API
func (s *Services) ServiceQuotas() servicequotasiface.ServiceQuotasAPI {
	return s.sq
}

// Route53 returns an interface to the AWS Route53 API
func (s *Services) Route53() route53iface.Route53API {
	return s.r53
}

// SSM returns an interface to the AWS SSM API
func (s *Services) SSM() ssmiface.SSMAPI {
	return s.ssm
}

// EC2 returns an interface to the AWS EC2 API
func (s *Services) EC2() ec2iface.EC2API {
	return s.ec2
}

// EKS returns an interface to the AWS EKS API
func (s *Services) EKS() eksiface.EKSAPI {
	return s.eks
}

// ELBV2 returns an interface to the AWS ELBV2 API
func (s *Services) ELBV2() elbv2iface.ELBV2API {
	return s.elbv2
}

// CloudFormation returns an interface to the AWS CloudFormation API
func (s *Services) CloudFormation() cloudformationiface.CloudFormationAPI {
	return s.cfn
}

// Region returns the configured AWS region
func (s *Services) Region() string {
	return s.region
}

// PrincipalARN return the principal arn of the authenticated party
func (s *Services) PrincipalARN() string {
	return s.principalARN
}
