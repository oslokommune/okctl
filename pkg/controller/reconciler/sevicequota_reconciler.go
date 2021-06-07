package reconciler

import (
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/oslokommune/okctl/pkg/servicequota"
	"github.com/pkg/errors"
)

// ServiceQuotaReconciler handles reconciliation for service quotas
type ServiceQuotaReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	provider v1alpha1.CloudProvider
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (r *ServiceQuotaReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeServiceQuota
}

// SetCommonMetadata knows how to store common metadata on the reconciler. This should do nothing if common metadata is
// not needed
func (r *ServiceQuotaReconciler) SetCommonMetadata(meta *resourcetree.CommonMetadata) {
	r.commonMetadata = meta
}

// Reconcile knows how to create, update and delete the relevant resource
func (r *ServiceQuotaReconciler) Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		vpc, err := state.Vpc.GetVpc(cfn.NewStackNamer().Vpc(r.commonMetadata.Declaration.Metadata.Name))
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			return result, fmt.Errorf("getting vpc: %w", err)
		}

		vpcProvisioned := vpc != nil

		// This should be a reconciler, e.g., a root node
		err = servicequota.CheckQuotas(
			servicequota.NewVpcCheck(vpcProvisioned, constant.DefaultRequiredVpcs, r.provider),
			servicequota.NewEipCheck(vpcProvisioned, constant.DefaultRequiredEpis, r.provider),
			servicequota.NewIgwCheck(vpcProvisioned, constant.DefaultRequiredIgws, r.provider),
			servicequota.NewFargateCheck(constant.DefaultRequiredFargateOnDemandPods, r.provider),
		)
		if err != nil {
			return result, fmt.Errorf("checking service quotas: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return result, nil
	}

	return result, nil
}

// NewServiceQuotaReconciler returns an initialised service quota reconciler
func NewServiceQuotaReconciler(provider v1alpha1.CloudProvider) *ServiceQuotaReconciler {
	return &ServiceQuotaReconciler{
		provider: provider,
	}
}
