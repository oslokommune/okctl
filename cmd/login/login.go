package login

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/login"
	"github.com/oslokommune/okctl/pkg/stager"
	"github.com/oslokommune/okctl/pkg/storage"
)

func BuildLoginCommand() *cobra.Command {
	var selectedCluster string

	repoCfg := &config.RepoConfig{}

	appCfg := &config.AppConfig{}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to AWS (deprecated)",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			var err error

			appCfg, repoCfg, err = config.Load()

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Login(appCfg, repoCfg, selectedCluster)
		},
	}

	cmd.PersistentFlags().StringVar(&selectedCluster, "cluster", "", "The target cluster")

	return cmd
}

func Login(appCfg *config.AppConfig, repoCfg *config.RepoConfig, clusterName string) error {
	store := storage.NewFileSystemStorage(appCfg.BaseDir)

	stagers, err := stager.FromConfig(appCfg.Binaries, appCfg.Host, store)
	if err != nil {
		return err
	}

	for _, s := range stagers {
		err = s.Run()
		if err != nil {
			return err
		}
	}

	var cluster config.Cluster

	for _, c := range repoCfg.Clusters {
		if c.Name == clusterName {
			cluster = c
			break
		}
	}

	if len(cluster.Name) == 0 {
		return fmt.Errorf("failed to get configuration for cluster: %s", clusterName)
	}

	err = login.New(cluster.Name, appCfg.User.Username).Login()
	if err != nil {
		return err
	}

	return nil
}
