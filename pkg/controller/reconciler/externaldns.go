package reconciler

import (
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ExternalDNSResourceState contains runtime data needed in Reconcile()
type ExternalDNSResourceState struct {
	HostedZoneID string
	Domain       string
}

// ExternalDNSReconciler contains service and metadata for the relevant resource
type ExternalDNSReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.ExternalDNSService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *ExternalDNSReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *ExternalDNSReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	resourceState, ok := node.ResourceState.(ExternalDNSResourceState)
	if !ok {
		return nil, errors.New("error casting External DNS resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateExternalDNS(z.commonMetadata.Ctx, client.CreateExternalDNSOpts{
			ID:           z.commonMetadata.ClusterID,
			HostedZoneID: resourceState.HostedZoneID,
			Domain:       resourceState.Domain,
		})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating external DNS: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteExternalDNS(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error deleting external DNS: %w", err)
		}
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewExternalDNSReconciler creates a new reconciler for the ExternalDNS resource
func NewExternalDNSReconciler(client client.ExternalDNSService) *ExternalDNSReconciler {
	return &ExternalDNSReconciler{
		client: client,
	}
}
