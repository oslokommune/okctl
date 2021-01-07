package state

import (
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type argoCDState struct {
	state state.RepositoryStateWithEnv
}

func (s *argoCDState) SaveArgoCD(cd *client.ArgoCD) (*store.Report, error) {
	argo := s.state.GetArgoCD()

	argo.SiteURL = cd.ArgoURL
	argo.Domain = cd.ArgoDomain
	argo.SecretKey.Version = cd.SecretKey.Version
	argo.SecretKey.Path = cd.SecretKey.Path
	argo.SecretKey.Name = cd.SecretKey.Name

	report, err := s.state.SaveArgoCD(argo)
	if err != nil {
		return nil, err
	}

	c := s.state.GetIdentityPoolClient(cd.IdentityClient.Purpose)

	c.ClientSecret.Name = cd.ClientSecret.Name
	c.ClientSecret.Version = cd.ClientSecret.Version
	c.ClientSecret.Path = cd.ClientSecret.Path

	r2, err := s.state.SaveIdentityPoolClient(c)
	if err != nil {
		return nil, err
	}

	report.Actions = append(report.Actions, r2.Actions...)

	report.Actions = append([]store.Action{
		{
			Name: "ArgoCD",
			Path: "cluster=" + cd.ID.ClusterName,
			Type: "StateUpdate[add]",
		},
		{
			Name: "IdentityPoolClient",
			Path: "purpose=" + c.Purpose,
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

// NewArgoCDState returns an initialised state layer
func NewArgoCDState(state state.RepositoryStateWithEnv) client.ArgoCDState {
	return &argoCDState{
		state: state,
	}
}
