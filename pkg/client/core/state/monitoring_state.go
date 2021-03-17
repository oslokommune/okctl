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

func (s *monitoringState) GetKubePromStack() (*client.KubePromStack, error) {
	m := s.state.GetMonitoring()

	return &client.KubePromStack{
		FargateCloudWatchPolicyARN:        m.FargateCloudWatchPolicyARN,
		FargateProfilePodExecutionRoleARN: m.FargateProfilePodExecutionRoleARN,
	}, nil
}

func (s *monitoringState) RemoveKubePromStack(_ api.ID) (*store.Report, error) {
	_ = s.state.GetMonitoring()

	return s.state.SaveMonitoring(state.Monitoring{})
}

func (s *monitoringState) SaveKubePromStack(kube *client.KubePromStack) (*store.Report, error) {
	m := s.state.GetMonitoring()
	m.DashboardURL = fmt.Sprintf("https://%s", kube.Hostname)
	m.FargateProfilePodExecutionRoleARN = kube.FargateProfilePodExecutionRoleARN
	m.FargateCloudWatchPolicyARN = kube.FargateCloudWatchPolicyARN

	return s.state.SaveMonitoring(m)
}

// NewMonitoringState returns an initialised state layer
func NewMonitoringState(state state.RepositoryStateWithEnv) client.MonitoringState {
	return &monitoringState{
		state: state,
	}
}
