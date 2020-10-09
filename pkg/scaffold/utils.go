package scaffold

import (
	"os"

	"github.com/oslokommune/okctl/pkg/ask"

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

func askForCluster(o *okctl.Okctl) *state.Cluster {
	asker := ask.New()

	environments := make([]*state.Cluster, len(o.RepoState.Clusters))

	index := 0

	for key := range o.RepoState.Clusters {
		current := o.RepoState.Clusters[key]
		environments[index] = &current

		index++
	}

	selectedEnvironment, err := asker.SelectEnvironment(environments)
	if err != nil {
		return nil
	}

	return selectedEnvironment
}

/*
GetCluster returns a certain cluster. If DefaultRepositoryConfig isnt available, GetCluster returns nil.
If env is specified, GetCluster tries to return the relevant cluster.
If only one cluster is available, GetCluster will return that.
If more than one cluster is available, GetCluster will ask which cluster to use
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

	availableClusters := len(o.RepoState.Clusters)

	var cluster state.Cluster

	switch {
	case env != "":
		cluster = o.RepoState.Clusters[env]
	case availableClusters == 0:
		return nil
	case availableClusters == 1:
		for key := range o.RepoState.Clusters {
			cluster = o.RepoState.Clusters[key]

			break
		}
	case availableClusters > 1:
		cluster = *askForCluster(o)
	}

	return &cluster
}

// GetIACRepoURL extracts the infrastructure as code repository URL from a cluster
func GetIACRepoURL(cluster *state.Cluster) string {
	url := ""
	for repo := range cluster.Github.Repositories {
		url = cluster.Github.Repositories[repo].GitURL

		break
	}

	return url
}

// GetHostedZoneDomain returns the first hosted zone domain for a cluster
func GetHostedZoneDomain(cluster *state.Cluster) string {
	if cluster == nil {
		return ""
	}

	var zone state.HostedZone
	for zoneName := range cluster.HostedZone {
		zone = cluster.HostedZone[zoneName]
		break
	}

	return zone.Domain
}
