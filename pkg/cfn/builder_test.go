package cfn_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestBuilderAndComposers(t *testing.T) {
	testCases := []struct {
		name     string
		golden   string
		composer cfn.Composer
	}{
		{
			name:     "Builder with VPC composer",
			golden:   "vpc-cloudformation.yaml",
			composer: components.NewVPCComposer("test", "test", "192.168.0.0/20", "eu-west-1"),
		},
		{
			name:     "Builder with Minimal VPC composer",
			golden:   "vpc-minimal-cf.yaml",
			composer: components.NewMinimalVPCComposer("test", "test", "192.168.0.0/20", "eu-west-1"),
		},
		{
			name:     "Builder with ExternalSecretsPolicy composer",
			golden:   "esp-cloudformation.yaml",
			composer: components.NewExternalSecretsPolicyComposer("repo", "test"),
		},
		{
			name:     "Builder with AlbIngressControllerPolicy composer",
			golden:   "alb-ingress-cloudformation.yaml",
			composer: components.NewAlbIngressControllerPolicyComposer("repo", "env"),
		},
		{
			name:     "Builder with ExternalDNSPolicy composer",
			golden:   "external-dns-cloudformation.yaml",
			composer: components.NewExternalDNSPolicyComposer("repo", "env"),
		},
		{
			name:     "Builder with PublicCertificate composer",
			golden:   "public-certificate-cf.yaml",
			composer: components.NewPublicCertificateComposer("test.oslo.systems.", "AZ12345"),
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
