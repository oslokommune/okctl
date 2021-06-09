package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// externalDNSReconciler contains service and metadata for the relevant resource
type externalDNSReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.ExternalDNSService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *externalDNSReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeExternalDNS
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *externalDNSReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *externalDNSReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		_, err = z.client.CreateExternalDNS(z.commonMetadata.Ctx, client.CreateExternalDNSOpts{
			ID:           z.commonMetadata.ClusterID,
			HostedZoneID: hz.HostedZoneID,
			Domain:       z.commonMetadata.Declaration.ClusterRootDomain,
		})
		if err != nil {
			result.Requeue = errors.IsKind(err, errors.Timeout)

			return result, fmt.Errorf("creating external DNS: %w", err)
		}
	case dependencytree.NodeStateAbsent:
		err = z.client.DeleteExternalDNS(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting external DNS: %w", err)
		}
	}

	return result, nil
}

// NewExternalDNSReconciler creates a new reconciler for the ExternalDNS resource
func NewExternalDNSReconciler(client client.ExternalDNSService) reconciliation.Reconciler {
	return &externalDNSReconciler{
		client: client,
	}
}
