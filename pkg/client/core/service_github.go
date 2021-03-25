package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/client"
)

type githubService struct {
	spinner spinner.Spinner
	api     client.GithubAPI
	report  client.GithubReport
	state   client.GithubState
}

// CreateDeployKey creates a new deploy key for a certain repository if it doesn't exist. If it exists, it returns the
// existing key
func (s *githubService) CreateRepositoryDeployKey(_ context.Context, repository *client.GithubRepository) (key *client.GithubDeployKey, err error) {
	existingRepository := s.state.GetRepositoryDeployKey(repository.ID)

	if existingRepository.Validate() == nil && existingRepository.DeployKey.Validate() == nil {
		return &client.GithubDeployKey{
			ID:           repository.ID,
			Organisation: repository.Organisation,
			Repository:   repository.Repository,
			Identifier:   existingRepository.DeployKey.ID,
			Title:        existingRepository.DeployKey.Title,
			PublicKey:    existingRepository.DeployKey.PublicKey,
			PrivateKeySecret: &client.GithubSecret{
				Name:    existingRepository.DeployKey.PrivateKeySecret.Name,
				Path:    existingRepository.DeployKey.PrivateKeySecret.Path,
				Version: existingRepository.DeployKey.PrivateKeySecret.Version,
			},
		}, nil
	}

	key, err = s.api.CreateRepositoryDeployKey(client.CreateGithubDeployKeyOpts{
		ID:           repository.ID,
		Organisation: repository.Organisation,
		Repository:   repository.Repository,
		Title:        fmt.Sprintf("okctl-iac-%s", repository.ID.ClusterName),
	})
	if err != nil {
		return nil, err
	}

	repository.DeployKey = key

	report, err := s.state.SaveRepositoryDeployKey(repository)
	if err != nil {
		return nil, fmt.Errorf("saving repository state: %w", err)
	}

	err = s.report.ReportRepositoryDeployKey(repository, report)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// NewGithubService returns an initialised service
func NewGithubService(spinner spinner.Spinner, api client.GithubAPI, report client.GithubReport, state client.GithubState) client.GithubService {
	return &githubService{
		spinner: spinner,
		api:     api,
		report:  report,
		state:   state,
	}
}
