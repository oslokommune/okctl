package controller

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

func TestTreeCreators(t *testing.T) {
	testCases := []struct {
		name string

		declaration      func() *v1alpha1.Cluster
		existingServices existingServices
	}{
		{
			name: "Should produce equal trees when all is enabled",
			declaration: func() *v1alpha1.Cluster {
				declaration := v1alpha1.NewDefaultCluster("", "", "", "", "", "")

				return &declaration
			},
			existingServices: existingServices{
				hasALBIngressController: true,
				hasCluster:              true,
				hasExternalDNS:          true,
				hasExternalSecrets:      true,
				hasGithubSetup:          true,
				hasIdentityManager:      true,
				hasPrimaryHostedZone:    true,
				hasVPC:                  true,
			},
		},
		{
			name: "Should produce equal trees when all but ExternalDNS is enabled",
			declaration: func() *v1alpha1.Cluster {
				declaration := v1alpha1.NewDefaultCluster("", "", "", "", "", "")
				declaration.Integrations.ExternalDNS = false

				return &declaration
			},
			existingServices: existingServices{
				hasALBIngressController: true,
				hasCluster:              true,
				hasExternalDNS:          false,
				hasExternalSecrets:      true,
				hasGithubSetup:          true,
				hasIdentityManager:      true,
				hasPrimaryHostedZone:    true,
				hasVPC:                  true,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			desiredTree := CreateDesiredStateGraph(tc.declaration())
			currentStateTree := CreateCurrentStateGraph(&tc.existingServices)

			assert.Equal(t, desiredTree.String(), currentStateTree.String())
		})
	}
}
