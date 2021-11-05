// Package mock provides mocks
package mock

import "C"

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"

	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"

	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"

	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"

	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/oslokommune/okctl/pkg/api/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	awspkg "github.com/oslokommune/okctl/pkg/credentials/aws"
)

const (
	// DefaultAccessKeyID is a mocked aws access key id
	DefaultAccessKeyID = "ASIAV3ZUEFP6EXAMPLE"
	// DefaultSecretAccessKey is a mocked aws secret access key
	// nolint: gosec
	DefaultSecretAccessKey = "8P+SQvWIuLnKhh8d++jpw0nNmQRBZvNEXAMPLEKEY"
	// DefaultSessionToken is a mocked aws session token
	// nolint: gosec, lll
	DefaultSessionToken = "IQoJb3JpZ2luX2VjEOz////////////////////wEXAMPLEtMSJHMEUCIDoKK3JH9uG\nQE1z0sINr5M4jk+Na8KHDcCYRVjJCZEvOAiEA3OvJGtw1EcViOleS2vhs8VdCKFJQWP\nQrmGdeehM4IC1NtBmUpp2wUE8phUZampKsburEDy0KPkyQDYwT7WZ0wq5VSXDvp75YU\n9HFvlRd8Tx6q6fE8YQcHNVXAkiY9q6d+xo0rKwT38xVqr7ZD0u0iPPkUL64lIZbqBAz\n+scqKmlzm8FDrypNC9Yjc8fPOLn9FX9KSYvKTr4rvx3iSIlTJabIQwj2ICCR/oLxBA=="
	// DefaultExpiration is a mocked aws expiration
	DefaultExpiration = "2019-11-01T20:26:47Z"
	// DefaultPrincipalARN is a mocked principal arn
	DefaultPrincipalARN = "arn:aws:sts::123456789012:assumed-role/admin/user"
	// DefaultRegion is a mocked default region
	DefaultRegion = "eu-west-1"
	// DefaultStackName is a mocked default stack name
	DefaultStackName = "myStack"
	// DefaultCloudFrontDistributionARN is a mocked default
	DefaultCloudFrontDistributionARN = "arn:::::/distribution/FHH78FAKE"
	// DefaultFargateProfilePodExecutionRoleARN is the default name of the pod execution role
	DefaultFargateProfilePodExecutionRoleARN = "arn:aws:iam::123456789012:role/fargatePodExecutionRole-GHEFFAKE"
	// DefaultListenerARN represents a listener arn
	DefaultListenerARN = "arn:aws:elasticloadbalancing:eu-west-1:123456789012:listener/app/d65FAKE/3f633e2FAKE/cb0f5FAKE"
	// DefaultLoadBalancerARN represents a load balancer arn
	DefaultLoadBalancerARN = "arn:aws:elasticloadbalancing:eu-west-1:123456789012:loadbalancer/app/vbhe933FAKE/145afFAKE"
	// DefaultCertificateARN represents a certificate arn
	DefaultCertificateARN = "arn:aws:acm:eu-west-1:123456789012:certificate/123456789012-1234-1234-1234-12345678"
	// DefaultDomain is the default domain
	DefaultDomain = "test.oslo.systems"
	// DefaultClusterName is the default cluster name
	DefaultClusterName = "test"
	// DefaultTargetGroupARN is the default target group arn
	DefaultTargetGroupARN = "arn:aws:elasticloadbalancing:eu-west-1:123456789012:targetgroup/sfihef303FAKE/69b3FAKE"
	// DefaultVpcID is the default vpc id
	DefaultVpcID = "ih0fj2f03vpcFAKE"
	// DefaultSecurityGroupID is the default security group id
	DefaultSecurityGroupID = "sg-dyf9uf03FAKE"
	// DefaultSecurityGroupName is the default security group name
	DefaultSecurityGroupName = "myCoolSG"
)

// DefaultCredentials returns a mocked set of aws credentials
func DefaultCredentials() *awspkg.Credentials {
	t, _ := time.Parse(time.RFC3339, DefaultExpiration)

	return &awspkg.Credentials{
		AccessKeyID:     DefaultAccessKeyID,
		SecretAccessKey: DefaultSecretAccessKey,
		SessionToken:    DefaultSessionToken,
		SecurityToken:   DefaultSessionToken,
		PrincipalARN:    DefaultPrincipalARN,
		Expires:         t.Local(),
		Region:          DefaultRegion,
	}
}

// DefaultStsCredentials returns a mocked set of aws sts credentials
func DefaultStsCredentials() *sts.Credentials {
	t, _ := time.Parse(time.RFC3339, DefaultExpiration)

	return &sts.Credentials{
		AccessKeyId:     aws.String(DefaultAccessKeyID),
		Expiration:      aws.Time(t.Local()),
		SecretAccessKey: aws.String(DefaultSecretAccessKey),
		SessionToken:    aws.String(DefaultSessionToken),
	}
}

// DefaultValidCredentials returns a mocked set of valid credentials
func DefaultValidCredentials() *awspkg.Credentials {
	creds := DefaultCredentials()

	creds.Expires = time.Now().Add(1 * time.Hour).Local()

	return creds
}

// DefaultValidStsCredentials returns a mocked set of valid aws sts credentials
func DefaultValidStsCredentials() *sts.Credentials {
	creds := DefaultStsCredentials()

	creds.Expiration = aws.Time(time.Now().Add(1 * time.Hour).Local())

	return creds
}

// SQAPI mocks the service quotas interface
type SQAPI struct {
	servicequotasiface.ServiceQuotasAPI

	GetServiceQuotaFn func(*servicequotas.GetServiceQuotaInput) (*servicequotas.GetServiceQuotaOutput, error)
}

// GetServiceQuota invokes the mocked response
func (a *SQAPI) GetServiceQuota(input *servicequotas.GetServiceQuotaInput) (*servicequotas.GetServiceQuotaOutput, error) {
	return a.GetServiceQuotaFn(input)
}

// STSAPI stores state for mocking out the STS API
type STSAPI struct {
	stsiface.STSAPI

	AssumeRoleWithSAMLFn func(*sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error)
}

// AssumeRoleWithSAML invokes the mocked function
func (a *STSAPI) AssumeRoleWithSAML(in *sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
	return a.AssumeRoleWithSAMLFn(in)
}

// NewGoodSTSAPI returns a mocked sts api that will succeed
func NewGoodSTSAPI() stsiface.STSAPI {
	return &STSAPI{
		AssumeRoleWithSAMLFn: func(input *sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
			return &sts.AssumeRoleWithSAMLOutput{
				AssumedRoleUser: &sts.AssumedRoleUser{
					Arn: aws.String(DefaultPrincipalARN),
				},
				Credentials: DefaultStsCredentials(),
			}, nil
		},
	}
}

// EKSAPI mocks eksiface.EKSAPI
type EKSAPI struct {
	eksiface.EKSAPI

	DescribeClusterFn        func(*eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error)
	DescribeFargateProfileFn func(*eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error)
}

// DescribeCluster mocks describe cluster invocation
func (a *EKSAPI) DescribeCluster(input *eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error) {
	return a.DescribeClusterFn(input)
}

// DescribeClusterWithContext mocks describe cluster with context invocation
func (a *EKSAPI) DescribeClusterWithContext(_ aws.Context, input *eks.DescribeClusterInput, _ ...request.Option) (*eks.DescribeClusterOutput, error) {
	return a.DescribeClusterFn(input)
}

// DescribeFargateProfile mocks the API invocation
func (a *EKSAPI) DescribeFargateProfile(input *eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error) {
	return a.DescribeFargateProfileFn(input)
}

// EC2API provides a mocked structure of the ec2 API
type EC2API struct {
	ec2iface.EC2API

	DescribeSubnetsFn               func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeAddressesFn             func(*ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error)
	DescribeInternetGatewaysFn      func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error)
	DescribeVpcsFn                  func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	DescribeSecurityGroupsFn        func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	AuthorizeSecurityGroupIngressFn func(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	RevokeSecurityGroupIngressFn    func(*ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error)
	DeleteSecurityGroupFn           func(*ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error)
}

// DeleteSecurityGroup returns mocked invocation
func (a *EC2API) DeleteSecurityGroup(input *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
	return a.DeleteSecurityGroupFn(input)
}

// DescribeSubnets invokes the mocked response
func (a *EC2API) DescribeSubnets(sub *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return a.DescribeSubnetsFn(sub)
}

// DescribeAddresses invokes the mocked response
func (a *EC2API) DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	return a.DescribeAddressesFn(input)
}

// DescribeInternetGateways invokes the mocked response
func (a *EC2API) DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	return a.DescribeInternetGatewaysFn(input)
}

// DescribeVpcs invokes the mocked response
func (a *EC2API) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return a.DescribeVpcsFn(input)
}

// DescribeSecurityGroups invokes the mocked response
func (a *EC2API) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return a.DescribeSecurityGroupsFn(input)
}

