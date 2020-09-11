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
	report, err := s.state.SaveArgoCD(&state.ArgoCD{
		SiteURL: cd.ArgoURL,
		Domain:  cd.ArgoDomain,
		SecretKey: &state.SecretKeySecret{
			Name:    cd.SecretKey.Name,
			Path:    cd.SecretKey.Path,
			Version: cd.SecretKey.Version,
		},
	})
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
