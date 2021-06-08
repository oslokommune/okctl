package reconciliation

import (
	"errors"
	"fmt"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// containerRepositoryReconciler contains service and metadata for the relevant resource
type containerRepositoryReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.ContainerRepositoryService
}

// NodeType returns the relevant NodeType for this reconciler
func (c *containerRepositoryReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeContainerRepository
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (c *containerRepositoryReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	c.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (c *containerRepositoryReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		_, err := state.ContainerRepository.GetContainerRepository(c.commonMetadata.ApplicationDeclaration.Image.Name)
		if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
			return reconciliation.Result{}, fmt.Errorf("getting container repository: %w", err)
		}

		if errors.Is(err, stormpkg.ErrNotFound) {
			_, err := c.client.CreateContainerRepository(c.commonMetadata.Ctx, client.CreateContainerRepositoryOpts{
				ClusterID: c.commonMetadata.ClusterID,
				ImageName: c.commonMetadata.ApplicationDeclaration.Image.Name,
			})
			if err != nil {
				return reconciliation.Result{}, fmt.Errorf("creating container repository: %w", err)
			}
		}

		return reconciliation.Result{}, nil

	case dependencytree.NodeStateAbsent:
		err := c.client.DeleteContainerRepository(c.commonMetadata.Ctx, client.DeleteContainerRepositoryOpts{
			ClusterID: c.commonMetadata.ClusterID,
			ImageName: c.commonMetadata.ApplicationDeclaration.Image.Name,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting container repository: %w", err)
		}

		return reconciliation.Result{}, nil
	}

	return reconciliation.Result{}, nil
}

// NewContainerRepositoryReconciler creates a new reconciler for the VPC resource
func NewContainerRepositoryReconciler(client client.ContainerRepositoryService) reconciliation.Reconciler {
	return &containerRepositoryReconciler{
		client: client,
	}
}