// DescribeSecurityGroupsWithContext invokes the mocked response
func (a *EC2API) DescribeSecurityGroupsWithContext(_ aws.Context, input *ec2.DescribeSecurityGroupsInput, _ ...request.Option) (*ec2.DescribeSecurityGroupsOutput, error) {
	return a.DescribeSecurityGroupsFn(input)
}

func (a *EC2API) DescribeSecurityGroupsWithContext2(_ aws.Context, input *ec2.DescribeSecurityGroupsInput, _ ...request.Option) (*ec2.DescribeSecurityGroupsOutput, error) {
	return a.DescribeSecurityGroupsFn(input)
}

// AuthorizeSecurityGroupIngress invokes the mocked response
func (a *EC2API) AuthorizeSecurityGroupIngress(input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return a.AuthorizeSecurityGroupIngressFn(input)
}

// RevokeSecurityGroupIngress invokes the mocked response
func (a *EC2API) RevokeSecurityGroupIngress(input *ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error) {
	return a.RevokeSecurityGroupIngressFn(input)
}

// CFAPI provides a mocked structure of the cf API
type CFAPI struct {
	cloudformationiface.CloudFormationAPI

	DescribeStacksFn   []func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
	nextDescribeStacks int

	CreateStackFn         func(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStackFn         func(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
	DescribeStackEventsFn func(*cloudformation.DescribeStackEventsInput) (*cloudformation.DescribeStackEventsOutput, error)
}

// DescribeStackEvents invokes a mocked response
func (a *CFAPI) DescribeStackEvents(input *cloudformation.DescribeStackEventsInput) (*cloudformation.DescribeStackEventsOutput, error) {
	return a.DescribeStackEventsFn(input)
}

// DeleteStack invokes a mocked a response
func (a *CFAPI) DeleteStack(in *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
	return a.DeleteStackFn(in)
}

// CreateStack invokes a mocked response
func (a *CFAPI) CreateStack(in *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return a.CreateStackFn(in)
}

// DescribeStacks invokes a mocked response
func (a *CFAPI) DescribeStacks(stack *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	i := a.nextDescribeStacks
	a.nextDescribeStacks++

	return a.DescribeStacksFn[i](stack)
}

// R53API mocks the Route53 API
type R53API struct {
	route53iface.Route53API

	ListHostedZonesFn        func(*route53.ListHostedZonesInput) (*route53.ListHostedZonesOutput, error)
	ListResourceRecordSetsFn func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
}

// ListResourceRecordSets returns a mocked response
func (a *R53API) ListResourceRecordSets(input *route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
	return a.ListResourceRecordSetsFn(input)
}

// ListHostedZones invokes a mocked response
func (a *R53API) ListHostedZones(input *route53.ListHostedZonesInput) (*route53.ListHostedZonesOutput, error) {
	return a.ListHostedZonesFn(input)
}

// CFRONTAPI mocks the CloudFront API
type CFRONTAPI struct {
	cloudfrontiface.CloudFrontAPI

	GetDistributionFn func(*cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error)
}

// GetDistribution invokes a mocked response
func (a *CFRONTAPI) GetDistribution(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
	return a.GetDistributionFn(input)
}

// CIPAPI mocks the CognitoIdentityProvider API
type CIPAPI struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI

	DescribeUserPoolDomainFn func(*cognitoidentityprovider.DescribeUserPoolDomainInput) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error)
}

// DescribeUserPoolDomain invokes a mocked response
func (a *CIPAPI) DescribeUserPoolDomain(input *cognitoidentityprovider.DescribeUserPoolDomainInput) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error) {
	return a.DescribeUserPoolDomainFn(input)
}

type CWAPI struct {
	GetMetricsStatisticsFn func(input *cloudwatch.GetMetricStatisticsInput) (*cloudwatch.GetMetricStatisticsOutput, error)
}

// DeleteAlarms returns a mocked response
func (c *CWAPI) DeleteAlarms(_ *cloudwatch.DeleteAlarmsInput) (*cloudwatch.DeleteAlarmsOutput, error) {
	panic("implement me")
}

// DeleteAlarmsWithContext returns a mocked response
func (c *CWAPI) DeleteAlarmsWithContext(_ aws.Context, _ *cloudwatch.DeleteAlarmsInput, _ ...request.Option) (*cloudwatch.DeleteAlarmsOutput, error) {
	panic("implement me")
}

// DeleteAlarmsRequest returns a mocked response
func (c *CWAPI) DeleteAlarmsRequest(_ *cloudwatch.DeleteAlarmsInput) (*request.Request, *cloudwatch.DeleteAlarmsOutput) {
	panic("implement me")
}

// DeleteAnomalyDetector returns a mocked response
func (c *CWAPI) DeleteAnomalyDetector(_ *cloudwatch.DeleteAnomalyDetectorInput) (*cloudwatch.DeleteAnomalyDetectorOutput, error) {
	panic("implement me")
}

// DeleteAnomalyDetectorWithContext returns a mocked response
func (c *CWAPI) DeleteAnomalyDetectorWithContext(_ aws.Context, _ *cloudwatch.DeleteAnomalyDetectorInput, _ ...request.Option) (*cloudwatch.DeleteAnomalyDetectorOutput, error) {
	panic("implement me")
}

// DeleteAnomalyDetectorRequest returns a mocked response
func (c *CWAPI) DeleteAnomalyDetectorRequest(_ *cloudwatch.DeleteAnomalyDetectorInput) (*request.Request, *cloudwatch.DeleteAnomalyDetectorOutput) {
	panic("implement me")
}

// DeleteDashboards returns a mocked response
func (c *CWAPI) DeleteDashboards(_ *cloudwatch.DeleteDashboardsInput) (*cloudwatch.DeleteDashboardsOutput, error) {
	panic("implement me")
}

// DeleteDashboardsWithContext returns a mocked response
func (c *CWAPI) DeleteDashboardsWithContext(_ aws.Context, _ *cloudwatch.DeleteDashboardsInput, _ ...request.Option) (*cloudwatch.DeleteDashboardsOutput, error) {
	panic("implement me")
}

// DeleteDashboardsRequest returns a mocked response
func (c *CWAPI) DeleteDashboardsRequest(_ *cloudwatch.DeleteDashboardsInput) (*request.Request, *cloudwatch.DeleteDashboardsOutput) {
	panic("implement me")
}

// DeleteInsightRules returns a mocked response
func (c *CWAPI) DeleteInsightRules(_ *cloudwatch.DeleteInsightRulesInput) (*cloudwatch.DeleteInsightRulesOutput, error) {
	panic("implement me")
}

// DeleteInsightRulesWithContext returns a mocked response
func (c *CWAPI) DeleteInsightRulesWithContext(_ aws.Context, _ *cloudwatch.DeleteInsightRulesInput, _ ...request.Option) (*cloudwatch.DeleteInsightRulesOutput, error) {
	panic("implement me")
}

// DeleteInsightRulesRequest returns a mocked response
func (c *CWAPI) DeleteInsightRulesRequest(_ *cloudwatch.DeleteInsightRulesInput) (*request.Request, *cloudwatch.DeleteInsightRulesOutput) {
	panic("implement me")
}

// DeleteMetricStream returns a mocked response
func (c *CWAPI) DeleteMetricStream(_ *cloudwatch.DeleteMetricStreamInput) (*cloudwatch.DeleteMetricStreamOutput, error) {
	panic("implement me")
}

// DeleteMetricStreamWithContext returns a mocked response
func (c *CWAPI) DeleteMetricStreamWithContext(_ aws.Context, _ *cloudwatch.DeleteMetricStreamInput, _ ...request.Option) (*cloudwatch.DeleteMetricStreamOutput, error) {
	panic("implement me")
}

// DeleteMetricStreamRequest returns a mocked response
func (c *CWAPI) DeleteMetricStreamRequest(_ *cloudwatch.DeleteMetricStreamInput) (*request.Request, *cloudwatch.DeleteMetricStreamOutput) {
	panic("implement me")
}

// DescribeAlarmHistory returns a mocked response
func (c *CWAPI) DescribeAlarmHistory(_ *cloudwatch.DescribeAlarmHistoryInput) (*cloudwatch.DescribeAlarmHistoryOutput, error) {
	panic("implement me")
}

// DescribeAlarmHistoryWithContext returns a mocked response
func (c *CWAPI) DescribeAlarmHistoryWithContext(_ aws.Context, _ *cloudwatch.DescribeAlarmHistoryInput, _ ...request.Option) (*cloudwatch.DescribeAlarmHistoryOutput, error) {
	panic("implement me")
}

// DescribeAlarmHistoryRequest returns a mocked response
func (c *CWAPI) DescribeAlarmHistoryRequest(_ *cloudwatch.DescribeAlarmHistoryInput) (*request.Request, *cloudwatch.DescribeAlarmHistoryOutput) {
	panic("implement me")
}

// DescribeAlarmHistoryPages returns a mocked response
func (c *CWAPI) DescribeAlarmHistoryPages(_ *cloudwatch.DescribeAlarmHistoryInput, _ func(*cloudwatch.DescribeAlarmHistoryOutput, bool) bool) error {
	panic("implement me")
}

