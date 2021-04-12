package reconciler

import (
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ApplicationState contains information only available at runtime
type ApplicationState struct {
	Declaration client.OkctlApplication

	PrimaryHostedZoneID     string
	PrimaryHostedZoneDomain string
}

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
func (a *applicationReconciler) Reconcile(node *resourcetree.ResourceNode) (ReconcilationResult, error) {
	state, ok := node.ResourceState.(*ApplicationState)
	if !ok {
		return ReconcilationResult{}, errors.New("casting application state")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		err := a.client.ScaffoldApplication(a.commonMetadata.Ctx, &client.ScaffoldApplicationOpts{
			OutputDir:        a.commonMetadata.Declaration.Github.OutputPath,
			ID:               &a.commonMetadata.ClusterID,
			HostedZoneID:     state.PrimaryHostedZoneID,
			HostedZoneDomain: state.PrimaryHostedZoneDomain,
			IACRepoURL:       a.commonMetadata.Declaration.Github.Repository,
			Application:      state.Declaration,
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
