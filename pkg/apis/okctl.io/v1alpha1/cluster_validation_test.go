package v1alpha1_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestValidateCluster(t *testing.T) {
	testCases := []struct {
		name        string
		withCluster func() v1alpha1.Cluster
		expectError string
	}{
		{
			name:        "Should pass when everything is A-ok",
			withCluster: newPassingCluster,
			expectError: "",
		},
		{
			name: "Should fail when name is empty",
			withCluster: func() v1alpha1.Cluster {
				c := newPassingCluster()

				c.Metadata.Name = ""

				return c
			},
			expectError: "metadata: (name: cannot be blank.).",
		},
		{
			name: "Should fail if clusterRootDomain is missing",
			withCluster: func() v1alpha1.Cluster {
				c := newPassingCluster()

				c.ClusterRootDomain = ""

				return c
			},
			expectError: "clusterRootDomain: cannot be blank.",
		},
		{
			name: "Should fail if clusterRootDomain have improper casing",
			withCluster: func() v1alpha1.Cluster {
				c := newPassingCluster()

				c.ClusterRootDomain = "ThisIsNotAllowed.oslo.systems"

				return c
			},
			expectError: "clusterRootDomain: must be in lower case.",
		},
		{
			name: "Cluster with argocd enabled, cognito disabled",
			withCluster: func() v1alpha1.Cluster {
				c := newPassingCluster()

				c.Integrations.Cognito = false

				return c
			},
			expectError: "integrations: (cognito: is required when argocd is enabled.).",
		},
		{
			// This happens when the user doesn't provide a clusterRootDomain
			name: "Should fail if clusterRootDomain starts with dash",
			withCluster: func() v1alpha1.Cluster {
				c := newPassingCluster()

				c.ClusterRootDomain = "-.oslo.systems"

				return c
			},
			expectError: "clusterRootDomain: invalid domain name: '-.oslo.systems'.",
		},
		{
			// This happens when the user doesn't provide a clusterRootDomain
			name: "Should fail if clusterRootDomain contains space",
			withCluster: func() v1alpha1.Cluster {
				c := newPassingCluster()

				c.ClusterRootDomain = "hel lo.oslo.systems"

				return c
			},
			expectError: "clusterRootDomain: invalid domain name: 'hel lo.oslo.systems'.",
		},
		{
			// This happens when the user doesn't provide a clusterRootDomain
			name: "Should fail if clusterRootDomain doesn't end with auto.oslo.systems and automerge enabled",
			withCluster: func() v1alpha1.Cluster {
				c := newPassingCluster()

				c.ClusterRootDomain = "hello.oslo.systems"
				c.Experimental = &v1alpha1.ClusterExperimental{
					AutomatizeZoneDelegation: true,
				}

				return c
			},
			expectError: "clusterRootDomain: with automatizeZoneDelegation enabled, must end with auto.oslo.systems.",
		},
		{
			// This happens when the user doesn't provide a clusterRootDomain
			name: "Should succeed if clusterRootDomain ends with auto.oslo.systems and automerge enabled",
			withCluster: func() v1alpha1.Cluster {
				c := newPassingCluster()

				c.ClusterRootDomain = "hello.auto.oslo.systems"
				c.Experimental = &v1alpha1.ClusterExperimental{
					AutomatizeZoneDelegation: true,
				}

				return c
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.withCluster().Validate()

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newPassingCluster() v1alpha1.Cluster {
	return newCluster(
		"myCluster",
		"mycluster-dev.oslo.systems",
		"my-organization",
		"my-github-repo",
		"123456789012",
	)
}