// DescribeAlarmHistoryPagesWithContext returns a mocked response
func (c *CWAPI) DescribeAlarmHistoryPagesWithContext(_ aws.Context, _ *cloudwatch.DescribeAlarmHistoryInput, _ func(*cloudwatch.DescribeAlarmHistoryOutput, bool) bool, _ ...request.Option) error {
	panic("implement me")
}

// DescribeAlarms returns a mocked response
func (c *CWAPI) DescribeAlarms(_ *cloudwatch.DescribeAlarmsInput) (*cloudwatch.DescribeAlarmsOutput, error) {
	panic("implement me")
}

// DescribeAlarmsWithContext returns a mocked response
func (c *CWAPI) DescribeAlarmsWithContext(_ aws.Context, _ *cloudwatch.DescribeAlarmsInput, _ ...request.Option) (*cloudwatch.DescribeAlarmsOutput, error) {
	panic("implement me")
}

// DescribeAlarmsRequest returns a mocked response
func (c *CWAPI) DescribeAlarmsRequest(_ *cloudwatch.DescribeAlarmsInput) (*request.Request, *cloudwatch.DescribeAlarmsOutput) {
	panic("implement me")
}

// DescribeAlarmsPages returns a mocked response
func (c *CWAPI) DescribeAlarmsPages(_ *cloudwatch.DescribeAlarmsInput, _ func(*cloudwatch.DescribeAlarmsOutput, bool) bool) error {
	panic("implement me")
}

// DescribeAlarmsPagesWithContext returns a mocked response
func (c *CWAPI) DescribeAlarmsPagesWithContext(_ aws.Context, _ *cloudwatch.DescribeAlarmsInput, _ func(*cloudwatch.DescribeAlarmsOutput, bool) bool, _ ...request.Option) error {
	panic("implement me")
}

// DescribeAlarmsForMetric returns a mocked response
func (c *CWAPI) DescribeAlarmsForMetric(_ *cloudwatch.DescribeAlarmsForMetricInput) (*cloudwatch.DescribeAlarmsForMetricOutput, error) {
	panic("implement me")
}

// DescribeAlarmsForMetricWithContext returns a mocked response
func (c *CWAPI) DescribeAlarmsForMetricWithContext(_ aws.Context, _ *cloudwatch.DescribeAlarmsForMetricInput, _ ...request.Option) (*cloudwatch.DescribeAlarmsForMetricOutput, error) {
	panic("implement me")
}

// DescribeAlarmsForMetricRequest returns a mocked response
func (c *CWAPI) DescribeAlarmsForMetricRequest(_ *cloudwatch.DescribeAlarmsForMetricInput) (*request.Request, *cloudwatch.DescribeAlarmsForMetricOutput) {
	panic("implement me")
}

// DescribeAnomalyDetectors returns a mocked response
func (c *CWAPI) DescribeAnomalyDetectors(_ *cloudwatch.DescribeAnomalyDetectorsInput) (*cloudwatch.DescribeAnomalyDetectorsOutput, error) {
	panic("implement me")
}

// DescribeAnomalyDetectorsWithContext returns a mocked response
func (c *CWAPI) DescribeAnomalyDetectorsWithContext(_ aws.Context, _ *cloudwatch.DescribeAnomalyDetectorsInput, _ ...request.Option) (*cloudwatch.DescribeAnomalyDetectorsOutput, error) {
	panic("implement me")
}

// DescribeAnomalyDetectorsRequest returns a mocked response
func (c *CWAPI) DescribeAnomalyDetectorsRequest(_ *cloudwatch.DescribeAnomalyDetectorsInput) (*request.Request, *cloudwatch.DescribeAnomalyDetectorsOutput) {
	panic("implement me")
}

// DescribeInsightRules returns a mocked response
func (c *CWAPI) DescribeInsightRules(_ *cloudwatch.DescribeInsightRulesInput) (*cloudwatch.DescribeInsightRulesOutput, error) {
	panic("implement me")
}

// DescribeInsightRulesWithContext returns a mocked response
func (c *CWAPI) DescribeInsightRulesWithContext(_ aws.Context, _ *cloudwatch.DescribeInsightRulesInput, _ ...request.Option) (*cloudwatch.DescribeInsightRulesOutput, error) {
	panic("implement me")
}

// DescribeInsightRulesRequest returns a mocked response
func (c *CWAPI) DescribeInsightRulesRequest(_ *cloudwatch.DescribeInsightRulesInput) (*request.Request, *cloudwatch.DescribeInsightRulesOutput) {
	panic("implement me")
}

// DescribeInsightRulesPages returns a mocked response
func (c *CWAPI) DescribeInsightRulesPages(_ *cloudwatch.DescribeInsightRulesInput, _ func(*cloudwatch.DescribeInsightRulesOutput, bool) bool) error {
	panic("implement me")
}

// DescribeInsightRulesPagesWithContext returns a mocked response
func (c *CWAPI) DescribeInsightRulesPagesWithContext(_ aws.Context, _ *cloudwatch.DescribeInsightRulesInput, _ func(*cloudwatch.DescribeInsightRulesOutput, bool) bool, _ ...request.Option) error {
	panic("implement me")
}

// DisableAlarmActions returns a mocked response
func (c *CWAPI) DisableAlarmActions(_ *cloudwatch.DisableAlarmActionsInput) (*cloudwatch.DisableAlarmActionsOutput, error) {
	panic("implement me")
}

// DisableAlarmActionsWithContext returns a mocked response
func (c *CWAPI) DisableAlarmActionsWithContext(_ aws.Context, _ *cloudwatch.DisableAlarmActionsInput, _ ...request.Option) (*cloudwatch.DisableAlarmActionsOutput, error) {
	panic("implement me")
}

// DisableAlarmActionsRequest returns a mocked response
func (c *CWAPI) DisableAlarmActionsRequest(_ *cloudwatch.DisableAlarmActionsInput) (*request.Request, *cloudwatch.DisableAlarmActionsOutput) {
	panic("implement me")
}

// DisableInsightRules returns a mocked response
func (c *CWAPI) DisableInsightRules(_ *cloudwatch.DisableInsightRulesInput) (*cloudwatch.DisableInsightRulesOutput, error) {
	panic("implement me")
}

// DisableInsightRulesWithContext returns a mocked response
func (c *CWAPI) DisableInsightRulesWithContext(_ aws.Context, _ *cloudwatch.DisableInsightRulesInput, _ ...request.Option) (*cloudwatch.DisableInsightRulesOutput, error) {
	panic("implement me")
}

// DisableInsightRulesRequest returns a mocked response
func (c *CWAPI) DisableInsightRulesRequest(_ *cloudwatch.DisableInsightRulesInput) (*request.Request, *cloudwatch.DisableInsightRulesOutput) {
	panic("implement me")
}

// EnableAlarmActions returns a mocked response
func (c *CWAPI) EnableAlarmActions(_ *cloudwatch.EnableAlarmActionsInput) (*cloudwatch.EnableAlarmActionsOutput, error) {
	panic("implement me")
}

// EnableAlarmActionsWithContext returns a mocked response
func (c *CWAPI) EnableAlarmActionsWithContext(_ aws.Context, _ *cloudwatch.EnableAlarmActionsInput, _ ...request.Option) (*cloudwatch.EnableAlarmActionsOutput, error) {
	panic("implement me")
}

// EnableAlarmActionsRequest returns a mocked response
func (c *CWAPI) EnableAlarmActionsRequest(_ *cloudwatch.EnableAlarmActionsInput) (*request.Request, *cloudwatch.EnableAlarmActionsOutput) {
	panic("implement me")
}

// EnableInsightRules returns a mocked response
func (c *CWAPI) EnableInsightRules(_ *cloudwatch.EnableInsightRulesInput) (*cloudwatch.EnableInsightRulesOutput, error) {
	panic("implement me")
}

// EnableInsightRulesWithContext returns a mocked response
func (c *CWAPI) EnableInsightRulesWithContext(_ aws.Context, _ *cloudwatch.EnableInsightRulesInput, _ ...request.Option) (*cloudwatch.EnableInsightRulesOutput, error) {
	panic("implement me")
}

// EnableInsightRulesRequest returns a mocked response
func (c *CWAPI) EnableInsightRulesRequest(_ *cloudwatch.EnableInsightRulesInput) (*request.Request, *cloudwatch.EnableInsightRulesOutput) {
	panic("implement me")
}

// GetDashboard returns a mocked response
func (c *CWAPI) GetDashboard(_ *cloudwatch.GetDashboardInput) (*cloudwatch.GetDashboardOutput, error) {
	panic("implement me")
}

// GetDashboardWithContext returns a mocked response
func (c *CWAPI) GetDashboardWithContext(_ aws.Context, _ *cloudwatch.GetDashboardInput, _ ...request.Option) (*cloudwatch.GetDashboardOutput, error) {
	panic("implement me")
}

// GetDashboardRequest returns a mocked response
func (c *CWAPI) GetDashboardRequest(_ *cloudwatch.GetDashboardInput) (*request.Request, *cloudwatch.GetDashboardOutput) {
	panic("implement me")
}

