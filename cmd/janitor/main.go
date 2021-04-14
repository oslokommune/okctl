package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "janitor",
		Short: "Simple tool for doing different types of maintenance",
	}

	cmd.AddCommand(buildHostedZoneCommand())

	return cmd
}
