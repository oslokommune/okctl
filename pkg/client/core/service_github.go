package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client"
)

type githubService struct {
	api    client.GithubAPI
	report client.GithubReport
	state  client.GithubState
}

func (s *githubService) ReadyGithubInfrastructureRepository(_ context.Context, opts client.ReadyGithubInfrastructureRepositoryOpts) (*client.GithubRepository, error) {
	err := opts.Validate()
	if err != nil {
		return nil, err
	}

	repository := s.state.GetGithubInfrastructureRepository(opts.ID)
	if repository != nil && repository.Validate() == nil {
		return repository, nil
	}

	selected, err := s.api.SelectGithubInfrastructureRepository(client.SelectGithubInfrastructureRepositoryOpts(opts))
	if err != nil {
		return nil, err
	}

	key, err := s.api.CreateGithubDeployKey(client.CreateGithubDeployKey{
		ID:           opts.ID,
		Organisation: opts.Organisation,
		Repository:   selected.Repository,
	})
	if err != nil {
		return nil, err
	}

	repository = &client.GithubRepository{
		ID:           selected.ID,
		Organisation: selected.Organisation,
		Repository:   selected.Repository,
		GitURL:       selected.GitURL,
		FullName:     selected.FullName,
		DeployKey:    key,
	}

	report, err := s.state.SaveGithubInfrastructureRepository(repository)
	if err != nil {
		return nil, err
	}

	err = s.report.ReadyGithubInfrastructureRepository(repository, report)
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (s *githubService) CreateGithubOauthApp(_ context.Context, opts client.CreateGithubOauthAppOpts) (*client.GithubOauthApp, error) {
	err := opts.Validate()
	if err != nil {
		return nil, err
	}

	app := s.state.GetGithubOauthApp(opts.Name, opts.ID)
	if app != nil && app.Validate() == nil {
		return app, nil
	}

	team, err := s.api.SelectGithubTeam(client.SelectGithubTeam{
		ID:           opts.ID,
		Organisation: opts.Organisation,
	})
	if err != nil {
		return nil, err
	}

	app, err = s.api.CreateGithubOauthApp(client.CreateGithubOauthAppOpts{
		ID:           opts.ID,
		Organisation: opts.Organisation,
		Team:         team,
		Name:         opts.Name,
		SiteURL:      opts.SiteURL,
		CallbackURL:  opts.CallbackURL,
	})
	if err != nil {
		return nil, err
	}

	report, err := s.state.SaveGithubOauthApp(app)
	if err != nil {
		return nil, err
	}

	err = s.report.CreateGithubOauthApp(app, report)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// NewGithubService returns an initialised service
func NewGithubService(api client.GithubAPI, report client.GithubReport, state client.GithubState) client.GithubService {
	return &githubService{
		api:    api,
		report: report,
		state:  state,
	}
}
