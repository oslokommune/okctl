package reconciliation

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/resourcetree"
)

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *awsLoadBalancerControllerReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeAWSLoadBalancerController
}

// albIngressReconciler contains service and metadata for the relevant resource
type awsLoadBalancerControllerReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	client         client.AWSLoadBalancerControllerService
}

// SetCommonMetadata stores common metadata for later use
func (z *awsLoadBalancerControllerReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *awsLoadBalancerControllerReconciler) Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		vpc, err := state.Vpc.GetVpc(
			cfn.NewStackNamer().Vpc(z.commonMetadata.Declaration.Metadata.Name),
		)
		if err != nil {
			return result, fmt.Errorf("getting vpc: %w", err)
		}

		_, err = z.client.CreateAWSLoadBalancerController(z.commonMetadata.Ctx, client.CreateAWSLoadBalancerControllerOpts{
			ID:    z.commonMetadata.ClusterID,
			VPCID: vpc.VpcID,
		})
		if err != nil {
			return result, fmt.Errorf("creating aws load balancer controller: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteAWSLoadBalancerController(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting aws load balancer controller: %w", err)
		}
	}

	return result, nil
}

// NewAWSLoadBalancerControllerReconciler creates a new reconciler for the aws load balancer controller resource
func NewAWSLoadBalancerControllerReconciler(client client.AWSLoadBalancerControllerService) reconciliation.Reconciler {
	return &awsLoadBalancerControllerReconciler{
		client: client,
	}
}
