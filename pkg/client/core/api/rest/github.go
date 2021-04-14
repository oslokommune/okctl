package rest

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/keypair"
)

type githubAPI struct {
	client       github.Githuber
	parameterAPI client.ParameterAPI
}

func githubDeployKeySecretName(org, repo string) string {
	return fmt.Sprintf("github/deploykeys/%s/%s/privatekey", org, repo)
}

func (a *githubAPI) DeleteRepositoryDeployKey(opts client.DeleteGithubDeployKeyOpts) error {
	err := a.parameterAPI.DeleteSecret(api.DeleteSecretOpts{
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

	param, err := a.parameterAPI.CreateSecret(api.CreateSecretOpts{
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

// NewGithubAPI returns an instantiated github API client
func NewGithubAPI(paramAPI client.ParameterAPI, client github.Githuber) client.GithubAPI {
	return &githubAPI{
		client:       client,
		parameterAPI: paramAPI,
	}
}
