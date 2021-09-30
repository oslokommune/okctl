package main

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildVersionCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: VersionShortDescription,
		Long:  VersionLongDescription,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			versioner := version.New(context.Background())

			ver, err := versioner.String()
			if err != nil {
				return fmt.Errorf("getting version: %w", err)
			}

			_, err = fmt.Fprintln(o.Out, ver)

			return err
		},
	}

	return cmd
}
