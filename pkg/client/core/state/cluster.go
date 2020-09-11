package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type clusterState struct {
	state state.Clusterer
}

func (s *clusterState) SaveCluster(c *api.Cluster) (*store.Report, error) {
	cluster := s.state.GetCluster()

	cluster.Name = c.ID.ClusterName
	cluster.Environment = c.ID.Environment
	cluster.AWSAccountID = c.ID.AWSAccountID

	report, err := s.state.SaveCluster(cluster)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "Cluster",
			Path: fmt.Sprintf("clusterName=%s", c.ID.ClusterName),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

func (s *clusterState) DeleteCluster(_ api.ID) (*store.Report, error) {
	return s.state.DeleteCluster()
}

// NewClusterState returns an initialised state handler
func NewClusterState(state state.Clusterer) client.ClusterState {
	return &clusterState{
		state: state,
	}
}
