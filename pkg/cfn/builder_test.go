package cfn_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestBuilderAndComposers(t *testing.T) {
	testCases := []struct {
		name     string
		golden   string
		composer cfn.Composer
	}{
		{
			name:     "VPC composer",
			golden:   "vpc-cloudformation.yaml",
			composer: components.NewVPCComposer("test", "test", "192.168.0.0/20", "eu-west-1"),
		},
		{
			name:     "Minimal VPC composer",
			golden:   "vpc-minimal-cf.yaml",
			composer: components.NewMinimalVPCComposer("test", "test", "192.168.0.0/20", "eu-west-1"),
		},
		{
			name:     "ExternalSecretsPolicy composer",
			golden:   "esp-cloudformation.yaml",
			composer: components.NewExternalSecretsPolicyComposer("repo", "test"),
		},
		{
			name:     "AlbIngressControllerPolicy composer",
			golden:   "alb-ingress-cloudformation.yaml",
			composer: components.NewAlbIngressControllerPolicyComposer("repo", "env"),
		},
		{
			name:     "AWSLoadBalancerControllerPolicy composer",
			golden:   "aws-load-balancer-controller.yaml",
			composer: components.NewAWSLoadBalancerControllerComposer("repo", "env"),
		},
		{
			name:     "ExternalDNSPolicy composer",
			golden:   "external-dns-cloudformation.yaml",
			composer: components.NewExternalDNSPolicyComposer("repo", "env"),
		},
		{
			name:     "AutoscalerPolicy composer",
			golden:   "autoscaler-cloudformation.yaml",
			composer: components.NewAutoscalerPolicyComposer("repo", "env"),
		},
		{
			name:     "BlockstoragePolicy composer",
			golden:   "blockstorage-cloudformation.yaml",
			composer: components.NewBlockstoragePolicyComposer("repo", "env"),
		},
		{
			name:     "PublicCertificate composer",
			golden:   "public-certificate-cf.yaml",
			composer: components.NewPublicCertificateComposer("test.oslo.systems.", "AZ12345"),
		},
		{
			name:     "CloudwatchDatasourcePolicyComposer",
			golden:   "cloudwatch-datasource-cf.yaml",
			composer: components.NewCloudwatchDatasourcePolicyComposer("repo", "env"),
		},
		{
			name:     "FargateCloudWatchPolicy",
			golden:   "fargate-cloudwatch.yaml",
			composer: components.NewFargateCloudwatchPolicyComposer("repo", "env"),
		},
		{
			name:   "UserPool composer",
			golden: "userpool.yaml",
			composer: components.NewUserPool(
				"env",
				"repo",
				"auth.oslo.systems",
				"HFJE38983FAKE",
				"arn://certificate/ihfieh9e9FAKE",
			),
		},
		{
			name:     "UserPoolClient",
			golden:   "userpool-client.yaml",
			composer: components.NewUserPoolClient("argocd", "test", "test", "https://argocd/callback", "GHFE723FAKE"),
		},
		{
			name:   "AliasRecordSet composer",
			golden: "alias-recordset.yaml",
			composer: components.NewAliasRecordSet(

				"DomainPoolAuth",
				"cloudfront-us-east-1.aws.com",
				"HJOJF678FAKE",
				"auth.oslo.systems",
				"GHFJE78FAKE",
			),
		},
		{
			name:   "RDSPostgres composer",
			golden: "rds-postgres.yaml",
			composer: components.NewRDSPostgresComposer(
				components.RDSPostgresComposerOpts{
					ApplicationDBName: "okctl",
					AWSAccountID:      "123456789012",
					Repository:        "repo",
					Environment:       "env",
					DBSubnetGroupName: "myDBSubnetGroupName",
					UserName:          "admin",
					VpcID:             "vpcid-w9ufe",
					VPCDBSubnetIDs:    []string{"dbsubnetid-123okf", "dbsubnetid-fjeo338"},
					VPCDBSubnetCIDRs:  []string{"192.168.1.0/20", "192.168.2.0/20"},
				},
			),
		},
		{
			name:     "S3Bucket composer",
			golden:   "s3-bucket.yaml",
			composer: components.NewS3BucketComposer("myBucket", "repo", "env"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			b := cfn.New(tc.composer)

			got, err := b.Build()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
