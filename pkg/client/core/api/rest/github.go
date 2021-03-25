package rest

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/keypair"
)

type githubAPI struct {
	client       github.Githuber
	parameterAPI client.ParameterAPI
	ask          *ask.Ask
	out          io.Writer
}

func (a *githubAPI) CreateRepositoryDeployKey(opts client.CreateGithubDeployKeyOpts) (*client.GithubDeployKey, error) {
	key, err := keypair.New(keypair.DefaultRandReader(), keypair.DefaultBitSize).Generate()
	if err != nil {
		return nil, err
	}

	param, err := a.parameterAPI.CreateSecret(api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   fmt.Sprintf("github/deploykeys/%s/%s/privatekey", opts.Organisation, opts.Repository),
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
		ID:           opts.ID,
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
func NewGithubAPI(out io.Writer, ask *ask.Ask, paramAPI client.ParameterAPI, client github.Githuber) client.GithubAPI {
	return &githubAPI{
		client:       client,
		parameterAPI: paramAPI,
		ask:          ask,
		out:          out,
	}
}
