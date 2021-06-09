package reconciliation

import (
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

// applicationReconciler contains service and metadata for the relevant resource
type applicationReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.ApplicationService
}

// NodeType returns the relevant NodeType for this reconciler
func (a *applicationReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeApplication
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (a *applicationReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	a.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (a *applicationReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
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
	case dependencytree.NodeStateAbsent:
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
