package main

import (
	"fmt"
	"github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

var (
	fullExample bool
)

func buildCreateApplicationCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use: "application",
		Short: "Create an application template",
		Long: "Scaffolds an application.yaml template which can be used to produce necessary Kubernetes and ArgoCD resources",
		RunE: func(_ *cobra.Command, _ []string) error {
			if fullExample {
				err := api.FetchFullExample(o.Out)
				if err != nil {
					return fmt.Errorf("failed fetching application.yaml example: %w", err)
				}
			} else {
				err := api.FetchMinimalExample(o.Out)
				if err != nil {
					return fmt.Errorf("failed fetching application.yaml example: %w", err)
				}
			}
			
			return nil
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
	
	cmd.Flags().BoolVarP(&fullExample, "full", "f", false, "Scaffold a full rather than a minimal example")
	
	return cmd
}
