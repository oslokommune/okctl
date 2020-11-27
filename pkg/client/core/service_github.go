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

// nolint: funlen
func (s *githubService) ReadyGithubInfrastructureRepository(_ context.Context, opts client.ReadyGithubInfrastructureRepositoryOpts) (*client.GithubRepository, error) {
	err := s.spinner.Start("github")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = opts.Validate()
	if err != nil {
		return nil, err
	}

	r := s.state.GetGithubInfrastructureRepository(opts.ID)
	if r.Validate() == nil {
		return &client.GithubRepository{
			ID:           opts.ID,
			Organisation: opts.Organisation,
			Repository:   r.Name,
			FullName:     r.FullName,
			GitURL:       r.GitURL,
			DeployKey: &client.GithubDeployKey{
				ID:           opts.ID,
				Organisation: opts.Organisation,
				Repository:   r.Name,
				Identifier:   r.DeployKey.ID,
				Title:        r.DeployKey.Title,
				PublicKey:    r.DeployKey.PublicKey,
				PrivateKeySecret: &client.GithubSecret{
					Name:    r.DeployKey.PrivateKeySecret.Name,
					Path:    r.DeployKey.PrivateKeySecret.Path,
					Version: r.DeployKey.PrivateKeySecret.Version,
				},
			},
		}, nil
	}

	selected, err := s.api.SelectGithubInfrastructureRepository(client.SelectGithubInfrastructureRepositoryOpts(opts))
	if err != nil {
		return nil, err
	}

	key, err := s.api.CreateGithubDeployKey(client.CreateGithubDeployKey{
		ID:           opts.ID,
		Organisation: opts.Organisation,
		Repository:   selected.Repository,
		Title:        fmt.Sprintf("okctl-iac-%s", opts.ID.ClusterName),
	})
	if err != nil {
		return nil, err
	}

	repo := &client.GithubRepository{
		ID:           selected.ID,
		Organisation: selected.Organisation,
		Repository:   selected.Repository,
		GitURL:       selected.GitURL,
		FullName:     selected.FullName,
		DeployKey:    key,
	}

	report, err := s.state.SaveGithubInfrastructureRepository(repo)
	if err != nil {
		return nil, err
	}

	err = s.report.ReadyGithubInfrastructureRepository(repo, report)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *githubService) ReadyGithubInfrastructureRepositoryWithoutUserinput(_ context.Context, opts client.ReadyGithubInfrastructureRepositoryOpts) (*client.GithubRepository, error) {
	err := s.spinner.Start("github")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = opts.Validate()
	if err != nil {
		return nil, err
	}

	r := s.state.GetGithubInfrastructureRepository(opts.ID)
	if r.Validate() == nil {
		return &client.GithubRepository{
			ID:           opts.ID,
			Organisation: opts.Organisation,
			Repository:   r.Name,
			FullName:     r.FullName,
			GitURL:       r.GitURL,
			DeployKey: &client.GithubDeployKey{
				ID:           opts.ID,
				Organisation: opts.Organisation,
				Repository:   r.Name,
				Identifier:   r.DeployKey.ID,
				Title:        r.DeployKey.Title,
				PublicKey:    r.DeployKey.PublicKey,
				PrivateKeySecret: &client.GithubSecret{
					Name:    r.DeployKey.PrivateKeySecret.Name,
					Path:    r.DeployKey.PrivateKeySecret.Path,
					Version: r.DeployKey.PrivateKeySecret.Version,
				},
			},
		}, nil
	}

	selected, err := s.api.GetGithubInfrastructureRepository(client.SelectGithubInfrastructureRepositoryOpts(opts))
	if err != nil {
		return nil, err
	}

	key, err := s.api.CreateGithubDeployKey(client.CreateGithubDeployKey{
		ID:           opts.ID,
		Organisation: opts.Organisation,
		Repository:   selected.Repository,
		Title:        fmt.Sprintf("okctl-iac-%s", opts.ID.ClusterName),
	})
	if err != nil {
		return nil, err
	}

	repo := &client.GithubRepository{
		ID:           selected.ID,
		Organisation: selected.Organisation,
		Repository:   selected.Repository,
		GitURL:       selected.GitURL,
		FullName:     selected.FullName,
		DeployKey:    key,
	}

	report, err := s.state.SaveGithubInfrastructureRepository(repo)
	if err != nil {
		return nil, err
	}

	err = s.report.ReadyGithubInfrastructureRepository(repo, report)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// nolint: funlen
func (s *githubService) CreateGithubOauthApp(_ context.Context, opts client.CreateGithubOauthAppOpts) (*client.GithubOauthApp, error) {
	err := s.spinner.Start("github")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = opts.Validate()
	if err != nil {
		return nil, err
	}

	app := s.state.GetGithubOauthApp(opts.Name, opts.ID)
	if app.Validate() == nil {
		return &client.GithubOauthApp{
			ID:           opts.ID,
			Organisation: opts.Organisation,
			Name:         app.Name,
			SiteURL:      app.SiteURL,
			CallbackURL:  app.CallbackURL,
			ClientID:     app.ClientID,
			ClientSecret: &client.GithubSecret{
				Name:    app.ClientSecret.Name,
				Path:    app.ClientSecret.Path,
				Version: app.ClientSecret.Version,
			},
			Team: &client.GithubTeam{
				ID:           opts.ID,
				Organisation: opts.Organisation,
				Name:         app.Team,
			},
		}, nil
	}

	team, err := s.api.SelectGithubTeam(client.SelectGithubTeam{
		ID:           opts.ID,
		Organisation: opts.Organisation,
	})
	if err != nil {
		return nil, err
	}

	a, err := s.api.CreateGithubOauthApp(client.CreateGithubOauthAppOpts{
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

	report, err := s.state.SaveGithubOauthApp(a)
	if err != nil {
		return nil, err
	}

	err = s.report.CreateGithubOauthApp(a, report)
	if err != nil {
		return nil, err
	}

	return a, nil
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
