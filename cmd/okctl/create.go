package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/manager"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/login"
	"github.com/oslokommune/okctl/pkg/stage"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/spf13/cobra"
)

const AWSAccountDigits = 12

func buildCreateCommand(appCfg *config.AppConfig, repoCfg *config.RepoConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create commands",
	}

	cmd.AddCommand(buildCreateClusterCommand(appCfg, repoCfg))

	return cmd
}

func buildCreateClusterCommand(appCfg *config.AppConfig, repoCfg *config.RepoConfig) *cobra.Command {
	var opts CreateClusterOpts

	cmd := &cobra.Command{
		Use:   "cluster [env] [aws account]",
		Short: "Create a cluster",
		Args:  cobra.ExactArgs(2), // nolint
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Env = args[0]
			account, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			opts.AWSAccount = account
			return opts.Validate()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return CreateCluster(appCfg, repoCfg, opts)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.Cidr, "cidr", "c", "192.168.0.0/20",
		"The CIDR used to create the cluster with")
	f.StringVarP(&opts.Region, "region", "r", "eu-west-1",
		"The AWS region we should create the cluster in")

	return cmd
}

type CreateClusterOpts struct {
	Env        string
	Cidr       string
	Region     string
	AWSAccount int
}

func (o *CreateClusterOpts) Validate() error {
	number := o.AWSAccount

	count := 0

	for number != 0 {
		number /= 10
		count += 1
	}

	if count != AWSAccountDigits {
		return fmt.Errorf("a valid AWS account consists of 12 digits, got: %d", o.AWSAccount)
	}

	if len(o.Env) == 0 {
		return errors.New("the name given to the cluster cannot be empty")
	}

	return nil
}

func CreateCluster(appCfg *config.AppConfig, repoCfg *config.RepoConfig, opts CreateClusterOpts) error {
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

	for _, c := range repoCfg.Clusters {
		if c.Name == opts.Env && c.AWS.Account == opts.AWSAccount {
			return fmt.Errorf("cluster: %s, already exists", opts.Env)
		}
	}

	creds, err := login.New(fmt.Sprintf("%d", opts.AWSAccount), appCfg.User.Username).Login()
	if err != nil {
		return err
	}

	m := manager.New(vpc.New(repoCfg.Name, opts.Env, opts.Cidr, opts.Region), creds)

	err = m.Create(fmt.Sprintf("%s-%s-cluster-vpc", repoCfg.Name, opts.Env), 5)
	if err != nil {
		return err
	}

	// 4. Run eksctl

	return nil
}
