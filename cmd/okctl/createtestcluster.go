package main

import (
	"fmt"
	"path"
	"regexp"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/servicequota"

	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"

	stateSaver "github.com/oslokommune/okctl/pkg/client/core/state"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/route53"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

const (
	createTestClusterArgs = 2
)

const testStartMsg = `We will now start to create your cluster, remember:

	- This can take upwards of %s to complete
	- Requires %s during the process

You can always ask for help in our slack channel:

	- %s

You can tail the logs to get more output:

$ tail -f %s

`

// createTestClusterOpts contains all the required inputs
type createTestClusterOpts struct {
	Environment    string
	AWSAccountID   string
	RepositoryName string
	Region         string
	ClusterName    string
	Cidr           string
}

// Validate the inputs
func (o *createTestClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z]{3,64}$")).Error("must consist of 3-64 characters (a-z, A-Z)")),
		validation.Field(&o.AWSAccountID, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{12}$")).Error("must consist of 12 digits")),
		validation.Field(&o.Cidr, validation.Required),
		validation.Field(&o.RepositoryName, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
	)
}

// nolint: funlen gocognit
func buildCreateTestClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &createTestClusterOpts{}

	cmd := &cobra.Command{
		Use:   "testcluster ENV AWS_ACCOUNT_ID",
		Short: "Create a lightweight cluster for testing and experimentation",
		Long: `This will create a lightweight cluster for testing and experimentation
that consumes a lot less resources, and it will not be tightly integrated
with Github or other production services.
`,
		Args: cobra.ExactArgs(createTestClusterArgs),
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

			formatErr := o.ErrorFormatter(fmt.Sprintf("create testcluster %s %s", opts.Environment, opts.AWSAccountID), userDir)

			id := api.ID{
				Region:       opts.Region,
				AWSAccountID: opts.AWSAccountID,
				Environment:  opts.Environment,
				Repository:   opts.RepositoryName,
				ClusterName:  opts.ClusterName,
			}

			_, err = fmt.Fprintf(o.Err, testStartMsg,
				aurora.Green("45 minutes"),
				aurora.Green("user input"),
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
				servicequota.NewVpcCheck(vpcProvisioned, constant.DefaultRequiredVpcsTestCluster, o.CloudProvider),
				servicequota.NewEipCheck(vpcProvisioned, constant.DefaultRequiredEpisTestCluster, o.CloudProvider),
				servicequota.NewIgwCheck(vpcProvisioned, constant.DefaultRequiredIgwsTestCluster, o.CloudProvider),
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

			vpc, err := services.Vpc.CreateVpc(o.Ctx, api.CreateVpcOpts{
				ID:      id,
				Cidr:    opts.Cidr,
				Minimal: true,
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
				Minimal:           true,
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
				K8sClusterVersion:       aurora.Green("1.17").String(),
			}
			txt, err := commands.GoTemplateToString(commands.CreateTestClusterEndMsg, msg)
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

	cmd.Hidden = true

	return cmd
}
