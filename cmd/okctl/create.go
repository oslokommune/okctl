package main

import (
	"fmt"
	"path"
	"regexp"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/domain"
	"github.com/oslokommune/okctl/pkg/git"

	"github.com/oslokommune/okctl/pkg/servicequota"

	"github.com/oslokommune/okctl/pkg/ask"
	stateSaver "github.com/oslokommune/okctl/pkg/client/core/state"

	"github.com/oslokommune/okctl/pkg/route53"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/github"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	defaultCidr = "192.168.0.0/20"
)

func buildCreateCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create commands",
		Long: `Create various kinds of resources.

Pre-requisites:
okctl creates various configuration files, and assumes that it is
running in a git repository. Initialize or clone a git repository
before running any of these commands.`,
	}

	cmd.AddCommand(buildCreateClusterCommand(o))
	cmd.AddCommand(buildCreateTestClusterCommand(o))

	return cmd
}

// createClusterOpts contains all the required inputs
type createClusterOpts struct {
	Environment    string
	AWSAccountID   string
	RepositoryName string
	Region         string
	ClusterName    string
	Cidr           string
	Organisation   string
}

// Validate the inputs
func (o *createClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z]{3,64}$")).Error("must consist of 3-64 characters (a-z, A-Z)")),
		validation.Field(&o.AWSAccountID, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{12}$")).Error("must consist of 12 digits")),
		validation.Field(&o.Cidr, validation.Required),
		validation.Field(&o.RepositoryName, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.Organisation, validation.Required),
	)
}

const startMsg = `We will now start to create your cluster, remember:

	- This can take upwards of %s to complete
	- Requires %s during the process

Requirements:

	- Be a member of the %s github organisation
	- Know the name of your %s on github
	- Have setup an %s on github
	- The infrastructure as code repository must be %s

If you are uncertain about any of these things, please
go to our slack channel and ask for help:

	- %s

You can tail the logs to get more output:

$ tail -f %s

`

const nsMsg = `
We could not detect any nameservers for your domain:

%s

We cannot continue with setting up the rest of the
cluster at this point, until we have delegated
the subdomain from the root account to your hosted
zone.

Ask in the %s slack channel for an update.

Once the hosted zone has been delegated you can
continue creating your cluster by rerunning the
same command:

$ okctl create cluster %s %s

`

