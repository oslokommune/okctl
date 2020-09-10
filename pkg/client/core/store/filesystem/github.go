package filesystem

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type githubStore struct {
	repoState *state.Repository
	repoPaths Paths
	fs        *afero.Afero
}

func (s githubStore) SaveGithubInfrastructureRepository(r *client.GithubRepository) (*store.Report, error) {
	cluster, ok := s.repoState.Clusters[r.ID.Environment]
	if !ok {
		return nil, fmt.Errorf("no cluster found for environment: %s", r.ID.Environment)
	}

	for _, repo := range cluster.Github.Repositories {
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

	if cluster.Github.Repositories == nil {
		cluster.Github.Repositories = map[string]*state.GithubRepository{}
	}

	cluster.Github.Repositories[local.FullName] = local

	report, err := store.NewFileSystem(s.repoPaths.BaseDir, s.fs).
		StoreStruct(s.repoPaths.ConfigFile, s.repoState, store.ToYAML()).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s githubStore) GetGithubInfrastructureRepository(id api.ID) (*client.GithubRepository, error) {
	cluster, ok := s.repoState.Clusters[id.Environment]
	if !ok {
		return nil, fmt.Errorf("no cluster found for environment: %s", id.Environment)
	}

	for _, repo := range cluster.Github.Repositories {
		for _, t := range repo.Types {
			if t == state.TypeInfrastructure {
				return &client.GithubRepository{
					ID:           id,
					Organisation: cluster.Github.Organisation,
					Repository:   repo.Name,
					FullName:     repo.FullName,
					GitURL:       repo.GitURL,
					DeployKey: &client.GithubDeployKey{
						ID:           id,
						Organisation: cluster.Github.Organisation,
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
				}, nil
			}
		}
	}

	return nil, nil
}

func (s githubStore) SaveGithubOauthApp(app *client.GithubOauthApp) (*store.Report, error) {
	cluster, ok := s.repoState.Clusters[app.ID.Environment]
	if !ok {
		return nil, fmt.Errorf("no cluster found for environment: %s", app.ID.Environment)
	}

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

	if cluster.Github.OauthApp == nil {
		cluster.Github.OauthApp = map[string]*state.GithubOauthApp{}
	}

	cluster.Github.OauthApp[app.Name] = local

	report, err := store.NewFileSystem(s.repoPaths.BaseDir, s.fs).
		StoreStruct(s.repoPaths.ConfigFile, s.repoState, store.ToYAML()).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s githubStore) GetGithubOauthApp(appName string, id api.ID) (*client.GithubOauthApp, error) {
	cluster, ok := s.repoState.Clusters[id.Environment]
	if !ok {
		return nil, fmt.Errorf("no cluster found for environment: %s", id.Environment)
	}

	if app, ok := cluster.Github.OauthApp[appName]; ok {
		return &client.GithubOauthApp{
			ID:           id,
			Organisation: cluster.Github.Organisation,
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
				Organisation: cluster.Github.Organisation,
				Name:         app.Team,
			},
		}, nil
	}

	return nil, nil
}

// NewGithubStore returns an initialised store
func NewGithubStore(repoPaths Paths, repoState *state.Repository, fs *afero.Afero) client.GithubStore {
	return &githubStore{
		repoState: repoState,
		repoPaths: repoPaths,
		fs:        fs,
	}
}
