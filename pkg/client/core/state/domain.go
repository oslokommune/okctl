package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type domainState struct {
	state state.HostedZoner
}

func (s *domainState) RemoveHostedZone(domain string) (*store.Report, error) {
	report, err := s.state.DeleteHostedZone(domain)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "HostedZone",
			Path: fmt.Sprintf("domain=%s", domain),
			Type: "StateUpdate[removed]",
		},
	}, report.Actions...)

	return report, nil
}

func (s *domainState) GetHostedZones() (zones []state.HostedZone) {
	for _, z := range s.state.GetHostedZones() {
		zones = append(zones, z)
	}

	return zones
}

func (s *domainState) SaveHostedZone(zone *client.HostedZone) (*store.Report, error) {
	hz := s.state.GetHostedZone(zone.HostedZone.Domain)

	hz.IsDelegated = zone.IsDelegated
	hz.Primary = zone.Primary
	hz.Managed = zone.HostedZone.Managed
	hz.Domain = zone.HostedZone.Domain
	hz.FQDN = zone.HostedZone.FQDN
	hz.NameServers = zone.HostedZone.NameServers

	report, err := s.state.SaveHostedZone(zone.HostedZone.Domain, hz)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "HostedZone",
			Path: fmt.Sprintf("domain=%s, clusterName=%s", zone.HostedZone.Domain, zone.HostedZone.ID.ClusterName),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

// NewDomainState returns a state implementation
func NewDomainState(set state.HostedZoner) client.DomainState {
	return &domainState{
		state: set,
	}
}
