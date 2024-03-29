package main

import (
	"fmt"

	common "github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/cmd/okctl/handlers"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// requiredApplyApplicationArguments defines number of arguments the ApplyApplication command expects
const requiredDeleteApplicationArguments = 0

//nolint funlen
func buildDeleteApplicationCommand(o *okctl.Okctl) *cobra.Command {
	opts := &handlers.HandleApplicationOpts{
		Out:           o.Out,
		Err:           o.Err,
		Ctx:           o.Ctx,
		DelayFunction: common.DefaultDelayFunction,
		Purge:         true,
	}

	cmd := &cobra.Command{
		Use:   "application",
		Short: deleteApplicationShortDescription,
		Long:  deleteApplicationLongDescription,
		Args:  cobra.ExactArgs(requiredDeleteApplicationArguments),
		PreRunE: hooks.RunECombinator(
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionDeleteApplication),
			hooks.LoadClusterDeclaration(o, &opts.ClusterDeclarationPath),
			hooks.InitializeOkctl(o),
			hooks.AcquireStateLock(o),
			hooks.DownloadState(o, true),
			func(cmd *cobra.Command, args []string) (err error) {
				err = commands.ValidateBinaryEqualsClusterVersion(o)
				if err != nil {
					return err
				}

				opts.ClusterManifest = *o.Declaration

				opts.ApplicationManifest, err = commands.InferApplicationFromStdinOrFile(*o.Declaration, o.In, o.FileSystem, opts.File)
				if err != nil {
					return fmt.Errorf("inferring application from stdin or file: %w", err)
				}

				opts.State = o.StateHandlers(o.StateNodes())

				opts.Services, err = o.ClientServices(opts.State)
				if err != nil {
					return fmt.Errorf("preparing client services: %w", err)
				}

				return nil
			},
		),
		RunE: handlers.HandleApplication(opts),
		PostRunE: hooks.RunECombinator(
			hooks.UploadState(o),
			hooks.ClearLocalState(o),
			hooks.ReleaseStateLock(o),
			hooks.EmitEndCommandExecutionEvent(metrics.ActionDeleteApplication),
		),
	}
	addAuthenticationFlags(cmd)
	addClusterDeclarationPathFlag(cmd, &opts.ClusterDeclarationPath)

	cmd.Flags().StringVarP(
		&opts.File,
		"file",
		"f", "",
		"Specify the file path for the application to delete. Use \"-\" for stdin",
	)
	cmd.Flags().BoolVarP(&opts.Confirm, "confirm", "y", false, "confirm all choices")

	return cmd
}
