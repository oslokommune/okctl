package reconciler

import (
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ExternalDNSResourceState contains runtime data needed in Reconcile()
type ExternalDNSResourceState struct {
	PrimaryHostedZoneID string
}

// externalDNSReconciler contains service and metadata for the relevant resource
type externalDNSReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.ExternalDNSService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *externalDNSReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeExternalDNS
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *externalDNSReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *externalDNSReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	resourceState, ok := node.ResourceState.(ExternalDNSResourceState)
	if !ok {
		return result, errors.New("casting External DNS resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateExternalDNS(z.commonMetadata.Ctx, client.CreateExternalDNSOpts{
			ID:           z.commonMetadata.ClusterID,
			HostedZoneID: resourceState.PrimaryHostedZoneID,
			Domain:       z.commonMetadata.Declaration.ClusterRootURL,
		})
		if err != nil {
			return result, fmt.Errorf("creating external DNS: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteExternalDNS(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting external DNS: %w", err)
		}
	}

	return result, nil
}

// NewExternalDNSReconciler creates a new reconciler for the ExternalDNS resource
func NewExternalDNSReconciler(client client.ExternalDNSService) Reconciler {
	return &externalDNSReconciler{
		client: client,
	}
}