// GetInsightRuleReport returns a mocked response
func (c *CWAPI) GetInsightRuleReport(_ *cloudwatch.GetInsightRuleReportInput) (*cloudwatch.GetInsightRuleReportOutput, error) {
	panic("implement me")
}

// GetInsightRuleReportWithContext returns a mocked response
func (c *CWAPI) GetInsightRuleReportWithContext(_ aws.Context, _ *cloudwatch.GetInsightRuleReportInput, _ ...request.Option) (*cloudwatch.GetInsightRuleReportOutput, error) {
	panic("implement me")
}

// GetInsightRuleReportRequest returns a mocked response
func (c *CWAPI) GetInsightRuleReportRequest(_ *cloudwatch.GetInsightRuleReportInput) (*request.Request, *cloudwatch.GetInsightRuleReportOutput) {
	panic("implement me")
}

// GetMetricData returns a mocked response
func (c *CWAPI) GetMetricData(_ *cloudwatch.GetMetricDataInput) (*cloudwatch.GetMetricDataOutput, error) {
	panic("implement me")
}

// GetMetricDataWithContext returns a mocked response
func (c *CWAPI) GetMetricDataWithContext(_ aws.Context, _ *cloudwatch.GetMetricDataInput, _ ...request.Option) (*cloudwatch.GetMetricDataOutput, error) {
	panic("implement me")
}

// GetMetricDataRequest returns a mocked response
func (c *CWAPI) GetMetricDataRequest(_ *cloudwatch.GetMetricDataInput) (*request.Request, *cloudwatch.GetMetricDataOutput) {
	panic("implement me")
}

// GetMetricDataPages returns a mocked response
func (c *CWAPI) GetMetricDataPages(_ *cloudwatch.GetMetricDataInput, _ func(*cloudwatch.GetMetricDataOutput, bool) bool) error {
	panic("implement me")
}

// GetMetricDataPagesWithContext returns a mocked response
func (c *CWAPI) GetMetricDataPagesWithContext(_ aws.Context, _ *cloudwatch.GetMetricDataInput, _ func(*cloudwatch.GetMetricDataOutput, bool) bool, _ ...request.Option) error {
	panic("implement me")
}

// GetMetricStatistics returns a mocked response
func (c *CWAPI) GetMetricStatistics(input *cloudwatch.GetMetricStatisticsInput) (*cloudwatch.GetMetricStatisticsOutput, error) {
	return c.GetMetricsStatisticsFn(input)
}

// GetMetricStatisticsWithContext returns a mocked response
func (c *CWAPI) GetMetricStatisticsWithContext(_ aws.Context, _ *cloudwatch.GetMetricStatisticsInput, _ ...request.Option) (*cloudwatch.GetMetricStatisticsOutput, error) {
	panic("implement me")
}

// GetMetricStatisticsRequest returns a mocked response
func (c *CWAPI) GetMetricStatisticsRequest(_ *cloudwatch.GetMetricStatisticsInput) (*request.Request, *cloudwatch.GetMetricStatisticsOutput) {
	panic("implement me")
}

// GetMetricStream returns a mocked response
func (c *CWAPI) GetMetricStream(_ *cloudwatch.GetMetricStreamInput) (*cloudwatch.GetMetricStreamOutput, error) {
	panic("implement me")
}

// GetMetricStreamWithContext returns a mocked response
func (c *CWAPI) GetMetricStreamWithContext(_ aws.Context, _ *cloudwatch.GetMetricStreamInput, _ ...request.Option) (*cloudwatch.GetMetricStreamOutput, error) {
	panic("implement me")
}

// GetMetricStreamRequest returns a mocked response
func (c *CWAPI) GetMetricStreamRequest(_ *cloudwatch.GetMetricStreamInput) (*request.Request, *cloudwatch.GetMetricStreamOutput) {
	panic("implement me")
}

// GetMetricWidgetImage returns a mocked response
func (c *CWAPI) GetMetricWidgetImage(_ *cloudwatch.GetMetricWidgetImageInput) (*cloudwatch.GetMetricWidgetImageOutput, error) {
	panic("implement me")
}

// GetMetricWidgetImageWithContext returns a mocked response
func (c *CWAPI) GetMetricWidgetImageWithContext(_ aws.Context, _ *cloudwatch.GetMetricWidgetImageInput, _ ...request.Option) (*cloudwatch.GetMetricWidgetImageOutput, error) {
	panic("implement me")
}

// GetMetricWidgetImageRequest returns a mocked response
func (c *CWAPI) GetMetricWidgetImageRequest(_ *cloudwatch.GetMetricWidgetImageInput) (*request.Request, *cloudwatch.GetMetricWidgetImageOutput) {
	panic("implement me")
}

// ListDashboards returns a mocked response
func (c *CWAPI) ListDashboards(_ *cloudwatch.ListDashboardsInput) (*cloudwatch.ListDashboardsOutput, error) {
	panic("implement me")
}

// ListDashboardsWithContext returns a mocked response
func (c *CWAPI) ListDashboardsWithContext(_ aws.Context, _ *cloudwatch.ListDashboardsInput, _ ...request.Option) (*cloudwatch.ListDashboardsOutput, error) {
	panic("implement me")
}

// ListDashboardsRequest returns a mocked response
func (c *CWAPI) ListDashboardsRequest(_ *cloudwatch.ListDashboardsInput) (*request.Request, *cloudwatch.ListDashboardsOutput) {
	panic("implement me")
}

// ListDashboardsPages returns a mocked response
func (c *CWAPI) ListDashboardsPages(_ *cloudwatch.ListDashboardsInput, _ func(*cloudwatch.ListDashboardsOutput, bool) bool) error {
	panic("implement me")
}

// ListDashboardsPagesWithContext returns a mocked response
func (c *CWAPI) ListDashboardsPagesWithContext(_ aws.Context, _ *cloudwatch.ListDashboardsInput, _ func(*cloudwatch.ListDashboardsOutput, bool) bool, _ ...request.Option) error {
	panic("implement me")
}

// ListMetricStreams returns a mocked response
func (c *CWAPI) ListMetricStreams(_ *cloudwatch.ListMetricStreamsInput) (*cloudwatch.ListMetricStreamsOutput, error) {
	panic("implement me")
}

// ListMetricStreamsWithContext returns a mocked response
func (c *CWAPI) ListMetricStreamsWithContext(_ aws.Context, _ *cloudwatch.ListMetricStreamsInput, _ ...request.Option) (*cloudwatch.ListMetricStreamsOutput, error) {
	panic("implement me")
}

// ListMetricStreamsRequest returns a mocked response
func (c *CWAPI) ListMetricStreamsRequest(_ *cloudwatch.ListMetricStreamsInput) (*request.Request, *cloudwatch.ListMetricStreamsOutput) {
	panic("implement me")
}

// ListMetricStreamsPages returns a mocked response
func (c *CWAPI) ListMetricStreamsPages(_ *cloudwatch.ListMetricStreamsInput, _ func(*cloudwatch.ListMetricStreamsOutput, bool) bool) error {
	panic("implement me")
}

// ListMetricStreamsPagesWithContext returns a mocked response
func (c *CWAPI) ListMetricStreamsPagesWithContext(_ aws.Context, _ *cloudwatch.ListMetricStreamsInput, _ func(*cloudwatch.ListMetricStreamsOutput, bool) bool, _ ...request.Option) error {
	panic("implement me")
}

// ListMetrics returns a mocked response
func (c *CWAPI) ListMetrics(_ *cloudwatch.ListMetricsInput) (*cloudwatch.ListMetricsOutput, error) {
	panic("implement me")
}

// ListMetricsWithContext returns a mocked response
func (c *CWAPI) ListMetricsWithContext(_ aws.Context, _ *cloudwatch.ListMetricsInput, _ ...request.Option) (*cloudwatch.ListMetricsOutput, error) {
	panic("implement me")
}

// ListMetricsRequest returns a mocked response
func (c *CWAPI) ListMetricsRequest(_ *cloudwatch.ListMetricsInput) (*request.Request, *cloudwatch.ListMetricsOutput) {
	panic("implement me")
}

// ListMetricsPages returns a mocked response
func (c *CWAPI) ListMetricsPages(_ *cloudwatch.ListMetricsInput, _ func(*cloudwatch.ListMetricsOutput, bool) bool) error {
	panic("implement me")
}

// ListMetricsPagesWithContext returns a mocked response
func (c *CWAPI) ListMetricsPagesWithContext(_ aws.Context, _ *cloudwatch.ListMetricsInput, _ func(*cloudwatch.ListMetricsOutput, bool) bool, _ ...request.Option) error {
	panic("implement me")
}

// ListTagsForResource returns a mocked response
func (c *CWAPI) ListTagsForResource(_ *cloudwatch.ListTagsForResourceInput) (*cloudwatch.ListTagsForResourceOutput, error) {
	panic("implement me")
}

