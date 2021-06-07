package reconciliation

import (
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// applicationReconciler contains service and metadata for the relevant resource
type applicationReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

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

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (a *applicationReconciler) Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		repoID := fmt.Sprintf("%s/%s",
			a.commonMetadata.Declaration.Github.Organisation,
			a.commonMetadata.Declaration.Github.Repository,
		)

		gh, err := state.Github.GetGithubRepository(repoID)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("retrieving Github information")
		}

		if a.commonMetadata.ApplicationDeclaration.Image.HasName() {
			repo, err := state.ContainerRepository.GetContainerRepository(a.commonMetadata.ApplicationDeclaration.Image.Name)
			if err != nil {
				return reconciliation.Result{}, fmt.Errorf("getting container repository: %w", err)
			}

			a.commonMetadata.ApplicationDeclaration.Image.Name = ""
			a.commonMetadata.ApplicationDeclaration.Image.URI = repo.URI()
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
			return reconciliation.Result{}, err
		}
	case resourcetree.ResourceNodeStateAbsent:
		return reconciliation.Result{}, errors.New("deletion of applications is not implemented")
	}

	return reconciliation.Result{}, nil
}

// NewApplicationReconciler creates a new reconciler for the VPC resource
func NewApplicationReconciler(client client.ApplicationService) reconciliation.Reconciler {
	return &applicationReconciler{
		client: client,
	}
}