// nolint: funlen gocyclo
func buildCreateClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &createClusterOpts{}

	cmd := &cobra.Command{
		Use:   "cluster ENV AWS_ACCOUNT_ID",
		Short: "Create a cluster",
		Long: `Fetch all tasks required to get an EKS cluster up and running on AWS.
This includes creating an EKS compatible VPC with private, public
and database subnets.`,
		// Args: cobra.ExactArgs(createClusterArgs),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			message := `We have made changes to the create cluster process and renamed the command.

To create a cluster, use %s instead.

For usage information, run %s
or see documentation on %s

`

			_, err := fmt.Fprintf(o.Err, message,
				aurora.Green("okctl apply cluster"),
				aurora.Green("okctl apply cluster --help"),
				aurora.Bold("https://okctl.io/usage/declarativecluster"),
			)
			if err != nil {
				return err
			}

			return errors.New("command removed")
		},
		PreRunE: func(_ *cobra.Command, args []string) error {
			environment := args[0]
			awsAccountID := args[1]

			err := validation.Validate(
				&environment,
				validation.Required,
				validation.Match(regexp.MustCompile("^[a-zA-Z]{3,64}$")).Error("the environment must consist of 3-64 characters (a-z, A-Z)"),
			)
			if err != nil {
				return err
			}

			err = validation.Validate(
				&awsAccountID,
				validation.Required,
				validation.Match(regexp.MustCompile("^[0-9]{12}$")).Error("the AWS Account ID must consist of 12 digits"),
			)
			if err != nil {
				return err
			}

			err = o.InitialiseWithEnvAndAWSAccountID(environment, awsAccountID)
			if err != nil {
				return err
			}

			meta := o.RepoStateWithEnv.GetMetadata()
			clusterName := o.RepoStateWithEnv.GetClusterName()

			opts.Environment = environment
			opts.AWSAccountID = awsAccountID
			opts.ClusterName = clusterName
			opts.RepositoryName = meta.Name
			opts.Region = meta.Region

			err = opts.Validate()
			if err != nil {
				return errors.E(err, "failed to validate create cluster options", errors.Invalid)
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			userDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			formatErr := o.ErrorFormatter(fmt.Sprintf("create cluster %s %s", opts.Environment, opts.AWSAccountID), userDir)

			repoDir, err := o.GetRepoDir()
			if err != nil {
				return formatErr(err)
			}

			id := api.ID{
				Region:       opts.Region,
				AWSAccountID: opts.AWSAccountID,
				Environment:  opts.Environment,
				Repository:   opts.RepositoryName,
				ClusterName:  opts.ClusterName,
			}

			_, err = fmt.Fprintf(o.Err, startMsg,
				aurora.Green("45 minutes"),
				aurora.Green("user input"),
				aurora.Green(opts.Organisation),
				aurora.Green("team"),
				aurora.Green("infrastructure as code repository"),
				aurora.Green("private"),
				aurora.Bold("#kjøremiljø-support"),
				path.Join(userDir, constant.DefaultLogDir, constant.DefaultLogName),
			)
			if err != nil {
				return formatErr(err)
			}

			// The cloud formation stack is created atomically, and the EIP and IGW
			// are created as part of this stack, therefore this check is sufficient
			// for all of these checks.
			vpcProvisioned := len(o.RepoStateWithEnv.GetVPC().VpcID) > 0

			err = servicequota.CheckQuotas(
				servicequota.NewVpcCheck(vpcProvisioned, constant.DefaultRequiredVpcs, o.CloudProvider),
				servicequota.NewEipCheck(vpcProvisioned, constant.DefaultRequiredEpis, o.CloudProvider),
				servicequota.NewIgwCheck(vpcProvisioned, constant.DefaultRequiredIgws, o.CloudProvider),
			)
			if err != nil {
				return formatErr(err)
			}

			ready := false
			prompt := &survey.Confirm{
				Message: "Are you ready to start?",
			}

			err = survey.AskOne(prompt, &ready)
			if err != nil {
				return formatErr(err)
			}

			if !ready {
				_, err = fmt.Fprintf(o.Err, "user wasn't ready to continue, aborting.")
				if err != nil {
					return formatErr(err)
				}

				return nil
			}

			spin, err := spinner.New("creating", o.Err)
			if err != nil {
				return formatErr(err)
			}

			services, err := o.ClientServices(spin)
			if err != nil {
				return err
			}

			var hostedZone *client.HostedZone

			hostedZone, err = services.Domain.GetPrimaryHostedZone(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			// Move this somewhere else
			if hostedZone == nil {
				zones, err := route53.New(o.CloudProvider).PublicHostedZones()
				if err != nil {
					return formatErr(err)
				}

				reuseHostedZone := func() error {
					zone, err := ask.New().SelectHostedZone(zones)
					if err != nil {
						return formatErr(err)
					}

					hostedZone = &client.HostedZone{
						IsDelegated: true,
						Primary:     true,
						HostedZone: &api.HostedZone{
							ID:           id,
							Managed:      false,
							FQDN:         zone.FQDN,
							Domain:       zone.Domain,
							HostedZoneID: zone.ID,
						},
					}

					state := stateSaver.NewDomainState(o.RepoStateWithEnv)
					_, err = state.SaveHostedZone(hostedZone)
					if err != nil {
						return formatErr(err)
					}

					return nil
				}

				doReuse := false

				if len(zones) > 0 {
					prompt = &survey.Confirm{
						Message: "We found existing hosted zones, do you want to reuse one?",
						Help:    "If you reuse an existing one, you won't have to wait for the zone to be delegated, and we will not remove it afterwards.",
					}

					err = survey.AskOne(prompt, &doReuse)
					if err != nil {
						return formatErr(err)
					}

					if doReuse {
						err = reuseHostedZone()
						if err != nil {
							return err
						}
					}
				}

				if !doReuse {
					hostedZone, err = services.Domain.CreatePrimaryHostedZone(o.Ctx, client.CreatePrimaryHostedZoneOpts{
						ID: id,
					})
					if err != nil {
						return formatErr(err)
					}
				}
			}

			pool, err := services.IdentityManager.CreateIdentityPool(o.Ctx, api.CreateIdentityPoolOpts{
				ID:           id,
				AuthDomain:   fmt.Sprintf("auth.%s", hostedZone.HostedZone.Domain),
				AuthFQDN:     fmt.Sprintf("auth.%s", hostedZone.HostedZone.FQDN),
				HostedZoneID: hostedZone.HostedZone.HostedZoneID,
			})
			if err != nil {
				return formatErr(err)
			}

			vpc, err := services.Vpc.CreateVpc(o.Ctx, api.CreateVpcOpts{
				ID:   id,
				Cidr: opts.Cidr,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.Cluster.CreateCluster(o.Ctx, api.ClusterCreateOpts{
				ID:                id,
				Cidr:              vpc.Cidr,
				Version:           constant.DefaultEKSKubernetesVersion,
				VpcID:             vpc.VpcID,
				VpcPrivateSubnets: vpc.PrivateSubnets,
				VpcPublicSubnets:  vpc.PublicSubnets,
				Minimal:           false,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.Blockstorage.CreateBlockstorage(o.Ctx, client.CreateBlockstorageOpts{
				ID: id,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.Autoscaler.CreateAutoscaler(o.Ctx, client.CreateAutoscalerOpts{
				ID: id,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.ExternalSecrets.CreateExternalSecrets(o.Ctx, client.CreateExternalSecretsOpts{
				ID: id,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.AWSLoadBalancerControllerService.CreateAWSLoadBalancerController(o.Ctx, client.CreateAWSLoadBalancerControllerOpts{
				ID:    id,
				VPCID: vpc.VpcID,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.ExternalDNS.CreateExternalDNS(o.Ctx, client.CreateExternalDNSOpts{
				ID:           id,
				HostedZoneID: hostedZone.HostedZone.HostedZoneID,
				Domain:       hostedZone.HostedZone.Domain,
			})
			if err != nil {
				return formatErr(err)
			}

			repo, err := git.GithubRepoFullName(opts.Organisation, repoDir)
			if err != nil {
				return err
			}

			githubRepo, err := services.Github.ReadyGithubInfrastructureRepository(o.Ctx, client.ReadyGithubInfrastructureRepositoryOpts{
				ID:           id,
				Organisation: opts.Organisation,
				Repository:   repo,
			})
			if err != nil {
				return err
			}

			gh := o.RepoStateWithEnv.GetGithub()
			gh.Organisation = opts.Organisation

			_, err = o.RepoStateWithEnv.SaveGithub(gh)
			if err != nil {
				return formatErr(err)
			}

			err = domain.ShouldHaveNameServers(hostedZone.HostedZone.Domain)
			if err != nil {
				_, err := fmt.Fprintf(o.Err, nsMsg,
					aurora.Green(hostedZone.HostedZone.Domain),
					aurora.Green("#kjøremiljø-support"),
					opts.Environment,
					opts.AWSAccountID,
				)
				if err != nil {
					return err
				}

				return err
			}

			argoCD, err := services.ArgoCD.CreateArgoCD(o.Ctx, client.CreateArgoCDOpts{
				ID:                 id,
				Domain:             hostedZone.HostedZone.Domain,
				FQDN:               hostedZone.HostedZone.FQDN,
				HostedZoneID:       hostedZone.HostedZone.HostedZoneID,
				GithubOrganisation: opts.Organisation,
				Repository:         githubRepo,
				UserPoolID:         pool.UserPoolID,
				AuthDomain:         pool.AuthDomain,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.Monitoring.CreateKubePromStack(o.Ctx, client.CreateKubePromStackOpts{
				ID:           id,
				Domain:       hostedZone.HostedZone.Domain,
				HostedZoneID: hostedZone.HostedZone.HostedZoneID,
				AuthDomain:   pool.AuthDomain,
				UserPoolID:   pool.UserPoolID,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.Monitoring.CreateLoki(o.Ctx, client.CreateLokiOpts{
				ID: id,
			})
			if err != nil {
				return formatErr(err)
			}

			_, err = services.Monitoring.CreatePromtail(o.Ctx, client.CreatePromtailOpts{
				ID: id,
			})
			if err != nil {
				return formatErr(err)
			}

			kubeConfig := path.Join(userDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterKubeConfig)
			awsConfig := path.Join(userDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterAwsConfig)
			awsCredentials := path.Join(userDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterAwsCredentials)

			exports := fmt.Sprintf(
				"export AWS_CONFIG_FILE=%s\nexport AWS_SHARED_CREDENTIALS_FILE=%s\nexport AWS_PROFILE=default\nexport KUBECONFIG=%s\n",
				awsConfig,
				awsCredentials,
				kubeConfig,
			)

			k, err := o.BinariesProvider.Kubectl(kubectl.Version)
			if err != nil {
				return formatErr(err)
			}

			a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
			if err != nil {
				return formatErr(err)
			}

			msg := commands.CreateClusterMsgOpts{
				KubernetesCluster:       aurora.Green("kubernetes cluster").String(),
				Exports:                 exports,
				Environment:             opts.Environment,
				VenvCmd:                 aurora.Green("okctl venv").String(),
				KubectlCmd:              aurora.Green("kubectl").String(),
				AwsIamAuthenticatorCmd:  aurora.Green("aws-iam-authenticator").String(),
				KubectlPath:             k.BinaryPath,
				AwsIamAuthenticatorPath: a.BinaryPath,
				K8sClusterVersion:       aurora.Green(constant.DefaultEKSKubernetesVersion).String(),
				ArgoCD:                  aurora.Green("ArgoCD").String(),
				ArgoCDURL:               aurora.Green(argoCD.ArgoURL).String(),
			}
			txt, err := commands.GoTemplateToString(commands.CreateClusterEndMsg, msg)
			_, err = fmt.Print(o.Err, txt)
			if err != nil {
				return formatErr(err)
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.Cidr, "cidr", "c", defaultCidr,
		"CIDR block the AWS VPC and subnets are created within")
	f.StringVarP(&opts.Organisation, "github-organisation", "o", github.DefaultOrg,
		"The Github organisation where we will look for your team and repository")

	return cmd
}
