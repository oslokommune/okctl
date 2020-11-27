package reconsiler

import (
	"fmt"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// GithubMetadata contains data from the desired state
type GithubMetadata struct {
	Organization string
	Repository string
}

// GithubGetter knows how to get the current state Github
type GithubGetter func() state.Github
// GithubSetter knows how to save a state.Github
type GithubSetter func(github state.Github) (*store.Report, error)

// GithubResourceState contains runtime data needed in Reconsile()
type GithubResourceState struct {
	Getter GithubGetter
	Saver GithubSetter
}

type githubReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.GithubService
}

// SetCommonMetadata saves common metadata for use in Reconsile()
func (z *githubReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconsile knows how to ensure the desired state is achieved
func (z *githubReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	metadata, ok := node.Metadata.(GithubMetadata)
	if !ok {
		return nil, errors.New("unable to cast Github metadata")
	}

	resourceState, ok := node.ResourceState.(GithubResourceState)
	if !ok {
		return nil, errors.New("unable to cast Github resource state")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.ReadyGithubInfrastructureRepositoryWithoutUserinput(z.commonMetadata.Ctx, client.ReadyGithubInfrastructureRepositoryOpts{
			ID:           z.commonMetadata.Id,
			Organisation: metadata.Organization,
			Repository:   metadata.Repository,
		})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating Github resource: %w", err)
		}

		gh := resourceState.Getter()
		gh.Organisation = metadata.Organization

		_, err = resourceState.Saver(gh)
		if err != nil {
		    return nil, fmt.Errorf("error saving github: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return nil, errors.New("deleting Github resource is not implemented")
	}

	return &ReconsilationResult{Requeue: false}, nil
}

// NewGithubReconsiler creates a new reconsiler for the Github resource
func NewGithubReconsiler(client client.GithubService) *githubReconsiler {
	return &githubReconsiler{
		client: client,
	}
}
