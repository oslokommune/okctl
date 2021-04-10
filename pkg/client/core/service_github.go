package core

import (
	"context"
	"errors"
	"fmt"

	stormpkg "github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/client"
)

type githubService struct {
	api   client.GithubAPI
	state client.GithubState
}

func (s *githubService) DeleteGithubRepository(_ context.Context, opts client.DeleteGithubRepositoryOpts) error {
	fullName := fmt.Sprintf("%s/%s", opts.Organisation, opts.Name)

	r, err := s.state.GetGithubRepository(fullName)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	err = s.api.DeleteRepositoryDeployKey(client.DeleteGithubDeployKeyOpts{
		ID:           opts.ID,
		Organisation: opts.Organisation,
		Repository:   opts.Name,
		Identifier:   r.DeployKey.Identifier,
	})
	if err != nil {
		return err
	}

	return s.state.RemoveGithubRepository(fullName)
}

// CreateGithubRepository creates a new github repository for a certain repository
// if it doesn't exist. If it exists, it returns the existing repository
func (s *githubService) CreateGithubRepository(_ context.Context, opts client.CreateGithubRepositoryOpts) (*client.GithubRepository, error) {
	fullName := fmt.Sprintf("%s/%s", opts.Organization, opts.Name)

	r, err := s.state.GetGithubRepository(fullName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return nil, err
	}

	if r.Validate() == nil {
		return r, nil
	}

	key, err := s.api.CreateRepositoryDeployKey(client.CreateGithubDeployKeyOpts{
		ID:           opts.ID,
		Organisation: opts.Organization,
		Repository:   opts.Name,
		Title:        fmt.Sprintf("okctl-iac-%s", opts.ID.ClusterName),
	})
	if err != nil {
		return nil, err
	}

	repo := &client.GithubRepository{
		ID:           opts.ID,
		Organisation: opts.Organization,
		Repository:   opts.Name,
		FullName:     fullName,
		GitURL:       fmt.Sprintf("%s:%s", opts.Host, fullName),
		DeployKey:    key,
	}

	err = s.state.SaveGithubRepository(repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// NewGithubService returns an initialised service
func NewGithubService(api client.GithubAPI, state client.GithubState) client.GithubService {
	return &githubService{
		api:   api,
		state: state,
	}
}
