package reconciliation

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/cleaner"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

type cleanupSGReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	provider v1alpha1.CloudProvider
}

// NodeType returns the relevant NodeType for this reconciler
func (z *cleanupSGReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeCleanupSG
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *cleanupSGReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *cleanupSGReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		// Nothing to do for present
		return result, nil
	case dependencytree.NodeStateAbsent:
		vpc, err := state.Vpc.GetVpc(cfn.NewStackNamer().Vpc(z.commonMetadata.Declaration.Metadata.Name))
		if err != nil {
			return result, fmt.Errorf("getting vpc: %w", err)
		}

		err = cleaner.New(z.provider).DeleteDanglingSecurityGroups(vpc.VpcID)
		if err != nil {
			return result, fmt.Errorf("cleaning up SGs: %w", err)
		}
	}

	return result, nil
}

// NewCleanupSGReconciler creates a new reconciler for cleaning up SGs
func NewCleanupSGReconciler(provider v1alpha1.CloudProvider) reconciliation.Reconciler {
	return &cleanupSGReconciler{
		provider: provider,
	}
}