// ListTagsForResourceWithContext returns a mocked response
func (c *CWAPI) ListTagsForResourceWithContext(_ aws.Context, _ *cloudwatch.ListTagsForResourceInput, _ ...request.Option) (*cloudwatch.ListTagsForResourceOutput, error) {
	panic("implement me")
}

// ListTagsForResourceRequest returns a mocked response
func (c *CWAPI) ListTagsForResourceRequest(_ *cloudwatch.ListTagsForResourceInput) (*request.Request, *cloudwatch.ListTagsForResourceOutput) {
	panic("implement me")
}

// PutAnomalyDetector returns a mocked response
func (c *CWAPI) PutAnomalyDetector(_ *cloudwatch.PutAnomalyDetectorInput) (*cloudwatch.PutAnomalyDetectorOutput, error) {
	panic("implement me")
}

// PutAnomalyDetectorWithContext returns a mocked response
func (c *CWAPI) PutAnomalyDetectorWithContext(_ aws.Context, _ *cloudwatch.PutAnomalyDetectorInput, _ ...request.Option) (*cloudwatch.PutAnomalyDetectorOutput, error) {
	panic("implement me")
}

// PutAnomalyDetectorRequest returns a mocked response
func (c *CWAPI) PutAnomalyDetectorRequest(_ *cloudwatch.PutAnomalyDetectorInput) (*request.Request, *cloudwatch.PutAnomalyDetectorOutput) {
	panic("implement me")
}

// PutCompositeAlarm returns a mocked response
func (c *CWAPI) PutCompositeAlarm(_ *cloudwatch.PutCompositeAlarmInput) (*cloudwatch.PutCompositeAlarmOutput, error) {
	panic("implement me")
}

// PutCompositeAlarmWithContext returns a mocked response
func (c *CWAPI) PutCompositeAlarmWithContext(_ aws.Context, _ *cloudwatch.PutCompositeAlarmInput, _ ...request.Option) (*cloudwatch.PutCompositeAlarmOutput, error) {
	panic("implement me")
}

// PutCompositeAlarmRequest returns a mocked response
func (c *CWAPI) PutCompositeAlarmRequest(_ *cloudwatch.PutCompositeAlarmInput) (*request.Request, *cloudwatch.PutCompositeAlarmOutput) {
	panic("implement me")
}

// PutDashboard returns a mocked response
func (c *CWAPI) PutDashboard(_ *cloudwatch.PutDashboardInput) (*cloudwatch.PutDashboardOutput, error) {
	panic("implement me")
}

// PutDashboardWithContext returns a mocked response
func (c *CWAPI) PutDashboardWithContext(_ aws.Context, _ *cloudwatch.PutDashboardInput, _ ...request.Option) (*cloudwatch.PutDashboardOutput, error) {
	panic("implement me")
}

// PutDashboardRequest returns a mocked response
func (c *CWAPI) PutDashboardRequest(_ *cloudwatch.PutDashboardInput) (*request.Request, *cloudwatch.PutDashboardOutput) {
	panic("implement me")
}

// PutInsightRule returns a mocked response
func (c *CWAPI) PutInsightRule(_ *cloudwatch.PutInsightRuleInput) (*cloudwatch.PutInsightRuleOutput, error) {
	panic("implement me")
}

// PutInsightRuleWithContext returns a mocked response
func (c *CWAPI) PutInsightRuleWithContext(_ aws.Context, _ *cloudwatch.PutInsightRuleInput, _ ...request.Option) (*cloudwatch.PutInsightRuleOutput, error) {
	panic("implement me")
}

// PutInsightRuleRequest returns a mocked response
func (c *CWAPI) PutInsightRuleRequest(_ *cloudwatch.PutInsightRuleInput) (*request.Request, *cloudwatch.PutInsightRuleOutput) {
	panic("implement me")
}

// PutMetricAlarm returns a mocked response
func (c *CWAPI) PutMetricAlarm(_ *cloudwatch.PutMetricAlarmInput) (*cloudwatch.PutMetricAlarmOutput, error) {
	panic("implement me")
}

// PutMetricAlarmWithContext returns a mocked response
func (c *CWAPI) PutMetricAlarmWithContext(_ aws.Context, _ *cloudwatch.PutMetricAlarmInput, _ ...request.Option) (*cloudwatch.PutMetricAlarmOutput, error) {
	panic("implement me")
}

// PutMetricAlarmRequest returns a mocked response
func (c *CWAPI) PutMetricAlarmRequest(_ *cloudwatch.PutMetricAlarmInput) (*request.Request, *cloudwatch.PutMetricAlarmOutput) {
	panic("implement me")
}

// PutMetricData returns a mocked response
func (c *CWAPI) PutMetricData(_ *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error) {
	panic("implement me")
}

// PutMetricDataWithContext returns a mocked response
func (c *CWAPI) PutMetricDataWithContext(_ aws.Context, _ *cloudwatch.PutMetricDataInput, _ ...request.Option) (*cloudwatch.PutMetricDataOutput, error) {
	panic("implement me")
}

// PutMetricDataRequest returns a mocked response
func (c *CWAPI) PutMetricDataRequest(_ *cloudwatch.PutMetricDataInput) (*request.Request, *cloudwatch.PutMetricDataOutput) {
	panic("implement me")
}

// PutMetricStream returns a mocked response
func (c *CWAPI) PutMetricStream(_ *cloudwatch.PutMetricStreamInput) (*cloudwatch.PutMetricStreamOutput, error) {
	panic("implement me")
}

// PutMetricStreamWithContext returns a mocked response
func (c *CWAPI) PutMetricStreamWithContext(_ aws.Context, _ *cloudwatch.PutMetricStreamInput, _ ...request.Option) (*cloudwatch.PutMetricStreamOutput, error) {
	panic("implement me")
}

// PutMetricStreamRequest returns a mocked response
func (c *CWAPI) PutMetricStreamRequest(_ *cloudwatch.PutMetricStreamInput) (*request.Request, *cloudwatch.PutMetricStreamOutput) {
	panic("implement me")
}

// SetAlarmState returns a mocked response
func (c *CWAPI) SetAlarmState(_ *cloudwatch.SetAlarmStateInput) (*cloudwatch.SetAlarmStateOutput, error) {
	panic("implement me")
}

// SetAlarmStateWithContext returns a mocked response
func (c *CWAPI) SetAlarmStateWithContext(_ aws.Context, _ *cloudwatch.SetAlarmStateInput, _ ...request.Option) (*cloudwatch.SetAlarmStateOutput, error) {
	panic("implement me")
}

// SetAlarmStateRequest returns a mocked response
func (c *CWAPI) SetAlarmStateRequest(_ *cloudwatch.SetAlarmStateInput) (*request.Request, *cloudwatch.SetAlarmStateOutput) {
	panic("implement me")
}

// StartMetricStreams returns a mocked response
func (c *CWAPI) StartMetricStreams(_ *cloudwatch.StartMetricStreamsInput) (*cloudwatch.StartMetricStreamsOutput, error) {
	panic("implement me")
}

// StartMetricStreamsWithContext returns a mocked response
func (c *CWAPI) StartMetricStreamsWithContext(_ aws.Context, _ *cloudwatch.StartMetricStreamsInput, _ ...request.Option) (*cloudwatch.StartMetricStreamsOutput, error) {
	panic("implement me")
}

// StartMetricStreamsRequest returns a mocked response
func (c *CWAPI) StartMetricStreamsRequest(_ *cloudwatch.StartMetricStreamsInput) (*request.Request, *cloudwatch.StartMetricStreamsOutput) {
	panic("implement me")
}

// StopMetricStreams returns a mocked response
func (c *CWAPI) StopMetricStreams(_ *cloudwatch.StopMetricStreamsInput) (*cloudwatch.StopMetricStreamsOutput, error) {
	panic("implement me")
}

// StopMetricStreamsWithContext returns a mocked response
func (c *CWAPI) StopMetricStreamsWithContext(_ aws.Context, _ *cloudwatch.StopMetricStreamsInput, _ ...request.Option) (*cloudwatch.StopMetricStreamsOutput, error) {
	panic("implement me")
}

// StopMetricStreamsRequest returns a mocked response
func (c *CWAPI) StopMetricStreamsRequest(_ *cloudwatch.StopMetricStreamsInput) (*request.Request, *cloudwatch.StopMetricStreamsOutput) {
	panic("implement me")
}

// TagResource returns a mocked response
func (c *CWAPI) TagResource(_ *cloudwatch.TagResourceInput) (*cloudwatch.TagResourceOutput, error) {
	panic("implement me")
}

// TagResourceWithContext returns a mocked response
func (c *CWAPI) TagResourceWithContext(_ aws.Context, _ *cloudwatch.TagResourceInput, _ ...request.Option) (*cloudwatch.TagResourceOutput, error) {
	panic("implement me")
}

// TagResourceRequest returns a mocked response
func (c *CWAPI) TagResourceRequest(_ *cloudwatch.TagResourceInput) (*request.Request, *cloudwatch.TagResourceOutput) {
	panic("implement me")
}

