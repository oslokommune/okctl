package cfn_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/api"

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
			composer: components.NewVPCComposer("test", "192.168.0.0/20", "eu-west-1"),
		},
		{
			name:     "Minimal VPC composer",
			golden:   "vpc-minimal-cf.yaml",
			composer: components.NewMinimalVPCComposer("test", "192.168.0.0/20", "eu-west-1"),
		},
		{
			name:     "ExternalSecretsPolicy composer",
			golden:   "esp-cloudformation.yaml",
			composer: components.NewExternalSecretsPolicyComposer("test"),
		},
		{
			name:     "AlbIngressControllerPolicy composer",
			golden:   "alb-ingress-cloudformation.yaml",
			composer: components.NewAlbIngressControllerPolicyComposer("tset"),
		},
		{
			name:     "AWSLoadBalancerControllerPolicy composer",
			golden:   "aws-load-balancer-controller.yaml",
			composer: components.NewAWSLoadBalancerControllerComposer("tset"),
		},
		{
			name:     "ExternalDNSPolicy composer",
			golden:   "external-dns-cloudformation.yaml",
			composer: components.NewExternalDNSPolicyComposer("tset"),
		},
		{
			name:     "AutoscalerPolicy composer",
			golden:   "autoscaler-cloudformation.yaml",
			composer: components.NewAutoscalerPolicyComposer("tset"),
		},
		{
			name:     "BlockstoragePolicy composer",
			golden:   "blockstorage-cloudformation.yaml",
			composer: components.NewBlockstoragePolicyComposer("tset"),
		},
		{
			name:     "PublicCertificate composer",
			golden:   "public-certificate-cf.yaml",
			composer: components.NewPublicCertificateComposer("test.oslo.systems.", "AZ12345"),
		},
		{
			name:     "CloudwatchDatasourcePolicyComposer",
			golden:   "cloudwatch-datasource-cf.yaml",
			composer: components.NewCloudwatchDatasourcePolicyComposer("test"),
		},
		{
			name:     "FargateCloudWatchPolicy",
			golden:   "fargate-cloudwatch.yaml",
			composer: components.NewFargateCloudwatchPolicyComposer("test"),
		},
		{
			name:   "UserPool composer",
			golden: "userpool.yaml",
			composer: components.NewUserPool(
				"test",
				"auth.oslo.systems",
				"HFJE38983FAKE",
				"arn://certificate/ihfieh9e9FAKE",
			),
		},
		{
			name:     "UserPoolClient",
			golden:   "userpool-client.yaml",
			composer: components.NewUserPoolClient("argocd", "test", "https://argocd/callback", "GHFE723FAKE"),
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
					ClusterName:       "cluster",
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
			composer: components.NewS3BucketComposer("myBucket", "S3Bucket", false),
		},
		{
			name:     "DynamoDB composer",
			golden:   "dynamodb.yaml",
			composer: components.NewDynamoDBTableComposer("myTable", "myID"),
		},
		{
			name:     "Loki S3BucketPolicy composer",
			golden:   "s3BucketPolicy.yaml",
			composer: components.NewLokiS3PolicyComposer("mock-cluster", "arn:but:is:not:an:arn"),
		},
		{
			name:     "Loki DynamoDB policy composer",
			golden:   "lokiDynamoDBPolicy.yaml",
			composer: components.NewLokiDynamoDBPolicyComposer(api.ID{ClusterName: "mock-cluster"}, "test-prefix_"),
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
