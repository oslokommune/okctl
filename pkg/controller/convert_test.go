package controller

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/controller/resourcetree"

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

// Disclaimer: hasDependency only checks one level
func hasDependency(tree *resourcetree.ResourceNode, target resourcetree.ResourceNodeType, potentialDependency resourcetree.ResourceNodeType) bool {
	dependencyNode := tree.GetNode(&resourcetree.ResourceNode{Type: potentialDependency})

	for _, child := range dependencyNode.Children {
		if child.Type == target {
			return true
		}
	}

	return false
}

func TestEnsureRelations(t *testing.T) {
	testCases := []struct {
		name string

		with resourcetree.ResourceNodeType

		expectDependencies []resourcetree.ResourceNodeType
		expectFail         bool
	}{
		{
			name: "Sanity check: hosted zone should not be dependent on cluster",

			with: resourcetree.ResourceNodeTypeZone,

			expectDependencies: []resourcetree.ResourceNodeType{resourcetree.ResourceNodeTypeCluster},
			expectFail:         true,
		},
		{
			name: "KubePromStack has the correct relations",

			with: resourcetree.ResourceNodeTypeKubePromStack,

			expectDependencies: []resourcetree.ResourceNodeType{resourcetree.ResourceNodeTypeIdentityManager},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tree := CreateCurrentStateTree(&ExistingServices{})

			for _, dependency := range tc.expectDependencies {
				if tc.expectFail {
					// Ensure no dependency
					assert.Equal(t, false, hasDependency(tree, tc.with, dependency))
				} else {
					// Ensure dependency
					assert.Equal(t, true, hasDependency(tree, tc.with, dependency))
					// Ensure no single level cyclic dependency
					assert.NotEqual(t, true, hasDependency(tree, dependency, tc.with))
				}
			}
		})
	}
}
