package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/keypair"

	"github.com/oslokommune/okctl/pkg/github"

	stormpkg "github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/client"
)

type githubService struct {
	parameterService api.ParameterService
	githubClient     github.Github
	state            client.GithubState
}

func githubDeployKeySecretName(org, repo string) string {
	return fmt.Sprintf("github/deploykeys/%s/%s/privatekey", org, repo)
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

	err = s.DeleteRepositoryDeployKey(client.DeleteGithubDeployKeyOpts{
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

	if r != nil && r.Validate() == nil {
		return r, nil
	}

	key, err := s.CreateRepositoryDeployKey(client.CreateGithubDeployKeyOpts{
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

func (s *githubService) DeleteRepositoryDeployKey(opts client.DeleteGithubDeployKeyOpts) error {
	err := s.parameterService.DeleteSecret(context.Background(), api.DeleteSecretOpts{
		ID:   opts.ID,
		Name: githubDeployKeySecretName(opts.Organisation, opts.Repository),
	})
	if err != nil {
		return err
	}

	return s.githubClient.DeleteDeployKey(opts.Organisation, opts.Repository, opts.Identifier)
}

func (s *githubService) ListReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	return s.githubClient.ListReleases(owner, repo)
}

func (s *githubService) CreateRepositoryDeployKey(opts client.CreateGithubDeployKeyOpts) (*client.GithubDeployKey, error) {
	key, err := keypair.Generate()
	if err != nil {
		return nil, err
	}

	param, err := s.parameterService.CreateSecret(context.Background(), api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   githubDeployKeySecretName(opts.Organisation, opts.Repository),
		Secret: string(key.PrivateKey),
	})
	if err != nil {
		return nil, err
	}

	deployKey, err := s.githubClient.CreateDeployKey(opts.Organisation, opts.Repository, opts.Title, string(key.PublicKey))
	if err != nil {
		return nil, err
	}

	return &client.GithubDeployKey{
		Organisation: opts.Organisation,
		Repository:   opts.Repository,
		Identifier:   deployKey.GetID(),
		Title:        deployKey.GetTitle(),
		PublicKey:    string(key.PublicKey),
		PrivateKeySecret: &client.GithubSecret{
			Name:    param.Name,
			Path:    param.Path,
			Version: param.Version,
		},
	}, nil
}

// NewGithubService returns an initialised service
func NewGithubService(parameterService api.ParameterService, githubClient github.Github, state client.GithubState) client.GithubService {
	return &githubService{
		parameterService: parameterService,
		githubClient:     githubClient,
		state:            state,
	}
}
