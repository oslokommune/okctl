package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"github.com/sebdah/goldie/v2"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

func TestCluster(t *testing.T) {
	testCases := []struct {
		name    string
		cluster v1alpha1.Cluster
		golden  string
	}{
		{
			name:    "Empty cluster",
			cluster: v1alpha1.Cluster{},
			golden:  "empty-cluster.yml",
		},
		{
			name: "Default cluster",
			cluster: v1alpha1.NewDefaultCluster(
				"okctl",
				"stage",
				"oslokommune",
				"okctl-iac",
				"kjøremiljø",
				"123456789012",
			),
			golden: "default-cluster.yml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := yaml.Marshal(tc.cluster)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}

func TestValidateCluster(t *testing.T) {
	testCases := []struct {
		name      string
		cluster   v1alpha1.Cluster
		expectErr bool
		expect    interface{}
	}{
		{
			name: "Default cluster",
			cluster: v1alpha1.NewDefaultCluster(
				"okctl",
				"stage",
				"oslokommune",
				"okctl-iac",
				"kjøremiljø",
				"123456789012",
			),
		},
		{
			name: "Cluster with argocd enabled, cognito disabled",
			cluster: func() v1alpha1.Cluster {
				c := v1alpha1.NewDefaultCluster(
					"okctl",
					"stage",
					"oslokommune",
					"okctl-iac",
					"kjøremiljø",
					"123456789012",
				)

				c.Integrations.Cognito = false

				return c
			}(),
			expect:    "integrations: (cognito: is required when argocd is enabled.).",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.cluster.Validate()

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
