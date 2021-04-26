package reconciler

import (
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// applicationReconciler contains service and metadata for the relevant resource
type applicationReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	handlers       *clientCore.StateHandlers

	client client.ApplicationService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (a *applicationReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeApplication
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (a *applicationReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	a.commonMetadata = metadata
}

// SetStateHandlers sets the state handlers
func (a *applicationReconciler) SetStateHandlers(handlers *clientCore.StateHandlers) {
	a.handlers = handlers
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (a *applicationReconciler) Reconcile(node *resourcetree.ResourceNode) (ReconcilationResult, error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		hz, err := a.handlers.Domain.GetPrimaryHostedZone()
		if err != nil {
			return ReconcilationResult{}, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		repoID := fmt.Sprintf("%s/%s",
			a.commonMetadata.Declaration.Github.Organisation,
			a.commonMetadata.Declaration.Github.Repository,
		)

		gh, err := a.handlers.Github.GetGithubRepository(repoID)
		if err != nil {
			return ReconcilationResult{}, fmt.Errorf("retrieving Github information")
		}

		if a.commonMetadata.ApplicationDeclaration.Image.HasName() {
			repo, err := a.handlers.ContainerRepository.GetContainerRepository(a.commonMetadata.ApplicationDeclaration.Image.Name)
			if err != nil {
				return ReconcilationResult{}, fmt.Errorf("getting container repository: %w", err)
			}

			repoURI := repo.URI()
			a.commonMetadata.ApplicationDeclaration.Image.URI = repoURI.String()
		}

		err = a.client.ScaffoldApplication(a.commonMetadata.Ctx, &client.ScaffoldApplicationOpts{
			OutputDir:        a.commonMetadata.Declaration.Github.OutputPath,
			ID:               &a.commonMetadata.ClusterID,
			HostedZoneID:     hz.HostedZoneID,
			HostedZoneDomain: hz.Domain,
			IACRepoURL:       gh.GitURL,
			Application:      a.commonMetadata.ApplicationDeclaration,
		})
		if err != nil {
			return ReconcilationResult{}, err
		}
	case resourcetree.ResourceNodeStateAbsent:
		return ReconcilationResult{}, errors.New("deletion of applications is not implemented")
	}

	return ReconcilationResult{}, nil
}

// NewApplicationReconciler creates a new reconciler for the VPC resource
func NewApplicationReconciler(client client.ApplicationService) Reconciler {
	return &applicationReconciler{
		client: client,
	}
}
