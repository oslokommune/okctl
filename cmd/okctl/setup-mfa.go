package main

import (
	"context"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/cognito"
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
			ctx := context.Background()

			userEmail := args[0]

			return cognito.RegisterMFADevice(cognito.RegisterMFADeviceOpts{
				Ctx:                    ctx,
				CognitoProvider:        o.CloudProvider.CognitoIdentityProvider(),
				ParameterStoreProvider: o.CloudProvider.SSM(),
				UserEmail:              userEmail,
				Cluster:                *o.Declaration,
			})
		},
	}

	addAuthenticationFlags(cmd)
	addClusterDeclarationPathFlag(cmd, &opts.ClusterManifestPath)

	return cmd
}
