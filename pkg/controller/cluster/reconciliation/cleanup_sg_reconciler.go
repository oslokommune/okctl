package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cleaner"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
)

const cleanupSGReconcilerIdentifier = "security group cleaner"

type cleanupSGReconciler struct {
	provider v1alpha1.CloudProvider
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *cleanupSGReconciler) Reconcile(_ context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	if !meta.Purge {
		return reconciliation.Result{Requeue: false}, nil
	}

	clusterExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.CheckClusterExistanceError, err)
	}

	if clusterExists {
		return reconciliation.Result{Requeue: true}, nil
	}

	vpcExists, err := state.Vpc.HasVPC(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.CheckVpcExistenceError, err)
	}

	if !vpcExists {
		return reconciliation.Result{}, nil
	}

	vpc, err := state.Vpc.GetVpc(cfn.NewStackNamer().Vpc(meta.ClusterDeclaration.Metadata.Name))
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.GetVpcError, err)
	}

	err = cleaner.New(z.provider).DeleteDanglingSecurityGroups(vpc.VpcID)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.DeleteSecurityGroupError, err)
	}

	return reconciliation.Result{}, nil
}

// String returns the identifier for this reconciler
func (z *cleanupSGReconciler) String() string {
	return cleanupSGReconcilerIdentifier
}

// NewCleanupSGReconciler creates a new reconciler for cleaning up SGs
func NewCleanupSGReconciler(provider v1alpha1.CloudProvider) reconciliation.Reconciler {
	return &cleanupSGReconciler{
		provider: provider,
	}
}
