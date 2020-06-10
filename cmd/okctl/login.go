package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"

	"github.com/oslokommune/okctl/pkg/login"
)

func buildLoginCommand(o *okctl.Okctl) *cobra.Command {
	var selectedCluster string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to AWS (deprecated)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Login(o, selectedCluster)
		},
	}

	cmd.PersistentFlags().StringVar(&selectedCluster, "cluster", "", "The target cluster")

	return cmd
}

func Login(o *okctl.Okctl, clusterName string) error {
	var cluster repository.Cluster

	for _, c := range o.RepoData.Clusters {
		if c.Name == clusterName {
			cluster = c
			break
		}
	}

	if len(cluster.Name) == 0 {
		return fmt.Errorf("failed to get configuration for cluster: %s", clusterName)
	}

	_, err := login.New(cluster.Name, o.AppData.User.Username).Login()
	if err != nil {
		return err
	}

	return nil
}
