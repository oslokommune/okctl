package direct

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/keypair"
)

type githubAPI struct {
	client           github.Githuber
	parameterService api.ParameterService
}

func githubDeployKeySecretName(org, repo string) string {
	return fmt.Sprintf("github/deploykeys/%s/%s/privatekey", org, repo)
}

func (a *githubAPI) DeleteRepositoryDeployKey(opts client.DeleteGithubDeployKeyOpts) error {
	err := a.parameterService.DeleteSecret(context.Background(), api.DeleteSecretOpts{
		Name: githubDeployKeySecretName(opts.Organisation, opts.Repository),
	})
	if err != nil {
		return err
	}

	return a.client.DeleteDeployKey(opts.Organisation, opts.Repository, opts.Identifier)
}

func (a *githubAPI) CreateRepositoryDeployKey(opts client.CreateGithubDeployKeyOpts) (*client.GithubDeployKey, error) {
	key, err := keypair.New(keypair.DefaultRandReader(), keypair.DefaultBitSize).Generate()
	if err != nil {
		return nil, err
	}

	param, err := a.parameterService.CreateSecret(context.Background(), api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   githubDeployKeySecretName(opts.Organisation, opts.Repository),
		Secret: string(key.PrivateKey),
	})
	if err != nil {
		return nil, err
	}

	deployKey, err := a.client.CreateDeployKey(opts.Organisation, opts.Repository, opts.Title, string(key.PublicKey))
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

func (a *githubAPI) ListReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	return a.client.ListReleases(owner, repo)
}

// NewGithubAPI returns an instantiated github API client
func NewGithubAPI(service api.ParameterService, client github.Githuber) client.GithubAPI {
	return &githubAPI{
		client:           client,
		parameterService: service,
	}
}
