package dependencytree_test

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// nolint funlen
func TestResourceNode_ApplyFunction(t *testing.T) {
	testCases := []struct {
		name                 string
		expectedCalls        int
		generateReceiverTree func() *dependencytree.Node
		generateTargetTree   func() *dependencytree.Node
	}{
		{
			name:          "Applies on all nodes in a tree of 3",
			expectedCalls: 3,
			generateReceiverTree: func() *dependencytree.Node {
				root := &dependencytree.Node{
					Type:     dependencytree.NodeTypeGroup,
					Children: []*dependencytree.Node{},
				}

				root.Children = append(root.Children, &dependencytree.Node{
					Type: dependencytree.NodeTypeVPC,
					Children: []*dependencytree.Node{
						{
							Type:     dependencytree.NodeTypeZone,
							Children: []*dependencytree.Node{},
						},
					},
				})

				return root
			},
			generateTargetTree: func() *dependencytree.Node {
				root := &dependencytree.Node{
					Type:     dependencytree.NodeTypeGroup,
					Children: []*dependencytree.Node{},
				}

				root.Children = append(root.Children, &dependencytree.Node{
					Type: dependencytree.NodeTypeVPC,
					Children: []*dependencytree.Node{
						{
							Type:     dependencytree.NodeTypeZone,
							Children: []*dependencytree.Node{},
						},
					},
				})

				return root
			},
		},
		{
			name:          "Applies to all nodes in a tree of 2",
			expectedCalls: 2,
			generateReceiverTree: func() *dependencytree.Node {
				root := &dependencytree.Node{
					Type:     dependencytree.NodeTypeGroup,
					Children: []*dependencytree.Node{},
				}

				root.Children = append(root.Children, &dependencytree.Node{
					Type:     dependencytree.NodeTypeVPC,
					Children: []*dependencytree.Node{},
				})

				return root
			},
			generateTargetTree: func() *dependencytree.Node {
				root := &dependencytree.Node{
					Type:     dependencytree.NodeTypeGroup,
					Children: []*dependencytree.Node{},
				}

				root.Children = append(root.Children, &dependencytree.Node{
					Type:     dependencytree.NodeTypeVPC,
					Children: []*dependencytree.Node{},
				})

				return root
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			desiredTree := tc.generateReceiverTree()
			currentStateTree := tc.generateTargetTree()

			numberOfCalls := 0
			desiredTree.ApplyFunctionWithTarget(func(receiver *dependencytree.Node, target *dependencytree.Node) {
				numberOfCalls += 1
			}, currentStateTree)

			assert.Equal(t, tc.expectedCalls, numberOfCalls)
		})
	}
}

func TestResourceNode_Equals(t *testing.T) {
	testCases := []struct {
		name  string
		nodeA *dependencytree.Node
		nodeB *dependencytree.Node

		expectEquality bool
	}{
		{
			name:           "Should be equal",
			nodeA:          &dependencytree.Node{Type: dependencytree.NodeTypeVPC},
			nodeB:          &dependencytree.Node{Type: dependencytree.NodeTypeVPC},
			expectEquality: true,
		},
		{
			name:           "Should not be equal",
			nodeA:          &dependencytree.Node{Type: dependencytree.NodeTypeVPC},
			nodeB:          &dependencytree.Node{Type: dependencytree.NodeTypeZone},
			expectEquality: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectEquality, tc.nodeA.Equals(tc.nodeB))
		})
	}
}
