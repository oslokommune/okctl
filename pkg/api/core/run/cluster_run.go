// Package run implements the runnable layer
package run

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/oslokommune/okctl/pkg/kubeconfig"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
)

type clusterRun struct {
	awsCredentialsPath string
	awsConfigPath      string
	provider           binaries.Provider
	cloud              v1alpha1.CloudProvider
	debug              bool
	kubeConfigStore    api.KubeConfigStore
}

// CreateCluster invokes a CLI for performing create
// nolint: funlen
func (c *clusterRun) CreateCluster(opts api.ClusterCreateOpts) (*api.Cluster, error) {
	a, err := c.provider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return nil, fmt.Errorf("retrieving aws-iam-authenticator binary: %w", err)
	}

	k, err := c.provider.Kubectl(kubectl.Version)
	if err != nil {
		return nil, fmt.Errorf("retrieving kubectl binary: %w", err)
	}

	cli, err := c.provider.Eksctl(eksctl.Version)
	if err != nil {
		return nil, fmt.Errorf("retrieving eksctl binary: %w", err)
	}

	var cfg *v1alpha5.ClusterConfig
	if opts.Minimal {
		cfg, err = clusterconfig.NewMinimal(&clusterconfig.MinimalArgs{
			ClusterName:            opts.ID.ClusterName,
			PermissionsBoundaryARN: v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
			PrivateSubnets:         opts.VpcPrivateSubnets,
			PublicSubnets:          opts.VpcPublicSubnets,
			Region:                 opts.ID.Region,
			Version:                opts.Version,
			VpcCidr:                opts.Cidr,
			VpcID:                  opts.VpcID,
		})
	} else {
		cfg, err = clusterconfig.New(&clusterconfig.Args{
			ClusterName:            opts.ID.ClusterName,
			PermissionsBoundaryARN: v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
			PrivateSubnets:         opts.VpcPrivateSubnets,
			PublicSubnets:          opts.VpcPublicSubnets,
			Region:                 opts.ID.Region,
			Version:                opts.Version,
			VpcCidr:                opts.Cidr,
			VpcID:                  opts.VpcID,
		})
	}

	if err != nil {
		return nil, errors.E(err, "creating cluster config", errors.Internal)
	}

	cli.Debug(c.debug)
	cli.AddToPath(a.BinaryPath, k.BinaryPath)
	cli.AddToEnv(
		fmt.Sprintf("AWS_CONFIG_FILE=%s", c.awsConfigPath),
		fmt.Sprintf("AWS_SHARED_CREDENTIALS_FILE=%s", c.awsCredentialsPath),
		"AWS_PROFILE=default",
	)

	exists, err := cli.HasCluster(cfg)
	if err != nil {
		return nil, fmt.Errorf("if cluster exists: %w", err)
	}

	cluster := &api.Cluster{
		ID:     opts.ID,
		Config: cfg,
	}

	if exists {
		return cluster, nil
	}

	_, err = cli.CreateCluster(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating the cluster: %w", err)
	}

	kubeConf, err := kubeconfig.New(cfg, c.cloud).Get()
	if err != nil {
		return nil, err
	}

	err = c.kubeConfigStore.SaveKubeConfig(kubeConf)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// DeleteCluster invokes a CLI for performing delete
func (c *clusterRun) DeleteCluster(opts api.ClusterDeleteOpts) error {
	cli, err := c.provider.Eksctl(eksctl.Version)
	if err != nil {
		return fmt.Errorf("failed to retrieve eksctl binary: %w", err)
	}

	// We want to continue even if delete fargate profiles fails,
	// delete cluster should succeed if fp-default is not found
	_, err = cli.DeleteFargateProfiles(opts.ID.ClusterName, opts.FargateProfileName)
	if err != nil {
		fmt.Println(err)
	}

	_, err = cli.DeleteCluster(opts.ID.ClusterName)
	if err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}

	return nil
}

// NewClusterRun returns a executor for clusterRun
func NewClusterRun(
	debug bool,
	kubeConfigStore api.KubeConfigStore,
	awsCredentialsPath, awsConfigPath string,
	provider binaries.Provider,
	cloud v1alpha1.CloudProvider,
) api.ClusterRun {
	return &clusterRun{
		kubeConfigStore:    kubeConfigStore,
		debug:              debug,
		awsCredentialsPath: awsCredentialsPath,
		awsConfigPath:      awsConfigPath,
		provider:           provider,
		cloud:              cloud,
	}
}
