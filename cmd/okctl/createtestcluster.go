package main

import (
	"fmt"
	"path"

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
	"github.com/oslokommune/okctl/pkg/config"
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

const testEndMsg = `Congratulations, your %s is now up and running.
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
`

// CreateTestClusterOpts contains all the required inputs
type CreateTestClusterOpts struct {
	Environment    string
	AWSAccountID   string
	RepositoryName string
	Region         string
	ClusterName    string
	Cidr           string
}

// Validate the inputs
func (o *CreateTestClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Cidr, validation.Required),
		validation.Field(&o.RepositoryName, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
	)
}

// nolint: funlen gocognit
func buildCreateTestClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &CreateTestClusterOpts{}

	cmd := &cobra.Command{
		Use:   "testcluster ENV AWS_ACCOUNT_ID",
		Short: "Create a lightweight cluster for testing and experimentation",
		Long: `This will create a lightweight cluster for testing and experimentation
that consumes a lot less resources, and it will not be tightly integrated
with Github or other production services.
`,
		Args: cobra.ExactArgs(createTestClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			if o.NoInput {
				return fmt.Errorf("we only support cluster creation with user input")
			}

			environment := args[0]
			awsAccountID := args[1]

			err := o.InitialiseWithEnvAndAWSAccountID(environment, awsAccountID)
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
				path.Join(userDir, config.DefaultLogDir, config.DefaultLogName),
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

			zones, err := route53.New(o.CloudProvider).PublicHostedZones()
			if err != nil {
				return formatErr(err)
			}

			var hostedZone *client.HostedZone
			if len(zones) > 0 {
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
			} else {
				hostedZone, err = services.Domain.CreatePrimaryHostedZone(o.Ctx, client.CreatePrimaryHostedZoneOpts{
					ID: id,
				})
				if err != nil {
					return formatErr(err)
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
				VpcID:             vpc.VpcID,
				VpcPrivateSubnets: vpc.PrivateSubnets,
				VpcPublicSubnets:  vpc.PublicSubnets,
				Minimal:           true,
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

			_, err = fmt.Fprintf(o.Err, testEndMsg,
				aurora.Green("kubernetes cluster"),
				exports,
				opts.Environment,
				aurora.Green("kubectl"),
				k.BinaryPath,
				k.BinaryPath,
				aurora.Green("aws-iam-authenticator"),
				a.BinaryPath,
				aurora.Green("1.17"),
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

	return cmd
}
