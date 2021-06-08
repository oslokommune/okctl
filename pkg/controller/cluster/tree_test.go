package cluster

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"

	"github.com/bmizerany/assert"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// nolint: funlen
func TestTreeCreators(t *testing.T) {
	testCases := []struct {
		name string

		declaration      func() *v1alpha1.Cluster
		existingServices ExistingResources
	}{
		{
			name: "Should produce equal trees when all is enabled",
			declaration: func() *v1alpha1.Cluster {
				declaration := v1alpha1.NewCluster()
				declaration.Databases = &v1alpha1.ClusterDatabases{
					Postgres: []v1alpha1.ClusterDatabasesPostgres{
						{
							Name:      "test",
							User:      "test",
							Namespace: "test",
						},
					},
				}

				return &declaration
			},
			existingServices: ExistingResources{
				hasServiceQuotaCheck:                  true,
				hasAWSLoadBalancerController:          true,
				hasCluster:                            true,
				hasExternalDNS:                        true,
				hasExternalSecrets:                    true,
				hasAutoscaler:                         true,
				hasBlockstorage:                       true,
				hasKubePromStack:                      true,
				hasLoki:                               true,
				hasPromtail:                           true,
				hasTempo:                              true,
				hasIdentityManager:                    true,
				hasArgoCD:                             true,
				hasPrimaryHostedZone:                  true,
				hasVPC:                                true,
				hasDelegatedHostedZoneNameservers:     true,
				hasDelegatedHostedZoneNameserversTest: true,
				hasUsers:                              false,
				hasPostgres: map[string]*v1alpha1.ClusterDatabasesPostgres{
					"test": {
						Name:      "test",
						User:      "test",
						Namespace: "test",
					},
				},
			},
		},
		{
			name: "Should produce equal trees when all but ExternalDNS is enabled",
			declaration: func() *v1alpha1.Cluster {
				declaration := v1alpha1.NewCluster()
				declaration.Integrations.ExternalDNS = false

				return &declaration
			},
			existingServices: ExistingResources{
				hasServiceQuotaCheck:                  true,
				hasAWSLoadBalancerController:          true,
				hasCluster:                            true,
				hasExternalDNS:                        false,
				hasExternalSecrets:                    true,
				hasAutoscaler:                         true,
				hasBlockstorage:                       true,
				hasKubePromStack:                      true,
				hasLoki:                               true,
				hasPromtail:                           true,
				hasTempo:                              true,
				hasIdentityManager:                    true,
				hasArgoCD:                             true,
				hasPrimaryHostedZone:                  true,
				hasVPC:                                true,
				hasDelegatedHostedZoneNameservers:     true,
				hasDelegatedHostedZoneNameserversTest: true,
				hasUsers:                              false,
				hasPostgres:                           nil,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			desiredTree := CreateResourceDependencyTree()
			desiredTree.ApplyFunction(applyDeclaration(tc.declaration()))

			currentStateTree := CreateResourceDependencyTree()
			currentStateTree.ApplyFunction(applyExistingState(tc.existingServices))

			assert.Equal(t, desiredTree.String(), currentStateTree.String())
		})
	}
}

// Disclaimer: hasDependency only checks one level
func hasDependency(tree *dependencytree.Node, target dependencytree.NodeType, potentialDependency dependencytree.NodeType) bool {
	dependencyNode := tree.GetNode(&dependencytree.Node{Type: potentialDependency})

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

		with dependencytree.NodeType

		expectDependencies []dependencytree.NodeType
		expectFail         bool
	}{
		{
			name: "Sanity check: hosted zone should not be dependent on cluster",

			with: dependencytree.NodeTypeZone,

			expectDependencies: []dependencytree.NodeType{dependencytree.NodeTypeCluster},
			expectFail:         true,
		},
		{
			name: "KubePromStack has the correct relations",

			with: dependencytree.NodeTypeKubePromStack,

			expectDependencies: []dependencytree.NodeType{dependencytree.NodeTypeIdentityManager},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tree := CreateResourceDependencyTree()

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