// UntagResource returns a mocked response
func (c *CWAPI) UntagResource(_ *cloudwatch.UntagResourceInput) (*cloudwatch.UntagResourceOutput, error) {
	panic("implement me")
}

// UntagResourceWithContext returns a mocked response
func (c *CWAPI) UntagResourceWithContext(_ aws.Context, _ *cloudwatch.UntagResourceInput, _ ...request.Option) (*cloudwatch.UntagResourceOutput, error) {
	panic("implement me")
}

// UntagResourceRequest returns a mocked response
func (c *CWAPI) UntagResourceRequest(_ *cloudwatch.UntagResourceInput) (*request.Request, *cloudwatch.UntagResourceOutput) {
	panic("implement me")
}

// WaitUntilAlarmExists returns a mocked response
func (c *CWAPI) WaitUntilAlarmExists(_ *cloudwatch.DescribeAlarmsInput) error {
	panic("implement me")
}

// WaitUntilAlarmExistsWithContext returns a mocked response
func (c *CWAPI) WaitUntilAlarmExistsWithContext(_ aws.Context, _ *cloudwatch.DescribeAlarmsInput, _ ...request.WaiterOption) error {
	panic("implement me")
}

// WaitUntilCompositeAlarmExists returns a mocked response
func (c *CWAPI) WaitUntilCompositeAlarmExists(_ *cloudwatch.DescribeAlarmsInput) error {
	panic("implement me")
}

// WaitUntilCompositeAlarmExistsWithContext returns a mocked response
func (c *CWAPI) WaitUntilCompositeAlarmExistsWithContext(_ aws.Context, _ *cloudwatch.DescribeAlarmsInput, _ ...request.WaiterOption) error {
	panic("implement me")
}

// IAMAPI mocks the IAM API
type IAMAPI struct {
	iamiface.IAMAPI

	AttachRolePolicyFn func(*iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error)
	DetachRolePolicyFn func(*iam.DetachRolePolicyInput) (*iam.DetachRolePolicyOutput, error)
}

// AttachRolePolicy mocks the invocation
func (a *IAMAPI) AttachRolePolicy(input *iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
	return a.AttachRolePolicyFn(input)
}

// DetachRolePolicy mocks the invocation
func (a *IAMAPI) DetachRolePolicy(input *iam.DetachRolePolicyInput) (*iam.DetachRolePolicyOutput, error) {
	return a.DetachRolePolicyFn(input)
}

// S3API mocks the S3 API
type S3API struct {
	s3iface.S3API

	PutObjectFn    func(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObjectFn func(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	GetObjectFn    func(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

// DeleteObject mocks the invocation
func (a *S3API) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return a.DeleteObjectFn(input)
}

// PutObject mocks the invocation
func (a *S3API) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return a.PutObjectFn(input)
}

// GetObject mocks the invocation
func (a *S3API) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return a.GetObjectFn(input)
}

// SMAPI mocks the SecretsManager API
type SMAPI struct {
	secretsmanageriface.SecretsManagerAPI

	RotateSecretFn       func(*secretsmanager.RotateSecretInput) (*secretsmanager.RotateSecretOutput, error)
	CancelRotateSecretFn func(*secretsmanager.CancelRotateSecretInput) (*secretsmanager.CancelRotateSecretOutput, error)
}

// RotateSecret mocks the invocation
func (a *SMAPI) RotateSecret(input *secretsmanager.RotateSecretInput) (*secretsmanager.RotateSecretOutput, error) {
	return a.RotateSecretFn(input)
}

// CancelRotateSecret mocks the invocation
func (a *SMAPI) CancelRotateSecret(input *secretsmanager.CancelRotateSecretInput) (*secretsmanager.CancelRotateSecretOutput, error) {
	return a.CancelRotateSecretFn(input)
}

// DescribeStackEventsSuccess sets a success response on the describe event
func (p *CloudProvider) DescribeStackEventsSuccess() *CloudProvider {
	p.CFAPI.DescribeStackEventsFn = func(input *cloudformation.DescribeStackEventsInput) (*cloudformation.DescribeStackEventsOutput, error) {
		return &cloudformation.DescribeStackEventsOutput{
			StackEvents: []*cloudformation.StackEvent{
				{
					ResourceStatus:       aws.String(cloudformation.ResourceStatusCreateFailed),
					ResourceStatusReason: aws.String("something went wrong"),
					ResourceType:         aws.String("ec2"),
				},
			},
		}, nil
	}

	return p
}

// ACMAPI mocks the ACMAPI interface
type ACMAPI struct {
	acmiface.ACMAPI

	DescribeCertificateFn func(*acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error)
	ListCertificatesFn    func(*acm.ListCertificatesInput) (*acm.ListCertificatesOutput, error)
}

// ListCertificates returns mocked invocation
func (a *ACMAPI) ListCertificates(input *acm.ListCertificatesInput) (*acm.ListCertificatesOutput, error) {
	return a.ListCertificatesFn(input)
}

// DescribeCertificate mocks the invocation
func (a *ACMAPI) DescribeCertificate(input *acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error) {
	return a.DescribeCertificateFn(input)
}

// ELBv2API mocks the ELBv2 API
type ELBv2API struct {
	elbv2iface.ELBV2API

	DescribeListenersFn     func(*elbv2.DescribeListenersInput) (*elbv2.DescribeListenersOutput, error)
	DeleteListenerFn        func(*elbv2.DeleteListenerInput) (*elbv2.DeleteListenerOutput, error)
	DescribeTargetGroupsFn  func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error)
	DescribeTagsFn          func(*elbv2.DescribeTagsInput) (*elbv2.DescribeTagsOutput, error)
	DeleteTargetGroupFn     func(*elbv2.DeleteTargetGroupInput) (*elbv2.DeleteTargetGroupOutput, error)
	DescribeLoadBalancersFn func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error)
	DeleteLoadBalancerFn    func(*elbv2.DeleteLoadBalancerInput) (*elbv2.DeleteLoadBalancerOutput, error)
}

// DeleteLoadBalancer returns mocked invocation
func (a *ELBv2API) DeleteLoadBalancer(input *elbv2.DeleteLoadBalancerInput) (*elbv2.DeleteLoadBalancerOutput, error) {
	return a.DeleteLoadBalancerFn(input)
}

// DescribeLoadBalancers returns mocked invocation
func (a *ELBv2API) DescribeLoadBalancers(input *elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	return a.DescribeLoadBalancersFn(input)
}

// DeleteTargetGroup returns mocked invocation
func (a *ELBv2API) DeleteTargetGroup(input *elbv2.DeleteTargetGroupInput) (*elbv2.DeleteTargetGroupOutput, error) {
	return a.DeleteTargetGroupFn(input)
}

// DescribeTags returns mocked invocation
func (a *ELBv2API) DescribeTags(input *elbv2.DescribeTagsInput) (*elbv2.DescribeTagsOutput, error) {
	return a.DescribeTagsFn(input)
}

// DescribeTargetGroups returns mocked invocation
func (a *ELBv2API) DescribeTargetGroups(input *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	return a.DescribeTargetGroupsFn(input)
}

// DeleteListener returns the mocked invocation
func (a *ELBv2API) DeleteListener(input *elbv2.DeleteListenerInput) (*elbv2.DeleteListenerOutput, error) {
	return a.DeleteListenerFn(input)
}

// DescribeListeners returns the mocked invocation
func (a *ELBv2API) DescribeListeners(input *elbv2.DescribeListenersInput) (*elbv2.DescribeListenersOutput, error) {
	return a.DescribeListenersFn(input)
}

// DeleteStackSuccess sets a success response on the mocked DeleteStack function
func (p *CloudProvider) DeleteStackSuccess() *CloudProvider {
	p.CFAPI.DeleteStackFn = func(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
		return &cloudformation.DeleteStackOutput{}, nil
	}

	return p
}

// CreateStackSuccess sets a success response on the mocked CreateStack function
func (p *CloudProvider) CreateStackSuccess() *CloudProvider {
	p.CFAPI.CreateStackFn = func(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
		return &cloudformation.CreateStackOutput{
			StackId: input.StackName,
		}, nil
	}

	return p
}

// DescribeStacksEmpty pushes an empty response onto the describe stacks list
func (p *CloudProvider) DescribeStacksEmpty() *CloudProvider {
	p.CFAPI.DescribeStacksFn = append(p.CFAPI.DescribeStacksFn, func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
		return nil, awserr.New(
			"ValidationError",
			fmt.Sprintf("Stack with id %s does not exist", DefaultStackName),
			fmt.Errorf("something"),
		)
	})

	return p
}

// DescribeStacksResponse pushes a success response onto the describe stacks list
func (p *CloudProvider) DescribeStacksResponse(status string) *CloudProvider {
	p.CFAPI.DescribeStacksFn = append(p.CFAPI.DescribeStacksFn, func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
		return &cloudformation.DescribeStacksOutput{
			Stacks: Stacks(status),
		}, nil
	})

	return p
}

