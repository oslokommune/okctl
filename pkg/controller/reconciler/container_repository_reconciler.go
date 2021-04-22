package reconciler

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// containerRepositoryReconciler contains service and metadata for the relevant resource
type containerRepositoryReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	handlers       *clientCore.StateHandlers

	client client.ContainerRepositoryService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (c *containerRepositoryReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeContainerRepository
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (c *containerRepositoryReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	c.commonMetadata = metadata
}

// SetStateHandlers sets the state handlers
func (c *containerRepositoryReconciler) SetStateHandlers(handlers *clientCore.StateHandlers) {
	c.handlers = handlers
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (c *containerRepositoryReconciler) Reconcile(node *resourcetree.ResourceNode) (ReconcilationResult, error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		fmt.Println("CREATING CONTAINER REPO")
		if !c.commonMetadata.ApplicationDeclaration.Image.HasName() {
			return ReconcilationResult{}, nil
		}

		_, err := c.handlers.ContainerRepository.GetContainerRepository(c.commonMetadata.ApplicationDeclaration.Image.Name + clustername something) // TODO use contains or something.. what is stackname...
		// aha see StackName: cfn.NewStackNamer().ContainerRepository(opts.ImageName, opts.ClusterID.ClusterName),
		// in service_container_repository_client.go
		// todo adjust this

		// wip comment: err != nil means not found
		if err != nil { // TODO: How do we detect difference in actual error and not found?
			createdRepo, err := c.client.CreateContainerRepository(c.commonMetadata.Ctx, client.CreateContainerRepositoryOpts{
				ClusterID: c.commonMetadata.ClusterID,
				ImageName: c.commonMetadata.ApplicationDeclaration.Image.Name,
			})
			if err != nil {
				return ReconcilationResult{}, fmt.Errorf("creating container repository: %w", err)
			}

			err = c.handlers.ContainerRepository.SaveContainerRepository(createdRepo)
			if err != nil {
			    return ReconcilationResult{}, fmt.Errorf("saving container repository: %w", err)
			}
		}

		return ReconcilationResult{}, nil

	case resourcetree.ResourceNodeStateAbsent:
		err := c.client.DeleteContainerRepository(c.commonMetadata.Ctx, client.DeleteContainerRepositoryOpts{
			ClusterID: c.commonMetadata.ClusterID,
			ImageName: c.commonMetadata.ApplicationDeclaration.Image.Name,
		})
		if err != nil {
			return ReconcilationResult{}, fmt.Errorf("deleting container repository: %w", err)
		}

		return ReconcilationResult{}, nil
	}

	return ReconcilationResult{}, nil
}

// NewApplicationReconciler creates a new reconciler for the VPC resource
func NewContainerRepositoryReconciler(client client.ContainerRepositoryService) Reconciler {
	return &containerRepositoryReconciler{
		client: client,
	}
}
