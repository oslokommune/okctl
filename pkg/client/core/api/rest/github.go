package rest

import (
	"fmt"
	"io"
	"strings"

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

func (a *githubAPI) GetGithubTeamMembers(opts client.GetGithubTeamMembers) (*[]client.GithubTeamMember, error) {
	members, err := a.client.ListTeamMembers(*opts.Team)
	if err != nil {
		return nil, err
	}

	return &members, nil
}

func (a *githubAPI) SelectGithubInfrastructureRepository(opts client.SelectGithubInfrastructureRepositoryOpts) (*client.SelectedGithubRepository, error) {
	repos, err := a.client.Repositories(opts.Organisation)
	if err != nil {
		return nil, err
	}

	repo, err := a.ask.SelectInfrastructureRepository(opts.Repository, repos)
	if err != nil {
		return nil, err
	}

	return &client.SelectedGithubRepository{
		ID:           opts.ID,
		Organisation: opts.Organisation,
		Repository:   repo.GetName(),
		FullName:     repo.GetFullName(),
		GitURL:       fmt.Sprintf("git@github.com:%s", strings.TrimPrefix(repo.GetGitURL(), "git://github.com/")),
	}, nil
}

func (a *githubAPI) CreateGithubDeployKey(opts client.CreateGithubDeployKey) (*client.GithubDeployKey, error) {
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

func (a *githubAPI) SelectGithubTeam(opts client.SelectGithubTeam) (*client.GithubTeam, error) {
	teams, err := a.client.Teams(opts.Organisation)
	if err != nil {
		return nil, err
	}

	team, err := a.ask.SelectTeam(teams)
	if err != nil {
		return nil, err
	}

	return &client.GithubTeam{
		ID:           opts.ID,
		Organisation: opts.Organisation,
		Name:         team.GetName(),
		TeamID:       *team.ID,
	}, nil
}

func (a *githubAPI) CreateGithubOauthApp(opts client.CreateGithubOauthAppOpts) (*client.GithubOauthApp, error) {
	app, err := a.ask.CreateOauthApp(a.out, ask.OauthAppOpts{
		Organisation: opts.Organisation,
		Name:         opts.Name,
		URL:          opts.SiteURL,
		CallbackURL:  opts.CallbackURL,
	})
	if err != nil {
		return nil, err
	}

	param, err := a.parameterAPI.CreateSecret(api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   fmt.Sprintf("github/oauthapps/%s/%s/clientsecret", opts.Organisation, opts.Name),
		Secret: app.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	return &client.GithubOauthApp{
		ID:           opts.ID,
		Organisation: opts.Organisation,
		Name:         opts.Name,
		SiteURL:      opts.SiteURL,
		CallbackURL:  opts.CallbackURL,
		ClientID:     app.ClientID,
		ClientSecret: &client.GithubSecret{
			Name:    param.Name,
			Path:    param.Path,
			Version: param.Version,
		},
		Team: opts.Team,
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