// DescribeStacksError pushes an unsuccessful response onto the describe stacks list
func (p *CloudProvider) DescribeStacksError(err error) *CloudProvider {
	p.CFAPI.DescribeStacksFn = append(p.CFAPI.DescribeStacksFn, func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
		return nil, err
	})

	return p
}

// DescribeSecurityGroupsResponse returns the provided outputs
func (p *CloudProvider) DescribeSecurityGroupsResponse(output *ec2.DescribeSecurityGroupsOutput, err error) *CloudProvider {
	p.EC2API.DescribeSecurityGroupsFn = func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
		return output, err
	}

	return p
}

// AuthorizeSecurityGroupIngressResponse returns the provided outputs
func (p *CloudProvider) AuthorizeSecurityGroupIngressResponse(output *ec2.AuthorizeSecurityGroupIngressOutput, err error) *CloudProvider {
	p.EC2API.AuthorizeSecurityGroupIngressFn = func(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
		return output, err
	}

	return p
}

// RevokeSecurityGroupIngressResponse returns the provided outputs
func (p *CloudProvider) RevokeSecurityGroupIngressResponse(output *ec2.RevokeSecurityGroupIngressOutput, err error) *CloudProvider {
	p.EC2API.RevokeSecurityGroupIngressFn = func(*ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error) {
		return output, err
	}

	return p
}

// CloudProvider provides a structure for the mocked CloudProvider
type CloudProvider struct {
	v1alpha1.CloudProvider

	SMAPI     *SMAPI
	S3API     *S3API
	IAMAPI    *IAMAPI
	EC2API    *EC2API
	CFAPI     *CFAPI
	EKSAPI    *EKSAPI
	R53API    *R53API
	CFRONTAPI *CFRONTAPI
	CIPAPI    *CIPAPI
	SQAPI     *SQAPI
	ACMAPI    *ACMAPI
	ELBv2API  *ELBv2API
	CWAPI     *CWAPI
}

// SecretsManager returns the mocked SecretsManager API
func (p *CloudProvider) SecretsManager() secretsmanageriface.SecretsManagerAPI {
	return p.SMAPI
}

// S3 returns the mocked S3 API
func (p *CloudProvider) S3() s3iface.S3API {
	return p.S3API
}

// IAM returns the mocked IAM API
func (p *CloudProvider) IAM() iamiface.IAMAPI {
	return p.IAMAPI
}

// ServiceQuotas returns the mocked SQ API
func (p *CloudProvider) ServiceQuotas() servicequotasiface.ServiceQuotasAPI {
	return p.SQAPI
}

// CloudFormation returns the mocked CF API
func (p *CloudProvider) CloudFormation() cloudformationiface.CloudFormationAPI {
	return p.CFAPI
}

// CloudWatch returns the mocked CW API
func (p *CloudProvider) CloudWatch() cloudwatchiface.CloudWatchAPI {
	return p.CWAPI
}

// EC2 returns the mocked EC2 API
func (p *CloudProvider) EC2() ec2iface.EC2API {
	return p.EC2API
}

// EKS returns the mocked EKS API
func (p *CloudProvider) EKS() eksiface.EKSAPI {
	return p.EKSAPI
}

// Route53 returns the mocked Route53 API
func (p *CloudProvider) Route53() route53iface.Route53API {
	return p.R53API
}

// CognitoIdentityProvider returns the mocked CognitoIdentityProvider API
func (p *CloudProvider) CognitoIdentityProvider() cognitoidentityprovideriface.CognitoIdentityProviderAPI {
	return p.CIPAPI
}

// CloudFront returns the mocked CloudFront API
func (p *CloudProvider) CloudFront() cloudfrontiface.CloudFrontAPI {
	return p.CFRONTAPI
}

// ACM returns the mocked ACM API
func (p *CloudProvider) ACM() acmiface.ACMAPI {
	return p.ACMAPI
}

// ELBV2 returns the mocked AWS ELBv2 API
func (p *CloudProvider) ELBV2() elbv2iface.ELBV2API {
	return p.ELBv2API
}

// PrincipalARN mocks the principal arn
func (p *CloudProvider) PrincipalARN() string {
	return "arn:::::/someuser"
}

// DescribeListenersResponse returns the provided values
func (p *CloudProvider) DescribeListenersResponse(output *elbv2.DescribeListenersOutput, err error) *CloudProvider {
	p.ELBv2API.DescribeListenersFn = func(*elbv2.DescribeListenersInput) (*elbv2.DescribeListenersOutput, error) {
		return output, err
	}

	return p
}

// NewCloudProvider returns a mocked cloud provider with no mocks sets
func NewCloudProvider() *CloudProvider {
	return &CloudProvider{
		SMAPI:     &SMAPI{},
		S3API:     &S3API{},
		IAMAPI:    &IAMAPI{},
		EC2API:    &EC2API{},
		CFAPI:     &CFAPI{},
		EKSAPI:    &EKSAPI{},
		R53API:    &R53API{},
		CFRONTAPI: &CFRONTAPI{},
		CIPAPI:    &CIPAPI{},
		SQAPI:     &SQAPI{},
		ACMAPI:    &ACMAPI{},
		ELBv2API:  &ELBv2API{},
		CWAPI:     &CWAPI{},
	}
}

