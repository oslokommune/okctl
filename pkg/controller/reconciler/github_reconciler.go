package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// githubReconciler contains service and metadata for the relevant resource
type githubReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.GithubService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *githubReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *githubReconciler) Reconcile(node *resourcetree.ResourceNode) (result *ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateDeployKey(z.commonMetadata.Ctx, client.NewGithubRepository(
			z.commonMetadata.ClusterID,
			config.DefaultGithubHost,
			z.commonMetadata.Declaration.Github.Organisation,
			z.commonMetadata.Declaration.Github.Repository,
		))
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating Github deploy key: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return nil, errors.New("deleting Github resource is not implemented")
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewGithubReconciler creates a new reconciler for the Github resource
func NewGithubReconciler(client client.GithubService) Reconciler {
	return &githubReconciler{
		client: client,
	}
}
