package state

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type domainState struct {
	state state.HostedZoner
}

func (s *domainState) GetHostedZones() (zones []*client.HostedZone) {
	for _, z := range s.state.GetHostedZones() {
		zones = append(zones, &client.HostedZone{
			IsDelegated: z.IsDelegated,
			Primary:     z.Primary,
			HostedZone: &api.HostedZone{
				FQDN:        z.FQDN,
				Domain:      z.Domain,
				NameServers: z.NameServers,
			},
		})
	}

	return zones
}

func (s *domainState) SaveHostedZone(zone *client.HostedZone) (*store.Report, error) {
	return s.state.SaveHostedZone(zone.HostedZone.Domain, &state.HostedZone{
		IsDelegated: zone.IsDelegated,
		Primary:     zone.Primary,
		Domain:      zone.HostedZone.Domain,
		FQDN:        zone.HostedZone.FQDN,
		NameServers: zone.HostedZone.NameServers,
	})
}

// NewDomainState returns a state implementation
func NewDomainState(set state.HostedZoner) client.DomainState {
	return &domainState{
		state: set,
	}
}
