package scaffold

import (
	"os"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func loadRepoData(o *okctl.Okctl, cmd *cobra.Command) error {
	repoDataNotFound := load.CreateOnRepoDataNotFound()

	if o.NoInput {
		repoDataNotFound = load.ErrOnRepoDataNotFound()
	}

	o.RepoDataLoader = load.RepoDataFromConfigFile(cmd, repoDataNotFound)

	return o.LoadRepoData()
}

/*
GetCluster returns a certain cluster if env is specified and DefaultRepositoryConfig is available.
If no env is specified and DefaultRepositoryConfig is available, GetCluster will return the first
cluster it finds.

If theres no DefaultRepositoryConfig available, GetCluster returns nil
*/
func GetCluster(o *okctl.Okctl, cmd *cobra.Command, env string) *state.Cluster {
	_, err := os.Stat(config.DefaultRepositoryConfig)
	if os.IsNotExist(err) {
		return nil
	}

	err = loadRepoData(o, cmd)
	if err != nil {
		return nil
	}

	var cluster state.Cluster

	if env == "" {
		for item := range o.RepoState.Clusters {
			cluster = o.RepoState.Clusters[item]

			break
		}
	} else {
		cluster = o.RepoState.Clusters[env]
	}

	return &cluster
}
