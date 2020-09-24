package argocd_test

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/oslokommune/okctl/pkg/helm/charts/argocd"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultValues(t *testing.T) {
	testCases := []struct {
		name   string
		values *argocd.Values
		golden string
	}{
		{
			name: "Default values should generate valid yaml",
			values: argocd.NewDefaultValues(argocd.ValuesOpts{
				URL:                  "https://argocd.test.oslo.systems",
				HostName:             "argocd.test.oslo.systems",
				CertificateARN:       "arn:aws:acm:eu-west-1:123456789012/certificate/12345abc",
				ClientID:             "client12345",
				Organisation:         "oslokommune",
				AuthDomain:           "auth.oslo.systems",
				UserPoolID:           "VBNJ6723FAKE",
				RepoURL:              "git@github.com:oslokommune/test.git",
				RepoName:             "test",
				PrivateKeySecretName: "argocd-test-oslokommune-private-key",
				PrivateKeySecretKey:  "ssh-private-key",
			}),
			golden: "argocd-values.yaml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := yaml.Marshal(tc.values)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
