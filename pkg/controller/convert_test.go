package controller

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// nolint: funlen
func TestTreeCreators(t *testing.T) {
	testCases := []struct {
		name string

		declaration      func() *v1alpha1.Cluster
		existingServices ExistingServices
	}{
		{
			name: "Should produce equal trees when all is enabled",
			declaration: func() *v1alpha1.Cluster {
				declaration := v1alpha1.NewDefaultCluster("", "", "", "", "", "")

				return &declaration
			},
			existingServices: ExistingServices{
				hasArgoCD:                         true,
				hasALBIngressController:           false,
				hasAWSLoadBalancerController:      true,
				hasCluster:                        true,
				hasExternalDNS:                    true,
				hasKubePromStack:                  true,
				hasExternalSecrets:                true,
				hasBlockstorage:                   true,
				hasAutoscaler:                     true,
				hasGithubSetup:                    true,
				hasIdentityManager:                true,
				hasPrimaryHostedZone:              true,
				hasVPC:                            true,
				hasDelegatedHostedZoneNameservers: true,
			},
		},
		{
			name: "Should produce equal trees when all but ExternalDNS is enabled",
			declaration: func() *v1alpha1.Cluster {
				declaration := v1alpha1.NewDefaultCluster("", "", "", "", "", "")
				declaration.Integrations.ExternalDNS = false

				return &declaration
			},
			existingServices: ExistingServices{
				hasArgoCD:                         true,
				hasALBIngressController:           false,
				hasAWSLoadBalancerController:      true,
				hasCluster:                        true,
				hasExternalDNS:                    false,
				hasExternalSecrets:                true,
				hasAutoscaler:                     true,
				hasKubePromStack:                  true,
				hasBlockstorage:                   true,
				hasGithubSetup:                    true,
				hasIdentityManager:                true,
				hasPrimaryHostedZone:              true,
				hasVPC:                            true,
				hasDelegatedHostedZoneNameservers: true,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			desiredTree := CreateDesiredStateTree(tc.declaration())
			currentStateTree := CreateCurrentStateTree(&tc.existingServices)

			assert.Equal(t, desiredTree.String(), currentStateTree.String())
		})
	}
}
