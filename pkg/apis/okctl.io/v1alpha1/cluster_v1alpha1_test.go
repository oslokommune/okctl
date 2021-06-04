package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"github.com/sebdah/goldie/v2"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

func TestMarshalCluster(t *testing.T) {
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
			cluster: newCluster(
				"okctl-stage",
				"okctl-stage.oslo.systems",
				"oslokommune",
				"okctl-iac",
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

func TestValidations(t *testing.T) {
	testCases := []struct {
		name      string
		with      v1alpha1.Cluster
		expectErr string
	}{
		{
			name: "Should not trigger KM302",
			with: func() v1alpha1.Cluster {
				c := newCluster(
					"test",
					"okctl.io",
					"oslokommune",
					"okctl",
					"012345678912",
				)

				c.VPC = &v1alpha1.ClusterVPC{
					CIDR:             "192.168.0.0/20",
					HighAvailability: false,
				}

				return c
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.with.Validate()

			if tc.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.expectErr, err)
			}
		})
	}
}

func newCluster(name, rootDomain, githubOrganization, githubRepository, awsAccountID string) v1alpha1.Cluster {
	c := v1alpha1.NewCluster()
	c.Metadata.Name = name
	c.Metadata.AccountID = awsAccountID
	c.Github.Organisation = githubOrganization
	c.Github.Repository = githubRepository
	c.ClusterRootDomain = rootDomain

	return c
}
