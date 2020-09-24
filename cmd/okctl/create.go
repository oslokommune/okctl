package main

import (
	"fmt"
	"path"
	"regexp"

	"github.com/oslokommune/okctl/pkg/servicequota"

	"github.com/oslokommune/okctl/pkg/ask"
	stateSaver "github.com/oslokommune/okctl/pkg/client/core/state"

	"github.com/oslokommune/okctl/pkg/git"

	"github.com/oslokommune/okctl/pkg/route53"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/domain"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"

	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/oslokommune/okctl/pkg/github"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	createClusterArgs = 2
	defaultCidr       = "192.168.0.0/20"
)

func buildCreateCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create commands",
	}

	cmd.AddCommand(buildCreateClusterCommand(o))
	cmd.AddCommand(buildCreateTestClusterCommand(o))

	return cmd
}

// CreateClusterOpts contains all the required inputs
type CreateClusterOpts struct {
	Environment    string
	AWSAccountID   string
	RepositoryName string
	Region         string
	ClusterName    string
	Cidr           string
	Organisation   string
}

// Validate the inputs
func (o *CreateClusterOpts) Validate() error {
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

const endMsg = `Congratulations, your %s is now up and running.
To get started with some basic interactions, you can paste the
following exports into a terminal:

%s

You can retrieve these credentials at any point by issuing the
command below, from within this repository:

$ okctl show credentials %s

Now you can use %s to list nodes, pods, etc. Try out some commands:

$ %s get pods --all-namespaces
$ %s get nodes

This also requires %s, which you can add to your PATH from here:

%s

Optionally, install kubectl and aws-iam-authenticator to your 
system from:

- https://kubernetes.io/docs/tasks/tools/install-kubectl/
- https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html

The installed version of kubectl needs to be within 2 versions of the
kubernetes cluster version, which is: %s.

We have also setup %s for continuous deployment, you can access
the UI at this URL by logging in with Github:

%s

It might take 5-10 minutes for the ArgoCD ALB to come up.
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
	opts := &CreateClusterOpts{}

	cmd := &cobra.Command{
		Use:   "cluster ENV AWS_ACCOUNT_ID",
		Short: "Create a cluster",
		Long: `Fetch all tasks required to get an EKS cluster up and running on AWS.
This includes creating an EKS compatible VPC with private, public
and database subnets.`,
		Args: cobra.ExactArgs(createClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			if o.NoInput {
				return fmt.Errorf("we currently don't support no user input")
			}

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
				path.Join(userDir, config.DefaultLogDir, config.DefaultLogName),
			)
			if err != nil {
				return formatErr(err)
			}

			var checkers []servicequota.Checker
			checkers = append(checkers,
				servicequota.NewVpcCheck(o.Err, o.CloudProvider, config.DefaultRequiredVpcs),
				servicequota.NewEipCheck(o.Err, o.CloudProvider, config.DefaultRequiredEpis),
				servicequota.NewIgwCheck(o.Err, o.CloudProvider, config.DefaultRequiredIgws))
			servicequota.CheckQuotas(checkers)

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
				VpcID:             vpc.VpcID,
				VpcPrivateSubnets: vpc.PrivateSubnets,
				VpcPublicSubnets:  vpc.PublicSubnets,
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

			_, err = services.ALBIngressController.CreateALBIngressController(o.Ctx, client.CreateALBIngressControllerOpts{
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
			})
			if err != nil {
				return formatErr(err)
			}

			kubeConfig := path.Join(userDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterKubeConfig)
			awsConfig := path.Join(userDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsConfig)
			awsCredentials := path.Join(userDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsCredentials)

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

			_, err = fmt.Fprintf(o.Err, endMsg,
				aurora.Green("kubernetes cluster"),
				exports,
				opts.Environment,
				aurora.Green("kubectl"),
				k.BinaryPath,
				k.BinaryPath,
				aurora.Green("aws-iam-authenticator"),
				a.BinaryPath,
				aurora.Green("1.17"),
				aurora.Green("ArgoCD"),
				argoCD.ArgoURL,
			)
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
