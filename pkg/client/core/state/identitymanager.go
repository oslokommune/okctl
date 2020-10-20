package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type identityManagerState struct {
	state state.IdentityPooler
}

func (s *identityManagerState) SaveIdentityPoolClient(client *api.IdentityPoolClient) (*store.Report, error) {
	c := s.state.GetIdentityPoolClient(client.Purpose)

	c.Purpose = client.Purpose
	c.CallbackURL = client.CallbackURL
	c.ClientID = client.ClientID

	report, err := s.state.SaveIdentityPoolClient(c)
	if err != nil {
		return nil, fmt.Errorf("saving state: %w", err)
	}

	report.Actions = append([]store.Action{
		{
			Name: "IdentityPoolClient",
			Path: fmt.Sprintf("purpose=%s, client_id=%s", c.Purpose, c.ClientID),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

func (s *identityManagerState) SaveIdentityPool(p *api.IdentityPool) (*store.Report, error) {
	pool := s.state.GetIdentityPool()

	pool.UserPoolID = p.UserPoolID
	pool.AuthDomain = p.AuthDomain
	pool.Alias = state.RecordSetAlias{
		AliasDomain:     p.RecordSetAlias.AliasDomain,
		AliasHostedZone: p.RecordSetAlias.AliasHostedZones,
	}

	report, err := s.state.SaveIdentityPool(pool)
	if err != nil {
		return nil, fmt.Errorf("saving state: %w", err)
	}

	report.Actions = append([]store.Action{
		{
			Name: "IdentityPool",
			Path: fmt.Sprintf("id=%s, url=%s", p.UserPoolID, p.AuthDomain),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

// NewIdentityManagerState returns an initialised state manager
func NewIdentityManagerState(state state.IdentityPooler) client.IdentityManagerState {
	return &identityManagerState{
		state: state,
	}
}
