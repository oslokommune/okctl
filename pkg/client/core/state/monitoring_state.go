package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type monitoringState struct {
	state state.RepositoryStateWithEnv
}

func (s *monitoringState) RemoveKubePromStack(_ api.ID) (*store.Report, error) {
	m := s.state.GetMonitoring()
	m.DashboardURL = ""

	return s.state.SaveMonitoring(m)
}

func (s *monitoringState) SaveKubePromStack(kube *client.KubePromStack) (*store.Report, error) {
	m := s.state.GetMonitoring()
	m.DashboardURL = fmt.Sprintf("https://%s", kube.Hostname)

	return s.state.SaveMonitoring(m)
}

// NewMonitoringState returns an initialised state layer
func NewMonitoringState(state state.RepositoryStateWithEnv) client.MonitoringState {
	return &monitoringState{
		state: state,
	}
}
