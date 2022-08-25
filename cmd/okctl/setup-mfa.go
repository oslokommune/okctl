package main

import (
	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/metrics"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

type buildSetupMFAOpts struct {
	ClusterManifestPath string
}

func buildSetupMFA(o *okctl.Okctl) *cobra.Command {
	opts := buildSetupMFAOpts{}

	cmd := &cobra.Command{
		Use:     "setup-mfa",
		Example: "okctl setup-mfa olly@okctl.io",
		Short:   "Register MFA device for Cognito user",
		Args:    cobra.ExactArgs(1),
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionApplyCluster),
			hooks.LoadClusterDeclaration(o, &opts.ClusterManifestPath),
			hooks.InitializeOkctl(o),
			hooks.AcquireStateLock(o),
			hooks.DownloadState(o, true),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	addAuthenticationFlags(cmd)
	addClusterDeclarationPathFlag(cmd, &opts.ClusterManifestPath)

	return cmd
}
