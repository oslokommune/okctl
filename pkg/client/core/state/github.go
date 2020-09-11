package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type githubState struct {
	state state.Githuber
}

func (s *githubState) GetGithubInfrastructureRepository(id api.ID) *client.GithubRepository {
	github := s.state.GetGithub()

	for _, repo := range github.Repositories {
		for _, t := range repo.Types {
			if t == state.TypeInfrastructure {
				return &client.GithubRepository{
					ID:           id,
					Organisation: github.Organisation,
					Repository:   repo.Name,
					FullName:     repo.FullName,
					GitURL:       repo.GitURL,
					DeployKey: &client.GithubDeployKey{
						ID:           id,
						Organisation: github.Organisation,
						Repository:   repo.Name,
						Identifier:   repo.DeployKey.ID,
						Title:        repo.DeployKey.Title,
						PublicKey:    repo.DeployKey.PublicKey,
						PrivateKeySecret: &client.GithubSecret{
							Name:    repo.DeployKey.PrivateKeySecret.Name,
							Path:    repo.DeployKey.PrivateKeySecret.Path,
							Version: repo.DeployKey.PrivateKeySecret.Version,
						},
					},
				}
			}
		}
	}

	return nil
}

func (s *githubState) GetGithubOauthApp(appName string, id api.ID) *client.GithubOauthApp {
	github := s.state.GetGithub()

	if app, ok := github.OauthApp[appName]; ok {
		return &client.GithubOauthApp{
			ID:           id,
			Organisation: github.Organisation,
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
				ID:           id,
				Organisation: github.Organisation,
				Name:         app.Team,
			},
		}
	}

	return nil
}

func (s *githubState) SaveGithubInfrastructureRepository(r *client.GithubRepository) (*store.Report, error) {
	github := s.state.GetGithub()

	for _, repo := range github.Repositories {
		for _, t := range repo.Types {
			if t == state.TypeInfrastructure && repo.FullName == r.FullName {
				return nil, fmt.Errorf("cluster already has an infrastructure repository: %s", repo.FullName)
			}
		}
	}

	local := &state.GithubRepository{
		Name:     r.Repository,
		FullName: r.FullName,
		Types:    []string{state.TypeInfrastructure},
		GitURL:   r.GitURL,
		DeployKey: &state.DeployKey{
			ID:        r.DeployKey.Identifier,
			Title:     r.DeployKey.Title,
			PublicKey: r.DeployKey.PublicKey,
			PrivateKeySecret: &state.PrivateKeySecret{
				Name:    r.DeployKey.PrivateKeySecret.Name,
				Path:    r.DeployKey.PrivateKeySecret.Path,
				Version: r.DeployKey.PrivateKeySecret.Version,
			},
		},
	}

	if github.Repositories == nil {
		github.Repositories = map[string]*state.GithubRepository{}
	}

	github.Repositories[local.FullName] = local

	report, err := s.state.SaveGithub(github)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "GithubRepository",
			Path: fmt.Sprintf("repository=%s, deploykey=%s, clusterName=%s",
				r.FullName,
				r.DeployKey.Title,
				r.ID.ClusterName,
			),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

func (s *githubState) SaveGithubOauthApp(app *client.GithubOauthApp) (*store.Report, error) {
	github := s.state.GetGithub()

	local := &state.GithubOauthApp{
		Team:        app.Team.Name,
		Name:        app.Name,
		SiteURL:     app.SiteURL,
		CallbackURL: app.CallbackURL,
		ClientID:    app.ClientID,
		ClientSecret: &state.ClientSecret{
			Name:    app.ClientSecret.Name,
			Path:    app.ClientSecret.Path,
			Version: app.ClientSecret.Version,
		},
	}

	if github.OauthApp == nil {
		github.OauthApp = map[string]*state.GithubOauthApp{}
	}

	github.OauthApp[app.Name] = local

	report, err := s.state.SaveGithub(github)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "GithubOauthApp",
			Path: fmt.Sprintf("name=%s, clusterName=%s", app.Name, app.ID.ClusterName),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

// NewGithubState returns an initialised state handler
func NewGithubState(state state.Githuber) client.GithubState {
	return &githubState{
		state: state,
	}
}
