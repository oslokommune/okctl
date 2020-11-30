package resourcetree_test

import (
	"github.com/bmizerany/assert"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"testing"
)

func TestResourceNode_ApplyFunction(t *testing.T) {
	testCases := []struct {
	    name string
	    expectedCalls int
	    generateReceiverTree func() *resourcetree.ResourceNode
	    generateTargetTree func() *resourcetree.ResourceNode
	}{
	    {
	        name: "Applies on all nodes in a tree of 3",
	        expectedCalls: 3,
	        generateReceiverTree: func() *resourcetree.ResourceNode {
				root := &resourcetree.ResourceNode{ 
					Type: resourcetree.ResourceNodeTypeGroup,
					Children: []*resourcetree.ResourceNode{},
				}
				
				root.Children = append(root.Children, &resourcetree.ResourceNode{
					Type: resourcetree.ResourceNodeTypeVPC,
					Children: []*resourcetree.ResourceNode{
						{
							Type: resourcetree.ResourceNodeTypeZone,
							Children: []*resourcetree.ResourceNode{},
						},
					},
				})
				
				return root
			},
			generateTargetTree: func() *resourcetree.ResourceNode {
				root := &resourcetree.ResourceNode{
					Type: resourcetree.ResourceNodeTypeGroup,
					Children: []*resourcetree.ResourceNode{},
				}

				root.Children = append(root.Children, &resourcetree.ResourceNode{
					Type: resourcetree.ResourceNodeTypeVPC,
					Children: []*resourcetree.ResourceNode{
						{
							Type: resourcetree.ResourceNodeTypeZone,
							Children: []*resourcetree.ResourceNode{},
						},
					},
				})

				return root
			},
		},
		{
			name: "Applies to all nodes in a tree of 2",
			expectedCalls: 2,
			generateReceiverTree: func() *resourcetree.ResourceNode {
				root := &resourcetree.ResourceNode{
					Type: resourcetree.ResourceNodeTypeGroup,
					Children: []*resourcetree.ResourceNode{},
				}

				root.Children = append(root.Children, &resourcetree.ResourceNode{
					Type: resourcetree.ResourceNodeTypeVPC,
					Children: []*resourcetree.ResourceNode{},
				})

				return root
			},
			generateTargetTree: func() *resourcetree.ResourceNode {
				root := &resourcetree.ResourceNode{
					Type: resourcetree.ResourceNodeTypeGroup,
					Children: []*resourcetree.ResourceNode{},
				}

				root.Children = append(root.Children, &resourcetree.ResourceNode{
					Type: resourcetree.ResourceNodeTypeVPC,
					Children: []*resourcetree.ResourceNode{},
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
	    	desiredTree.ApplyFunction(func(receiver *resourcetree.ResourceNode, target *resourcetree.ResourceNode) {
	    		numberOfCalls += 1
			}, currentStateTree)

	    	assert.Equal(t, tc.expectedCalls, numberOfCalls)
	    })
	}
}

func TestResourceNode_Equals(t *testing.T) {
	testCases := []struct {
	    name string
	    nodeA *resourcetree.ResourceNode
	    nodeB *resourcetree.ResourceNode
	    
	    expectEquality bool
	}{
	    {
	        name: "Should be equal",
	        nodeA: &resourcetree.ResourceNode{ Type: resourcetree.ResourceNodeTypeVPC },
	        nodeB: &resourcetree.ResourceNode{ Type: resourcetree.ResourceNodeTypeVPC },
	        expectEquality: true,
	    },
	    {
	    	name: "Should not be equal",
	    	nodeA: &resourcetree.ResourceNode{ Type: resourcetree.ResourceNodeTypeVPC },
	    	nodeB: &resourcetree.ResourceNode{ Type: resourcetree.ResourceNodeTypeZone },
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