// NewGoodCloudProvider returns a mocked cloud provider with success set on all
// nolint: funlen
func NewGoodCloudProvider() *CloudProvider {
	return &CloudProvider{
		EC2API: &EC2API{
			EC2API: nil,
			DescribeSubnetsFn: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
				return &ec2.DescribeSubnetsOutput{
					Subnets: Subnets(),
				}, nil
			},
			DescribeAddressesFn: func(*ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
				return &ec2.DescribeAddressesOutput{
					Addresses: []*ec2.Address{
						{},
					},
				}, nil
			},
			DescribeInternetGatewaysFn: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
				return &ec2.DescribeInternetGatewaysOutput{
					InternetGateways: []*ec2.InternetGateway{
						{},
					},
				}, nil
			},
			DescribeVpcsFn: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
				return &ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{},
					},
				}, nil
			},
			DescribeSecurityGroupsFn: func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
				return &ec2.DescribeSecurityGroupsOutput{
					SecurityGroups: []*ec2.SecurityGroup{
						{
							GroupId:   aws.String(DefaultSecurityGroupID),
							GroupName: aws.String(DefaultSecurityGroupName),
							VpcId:     aws.String(DefaultVpcID),
						},
					},
				}, nil
			},
			AuthorizeSecurityGroupIngressFn: nil,
			RevokeSecurityGroupIngressFn:    nil,
			DeleteSecurityGroupFn: func(*ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
				return &ec2.DeleteSecurityGroupOutput{}, nil
			},
		},
		EKSAPI: &EKSAPI{
			DescribeClusterFn: func(*eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error) {
				return &eks.DescribeClusterOutput{
					Cluster: &eks.Cluster{
						Arn: aws.String("arn:::something"),
						CertificateAuthority: &eks.Certificate{
							Data: aws.String(base64.StdEncoding.EncodeToString([]byte("something"))),
						},
						Endpoint: aws.String("https://something"),
						Name:     aws.String(mock.DefaultClusterName),
						Status:   aws.String(eks.ClusterStatusActive),
						ResourcesVpcConfig: &eks.VpcConfigResponse{
							ClusterSecurityGroupId: aws.String(mock.DefaultClusterSecurityGroupIDValue),
						},
					},
				}, nil
			},
			DescribeFargateProfileFn: func(*eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error) {
				return &eks.DescribeFargateProfileOutput{
					FargateProfile: &eks.FargateProfile{
						PodExecutionRoleArn: aws.String(DefaultFargateProfilePodExecutionRoleARN),
					},
				}, nil
			},
		},
		R53API: &R53API{
			ListHostedZonesFn: func(*route53.ListHostedZonesInput) (*route53.ListHostedZonesOutput, error) {
				return &route53.ListHostedZonesOutput{
					HostedZones: []*route53.HostedZone{
						{
							Config: &route53.HostedZoneConfig{
								PrivateZone: aws.Bool(false),
							},
							Id:   aws.String("/hostedzone/AABBCCDD"),
							Name: aws.String("test.oslo.systems."),
						},
					},
				}, nil
			},
			ListResourceRecordSetsFn: func(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
				return &route53.ListResourceRecordSetsOutput{
					IsTruncated: aws.Bool(false),
					ResourceRecordSets: []*route53.ResourceRecordSet{
						{
							Name: aws.String("mine.oslo.systems"),
							ResourceRecords: []*route53.ResourceRecord{
								{
									Value: aws.String("ns1.something.com"),
								},
							},
							Type: aws.String("NS"),
						},
					},
				}, nil
			},
		},
		CFRONTAPI: &CFRONTAPI{
			GetDistributionFn: func(*cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
				return &cloudfront.GetDistributionOutput{
					Distribution: &cloudfront.Distribution{
						ARN:        aws.String(DefaultCloudFrontDistributionARN),
						DomainName: aws.String("cloudfront-us-east-1.something.aws.com"),
					},
				}, nil
			},
		},
		CIPAPI: &CIPAPI{
			DescribeUserPoolDomainFn: func(*cognitoidentityprovider.DescribeUserPoolDomainInput) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error) {
				return &cognitoidentityprovider.DescribeUserPoolDomainOutput{
					DomainDescription: &cognitoidentityprovider.DomainDescriptionType{
						CloudFrontDistribution: aws.String("cloudfront-us-east-1.something.aws.com"),
						Domain:                 aws.String("auth.oslo.systems"),
					},
				}, nil
			},
		},
		SQAPI: &SQAPI{
			GetServiceQuotaFn: func(*servicequotas.GetServiceQuotaInput) (*servicequotas.GetServiceQuotaOutput, error) {
				return &servicequotas.GetServiceQuotaOutput{
					Quota: &servicequotas.ServiceQuota{
						Value: aws.Float64(3), // nolint: gomnd
					},
				}, nil
			},
		},
		IAMAPI: &IAMAPI{
			AttachRolePolicyFn: func(*iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
				return &iam.AttachRolePolicyOutput{}, nil
			},
			DetachRolePolicyFn: func(*iam.DetachRolePolicyInput) (*iam.DetachRolePolicyOutput, error) {
				return &iam.DetachRolePolicyOutput{}, nil
			},
		},
		S3API: &S3API{
			PutObjectFn: func(*s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
			DeleteObjectFn: func(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			},
			GetObjectFn: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, nil
			},
		},
		SMAPI: &SMAPI{
			RotateSecretFn: func(*secretsmanager.RotateSecretInput) (*secretsmanager.RotateSecretOutput, error) {
				return &secretsmanager.RotateSecretOutput{}, nil
			},
			CancelRotateSecretFn: func(*secretsmanager.CancelRotateSecretInput) (*secretsmanager.CancelRotateSecretOutput, error) {
				return &secretsmanager.CancelRotateSecretOutput{}, nil
			},
		},
		ACMAPI: &ACMAPI{
			DescribeCertificateFn: func(*acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error) {
				return &acm.DescribeCertificateOutput{
					Certificate: &acm.CertificateDetail{
						CertificateArn: aws.String(DefaultCertificateARN),
						InUseBy: []*string{
							aws.String(DefaultLoadBalancerARN),
						},
					},
				}, nil
			},
			ListCertificatesFn: func(*acm.ListCertificatesInput) (*acm.ListCertificatesOutput, error) {
				return &acm.ListCertificatesOutput{
					CertificateSummaryList: []*acm.CertificateSummary{
						{
							CertificateArn: aws.String(DefaultCertificateARN),
							DomainName:     aws.String(DefaultDomain),
						},
					},
					NextToken: nil,
				}, nil
			},
		},
		ELBv2API: &ELBv2API{
			DescribeListenersFn: func(*elbv2.DescribeListenersInput) (*elbv2.DescribeListenersOutput, error) {
				return &elbv2.DescribeListenersOutput{
					Listeners: []*elbv2.Listener{
						{
							Certificates: []*elbv2.Certificate{
								{
									CertificateArn: aws.String(DefaultCertificateARN),
								},
							},
							ListenerArn: aws.String(DefaultListenerARN),
						},
					},
				}, nil
			},
			DeleteListenerFn: func(*elbv2.DeleteListenerInput) (*elbv2.DeleteListenerOutput, error) {
				return &elbv2.DeleteListenerOutput{}, nil
			},
			DescribeTargetGroupsFn: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
				return &elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{
							TargetGroupArn: aws.String(DefaultTargetGroupARN),
						},
					},
				}, nil
			},
			DescribeTagsFn: func(*elbv2.DescribeTagsInput) (*elbv2.DescribeTagsOutput, error) {
				return &elbv2.DescribeTagsOutput{
					TagDescriptions: []*elbv2.TagDescription{
						{
							Tags: []*elbv2.Tag{
								{
									Key:   aws.String("tag:elbv2.k8s.aws/cluster"),
									Value: aws.String(DefaultClusterName),
								},
							},
						},
					},
				}, nil
			},
			DeleteTargetGroupFn: func(*elbv2.DeleteTargetGroupInput) (*elbv2.DeleteTargetGroupOutput, error) {
				return &elbv2.DeleteTargetGroupOutput{}, nil
			},
			DescribeLoadBalancersFn: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
				return &elbv2.DescribeLoadBalancersOutput{
					LoadBalancers: []*elbv2.LoadBalancer{
						{
							LoadBalancerArn: aws.String(DefaultLoadBalancerARN),
							VpcId:           aws.String(DefaultVpcID),
						},
					},
				}, nil
			},
			DeleteLoadBalancerFn: func(*elbv2.DeleteLoadBalancerInput) (*elbv2.DeleteLoadBalancerOutput, error) {
				return &elbv2.DeleteLoadBalancerOutput{}, nil
			},
		},
		CWAPI: &CWAPI{
			GetMetricsStatisticsFn: func(input *cloudwatch.GetMetricStatisticsInput) (*cloudwatch.GetMetricStatisticsOutput, error) {
				return &cloudwatch.GetMetricStatisticsOutput{}, nil
			},
		},
		CFAPI: &CFAPI{DescribeStacksFn: []func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error){}},
	}
}

var errBad = fmt.Errorf("something bad")

// NewBadCloudProvider returns a mocked cloud provider with failure set on all
// nolint: funlen
func NewBadCloudProvider() *CloudProvider {
	return &CloudProvider{
		EC2API: &EC2API{
			DescribeSubnetsFn: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
				return nil, errBad
			},
			DescribeAddressesFn:        nil,
			DescribeInternetGatewaysFn: nil,
			DescribeVpcsFn:             nil,
			DescribeSecurityGroupsFn: func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
				return nil, errBad
			},
			AuthorizeSecurityGroupIngressFn: nil,
			RevokeSecurityGroupIngressFn:    nil,
			DeleteSecurityGroupFn: func(*ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
				return nil, errBad
			},
		},
		EKSAPI: &EKSAPI{
			DescribeFargateProfileFn: func(*eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error) {
				return nil, errBad
			},
		},
		IAMAPI: &IAMAPI{
			AttachRolePolicyFn: func(*iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
				return nil, errBad
			},
			DetachRolePolicyFn: func(*iam.DetachRolePolicyInput) (*iam.DetachRolePolicyOutput, error) {
				return nil, errBad
			},
		},
		S3API: &S3API{
			PutObjectFn: func(*s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return nil, errBad
			},
			DeleteObjectFn: func(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
				return nil, errBad
			},
			GetObjectFn: func(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return nil, errBad
			},
		},
		SMAPI: &SMAPI{
			RotateSecretFn: func(*secretsmanager.RotateSecretInput) (*secretsmanager.RotateSecretOutput, error) {
				return nil, errBad
			},
			CancelRotateSecretFn: func(*secretsmanager.CancelRotateSecretInput) (*secretsmanager.CancelRotateSecretOutput, error) {
				return nil, errBad
			},
		},
		ACMAPI: &ACMAPI{
			DescribeCertificateFn: func(*acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error) {
				return nil, errBad
			},
			ListCertificatesFn: func(*acm.ListCertificatesInput) (*acm.ListCertificatesOutput, error) {
				return nil, errBad
			},
		},
		ELBv2API: &ELBv2API{
			DescribeListenersFn: func(*elbv2.DescribeListenersInput) (*elbv2.DescribeListenersOutput, error) {
				return nil, errBad
			},
			DeleteListenerFn: func(*elbv2.DeleteListenerInput) (*elbv2.DeleteListenerOutput, error) {
				return nil, errBad
			},
			DescribeTargetGroupsFn: func(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
				return nil, errBad
			},
			DescribeTagsFn: func(*elbv2.DescribeTagsInput) (*elbv2.DescribeTagsOutput, error) {
				return nil, errBad
			},
			DeleteTargetGroupFn: func(*elbv2.DeleteTargetGroupInput) (*elbv2.DeleteTargetGroupOutput, error) {
				return nil, errBad
			},
			DescribeLoadBalancersFn: func(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
				return nil, errBad
			},
			DeleteLoadBalancerFn: func(*elbv2.DeleteLoadBalancerInput) (*elbv2.DeleteLoadBalancerOutput, error) {
				return nil, errBad
			},
		},
		CFAPI: &CFAPI{DescribeStacksFn: []func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error){}},
	}
}

// Subnets returns a valid subnet describe response
func Subnets() []*ec2.Subnet {
	return []*ec2.Subnet{
		{
			AvailabilityZone: aws.String("eu-west-1a"),
			CidrBlock:        aws.String("192.168.0.0/24"),
			SubnetId:         aws.String("subnet-0bb1c79de3EXAMPLE"),
		},
	}
}

// Stacks returns a valid stack describe response
func Stacks(status string) []*cloudformation.Stack {
	return []*cloudformation.Stack{
		{
			StackId:           aws.String(DefaultStackName),
			StackName:         aws.String(DefaultStackName),
			StackStatus:       aws.String(status),
			StackStatusReason: aws.String("something"),
		},
	}
}
