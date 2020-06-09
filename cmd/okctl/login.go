package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/stage"
	"github.com/spf13/cobra"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/login"
	"github.com/oslokommune/okctl/pkg/storage"
)

func buildLoginCommand(appCfg *config.AppConfig, repoCfg *config.RepoConfig) *cobra.Command {
	var selectedCluster string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to AWS (deprecated)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Login(appCfg, repoCfg, selectedCluster)
		},
	}

	cmd.PersistentFlags().StringVar(&selectedCluster, "cluster", "", "The target cluster")

	return cmd
}

func Login(appCfg *config.AppConfig, repoCfg *config.RepoConfig, clusterName string) error {
	store := storage.NewFileSystemStorage(appCfg.BaseDir)

	stagers, err := stage.FromConfig(appCfg.Binaries, appCfg.Host, store)
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

	_, err = login.New(cluster.Name, appCfg.User.Username).Login()
	if err != nil {
		return err
	}

	return nil
}
