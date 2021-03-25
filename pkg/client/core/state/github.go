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

func (s *githubState) GetRepositoryDeployKey(_ api.ID) state.GithubRepository {
	github := s.state.GetGithub()

	for _, repo := range github.Repositories {
		for _, t := range repo.Types {
			if t == state.TypeInfrastructure {
				return repo
			}
		}
	}

	return state.GithubRepository{}
}

func (s *githubState) SaveRepositoryDeployKey(r *client.GithubRepository) (*store.Report, error) {
	github := s.state.GetGithub()

	for _, repo := range github.Repositories {
		for _, t := range repo.Types {
			if t == state.TypeInfrastructure && repo.FullName == r.FullName {
				return nil, fmt.Errorf("cluster already has an infrastructure repository: %s", repo.FullName)
			}
		}
	}

	local := state.GithubRepository{
		Name:         r.Repository,
		FullName:     r.FullName,
		Organization: r.Organisation,
		Types:        []string{state.TypeInfrastructure},
		GitURL:       r.GitURL,
		DeployKey: state.DeployKey{
			ID:        r.DeployKey.Identifier,
			Title:     r.DeployKey.Title,
			PublicKey: r.DeployKey.PublicKey,
			PrivateKeySecret: state.PrivateKeySecret{
				Name:    r.DeployKey.PrivateKeySecret.Name,
				Path:    r.DeployKey.PrivateKeySecret.Path,
				Version: r.DeployKey.PrivateKeySecret.Version,
			},
		},
	}

	if github.Repositories == nil {
		github.Repositories = map[string]state.GithubRepository{}
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

// NewGithubState returns an initialised state handler
func NewGithubState(state state.Githuber) client.GithubState {
	return &githubState{
		state: state,
	}
}
