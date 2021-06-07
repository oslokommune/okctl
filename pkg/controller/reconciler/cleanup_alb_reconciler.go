package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cleaner"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type cleanupALBReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	provider v1alpha1.CloudProvider
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *cleanupALBReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeCleanupALB
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *cleanupALBReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *cleanupALBReconciler) Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		// Nothing to do for present
		return result, nil
	case resourcetree.ResourceNodeStateAbsent:
		vpc, err := state.Vpc.GetVpc(cfn.NewStackNamer().Vpc(z.commonMetadata.ClusterID.ClusterName))
		if err != nil {
			return result, fmt.Errorf("getting vpc: %w", err)
		}

		clean := cleaner.New(z.provider)

		err = clean.DeleteDanglingALBs(vpc.VpcID)
		if err != nil {
			return result, fmt.Errorf("cleaning up ALBs: %w", err)
		}

		err = clean.DeleteDanglingTargetGroups(z.commonMetadata.ClusterID.ClusterName)
		if err != nil {
			return result, fmt.Errorf("cleaning target groups: %w", err)
		}
	}

	return result, nil
}

// NewCleanupALBReconciler creates a new reconciler for cleaning up ALBs
func NewCleanupALBReconciler(provider v1alpha1.CloudProvider) Reconciler {
	return &cleanupALBReconciler{
		provider: provider,
	}
}
