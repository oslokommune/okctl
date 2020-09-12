package state

import (
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type argoCDState struct {
	state state.Argocder
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

	report.Actions = append([]store.Action{
		{
			Name: "ArgoCD",
			Path: "cluster=" + cd.ID.ClusterName,
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

// NewArgoCDState returns an initialised state layer
func NewArgoCDState(state state.Argocder) client.ArgoCDState {
	return &argoCDState{
		state: state,
	}
}
