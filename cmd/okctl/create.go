package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"github.com/oslokommune/okctl/pkg/domain"

	"github.com/oslokommune/okctl/pkg/spinner"

	stateSaver "github.com/oslokommune/okctl/pkg/client/core/state"

	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"

	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/client/core/report/console"

	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/spf13/afero"

	"github.com/theckman/yacspin"

	"github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client/core/api/rest"
	"github.com/oslokommune/okctl/pkg/client/core/store/filesystem"
	"github.com/oslokommune/okctl/pkg/config"

	"github.com/oslokommune/okctl/pkg/git"
	"github.com/oslokommune/okctl/pkg/github"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/ask"
	"github.com/oslokommune/okctl/pkg/client"
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
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
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

const errMsg = `
Setup of cluster was aborted, because of:

%s

If you want to see more information about this error
you can %s of the command by running:

$ OKCTL_DEBUG=true okctl create cluster %s %s

These commands are %s and can be run as
many times as you like.

You can also inspect the %s, which contain more
information already, at:

%s

It is also possible to log onto your AWS account
and look for failures on the cloud formation stacks,
or inspect the pods using kubectl.

If you need any help to debug this failure, please
go to our slack channel and ask for help:

	- %s

Have the command you ran, logs, etc., ready.
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

			prettyErr := func(err error) error {
				return fmt.Errorf(errMsg,
					err,
					aurora.Blue("enable debugging"),
					opts.Environment,
					opts.AWSAccountID,
					aurora.Blue("idempotent"),
					aurora.Blue("logs"),
					path.Join(userDir, config.DefaultLogDir, config.DefaultLogName),
					aurora.Bold("#kjøremiljø-support"),
				)
			}

			repoDir, err := o.GetRepoDir()
			if err != nil {
				return prettyErr(err)
			}

			outputDir, err := o.GetRepoOutputDir(opts.Environment)
			if err != nil {
				return prettyErr(err)
			}

			c := rest.New(o.Debug, ioutil.Discard, o.ServerURL)

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
				aurora.Blue(opts.Organisation),
				aurora.Blue("team"),
				aurora.Blue("infrastructure as code repository"),
				aurora.Blue("private"),
				aurora.Bold("#kjøremiljø-support"),
				path.Join(userDir, config.DefaultLogDir, config.DefaultLogName),
			)
			if err != nil {
				return prettyErr(err)
			}

			ready := false
			prompt := &survey.Confirm{
				Message: "Are you ready to start?",
			}

			err = survey.AskOne(prompt, &ready)
			if err != nil {
				return prettyErr(err)
			}

			if !ready {
				_, err = fmt.Fprintf(o.Err, "user wasn't ready to continue, aborting.")
				if err != nil {
					return prettyErr(err)
				}

				return nil
			}

			spin, err := spinner.New("creating", o.Err)
			if err != nil {
				return prettyErr(err)
			}

			common := &common{
				Client:    c,
				Ctx:       o.Ctx,
				Fs:        o.FileSystem,
				ID:        id,
				OutputDir: outputDir,
				Out:       o.Err,
				Spinner:   spin,
				State:     o.RepoStateWithEnv,
			}

			ghClient, err := github.New(o.Ctx, o.CredentialsProvider.Github())
			if err != nil {
				return err
			}

			hostedZone, err := createPrimaryHostedZone(&hostedZone{
				common: common,
			})
			if err != nil {
				return prettyErr(err)
			}

			vpc, err := createVPC(&vpc{
				common: common,
				CIDR:   opts.Cidr,
			})
			if err != nil {
				return prettyErr(err)
			}

			_, err = createCluster(&cluster{
				common: common,
				CIDR:   opts.Cidr,
				VPC:    vpc,
			})
			if err != nil {
				return prettyErr(err)
			}

			_, err = createExternalSecrets(&externalSecrets{
				common: common,
			})
			if err != nil {
				return prettyErr(err)
			}

			_, err = createAlbIngressController(&ingressController{
				common: common,
				VPC:    vpc,
			})
			if err != nil {
				return prettyErr(err)
			}

			_, err = createExternalDNS(&externalDNS{
				common:     common,
				HostedZone: hostedZone.HostedZone,
			})
			if err != nil {
				return prettyErr(err)
			}

			githubRepo, err := createGithubRepo(&githubClient{
				common:       common,
				Github:       ghClient,
				Organisation: opts.Organisation,
				RepoDir:      repoDir,
			})

			gh := o.RepoStateWithEnv.GetGithub()
			gh.Organisation = opts.Organisation

			_, err = o.RepoStateWithEnv.SaveGithub(gh)
			if err != nil {
				return prettyErr(err)
			}

			err = domain.ShouldHaveNameServers(hostedZone.HostedZone.Domain)
			if err != nil {
				_, err := fmt.Fprintf(o.Err, nsMsg,
					aurora.Blue(hostedZone.HostedZone.Domain),
					aurora.Blue("#kjøremiljø-support"),
					opts.Environment,
					opts.AWSAccountID,
				)
				if err != nil {
					return err
				}

				return err
			}

			argoCD, err := createArgoCD(&argocdSetup{
				common:       common,
				HostedZone:   hostedZone.HostedZone,
				GithubRepo:   githubRepo,
				Organisation: opts.Organisation,
				Github:       ghClient,
			})
			if err != nil {
				return prettyErr(err)
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
				return prettyErr(err)
			}

			a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
			if err != nil {
				return prettyErr(err)
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
				return prettyErr(err)
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

type common struct {
	Client    *rest.HTTPClient
	Ctx       context.Context
	Fs        *afero.Afero
	ID        api.ID
	OutputDir string
	Out       io.Writer
	Spinner   *yacspin.Spinner
	State     state.RepositoryStateWithEnv
}

type vpc struct {
	*common
	CIDR string
}

func createVPC(v *vpc) (*api.Vpc, error) {
	_ = v.Spinner.Start()
	exit := spinner.Timer("vpc", v.Spinner)

	vpcService := core.NewVPCService(
		rest.NewVPCAPI(v.Client),
		filesystem.NewVpcStore(
			filesystem.Paths{
				OutputFile:         config.DefaultVpcOutputs,
				CloudFormationFile: config.DefaultVpcCloudFormationTemplate,
				BaseDir:            path.Join(v.OutputDir, config.DefaultVpcBaseDir),
			},
			v.Fs,
		),
		console.NewVPCReport(v.Out, v.Spinner, exit),
		stateSaver.NewVpcState(v.State),
	)

	return vpcService.CreateVpc(v.Ctx, api.CreateVpcOpts{
		ID:   v.ID,
		Cidr: v.CIDR,
	})
}

type cluster struct {
	*common
	CIDR string
	VPC  *api.Vpc
}

func createCluster(c *cluster) (*api.Cluster, error) {
	_ = c.Spinner.Start()
	exit := spinner.Timer("cluster", c.Spinner)

	clusterService := core.NewClusterService(
		rest.NewClusterAPI(c.Client),
		filesystem.NewClusterStore(
			filesystem.Paths{
				ConfigFile: config.DefaultClusterConfig,
				BaseDir:    path.Join(c.OutputDir, config.DefaultClusterBaseDir),
			},
			c.Fs,
		),
		console.NewClusterReport(c.Out, exit, c.Spinner),
		stateSaver.NewClusterState(c.State),
	)

	return clusterService.CreateCluster(c.Ctx, api.ClusterCreateOpts{
		ID:                c.ID,
		Cidr:              c.CIDR,
		VpcID:             c.VPC.VpcID,
		VpcPrivateSubnets: c.VPC.PrivateSubnets,
		VpcPublicSubnets:  c.VPC.PublicSubnets,
	})
}

type externalSecrets struct {
	*common
}

func createExternalSecrets(e *externalSecrets) (*client.ExternalSecrets, error) {
	_ = e.Spinner.Start()
	exit := spinner.Timer("external-secrets", e.Spinner)

	externalSecretsService := core.NewExternalSecretsService(
		rest.NewExternalSecretsAPI(e.Client),
		filesystem.NewExternalSecretsStore(
			filesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(e.OutputDir, config.DefaultExternalSecretsBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(e.OutputDir, config.DefaultExternalSecretsBaseDir),
			},
			filesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(e.OutputDir, config.DefaultExternalSecretsBaseDir),
			},
			e.Fs,
		),
		console.NewExternalSecretsReport(e.Out, exit, e.Spinner),
	)

	return externalSecretsService.CreateExternalSecrets(e.Ctx, client.CreateExternalSecretsOpts{
		ID: e.ID,
	})
}

type ingressController struct {
	*common
	VPC *api.Vpc
}

func createAlbIngressController(c *ingressController) (*client.ALBIngressController, error) {
	_ = c.Spinner.Start()
	exit := spinner.Timer("alb-ingress-controller", c.Spinner)

	albIngressControllerService := core.NewALBIngressControllerService(
		rest.NewALBIngressControllerAPI(c.Client),
		filesystem.NewALBIngressControllerStore(
			filesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(c.OutputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(c.OutputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			filesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(c.OutputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			c.Fs,
		),
		console.NewAlbIngressControllerReport(c.Out, exit, c.Spinner),
	)

	return albIngressControllerService.CreateALBIngressController(c.Ctx, client.CreateALBIngressControllerOpts{
		ID:    c.ID,
		VPCID: c.VPC.VpcID,
	})
}

type hostedZone struct {
	*common
}

func createPrimaryHostedZone(h *hostedZone) (*client.HostedZone, error) {
	_ = h.Spinner.Start()
	exit := spinner.Timer("primary-hosted-zone", h.Spinner)

	a := ask.New().WithSpinner(h.Spinner)

	domainService := core.NewDomainService(
		h.Out,
		a,
		rest.NewDomainAPI(h.Client),
		filesystem.NewDomainStore(
			filesystem.Paths{
				OutputFile:         config.DefaultDomainOutputsFile,
				CloudFormationFile: config.DefaultDomainCloudFormationTemplate,
				BaseDir:            path.Join(h.OutputDir, config.DefaultDomainBaseDir),
			},
			h.Fs,
		),
		console.NewDomainReport(h.Out, exit, h.Spinner),
		stateSaver.NewDomainState(h.State),
	)

	return domainService.CreatePrimaryHostedZone(h.Ctx, client.CreatePrimaryHostedZoneOpts{
		ID: h.ID,
	})
}

type externalDNS struct {
	*common
	HostedZone *api.HostedZone
}

func createExternalDNS(e *externalDNS) (*client.ExternalDNS, error) {
	_ = e.Spinner.Start()
	exit := spinner.Timer("external-dns", e.Spinner)

	externalDNSService := core.NewExternalDNSService(
		rest.NewExternalDNSAPI(e.Client),
		filesystem.NewExternalDNSStore(
			filesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(e.OutputDir, config.DefaultExternalDNSBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(e.OutputDir, config.DefaultExternalDNSBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultKubeOutputsFile,
				BaseDir:    path.Join(e.OutputDir, config.DefaultExternalDNSBaseDir),
			},
			e.Fs,
		),
		console.NewExternalDNSReport(e.Out, exit, e.Spinner),
	)

	return externalDNSService.CreateExternalDNS(e.Ctx, client.CreateExternalDNSOpts{
		ID:           e.ID,
		HostedZoneID: e.HostedZone.HostedZoneID,
		Domain:       e.HostedZone.Domain,
	})
}

type githubClient struct {
	*common
	Github       github.Githuber
	Organisation string
	RepoDir      string
}

func createGithubRepo(c *githubClient) (*client.GithubRepository, error) {
	repo, err := git.GithubRepoFullName(c.Organisation, c.RepoDir)
	if err != nil {
		return nil, err
	}

	_ = c.Spinner.Start()
	exit := spinner.Timer("github", c.Spinner)

	githubService := core.NewGithubService(
		rest.NewGithubAPI(
			c.Out,
			ask.New().WithSpinner(c.Spinner),
			rest.NewParameterAPI(c.Client),
			c.Github,
		),
		console.NewGithubReport(c.Out, exit, c.Spinner),
		stateSaver.NewGithubState(c.State),
	)

	githubRepo, err := githubService.ReadyGithubInfrastructureRepository(c.Ctx, client.ReadyGithubInfrastructureRepositoryOpts{
		ID:           c.ID,
		Organisation: c.Organisation,
		Repository:   repo,
	})

	return githubRepo, err
}

type argocdSetup struct {
	*common
	HostedZone   *api.HostedZone
	Github       github.Githuber
	GithubRepo   *client.GithubRepository
	Organisation string
}

// nolint: funlen
func createArgoCD(a *argocdSetup) (*client.ArgoCD, error) {
	argoBaseDir := path.Join(a.OutputDir, config.DefaultArgoCDBaseDir)

	_ = a.Spinner.Start()
	exit := spinner.Timer("argocd", a.Spinner)

	githubService := core.NewGithubService(
		rest.NewGithubAPI(
			a.Out,
			ask.New().WithSpinner(a.Spinner),
			rest.NewParameterAPI(a.Client),
			a.Github,
		),
		console.NewGithubReport(a.Out, nil, a.Spinner),
		stateSaver.NewGithubState(a.State),
	)

	certService := core.NewCertificateService(
		rest.NewCertificateAPI(a.Client),
		filesystem.NewCertificateStore(
			filesystem.Paths{
				OutputFile:         config.DefaultCertificateOutputsFile,
				CloudFormationFile: config.DefaultCertificateCloudFormationTemplate,
				BaseDir:            path.Join(argoBaseDir, config.DefaultCertificateBaseDir),
			},
			a.Fs,
		),
		stateSaver.NewCertificateState(a.State),
		console.NewCertificateReport(a.Out, nil, a.Spinner),
	)

	manifestService := core.NewManifestService(
		rest.NewManifestAPI(a.Client),
		filesystem.NewManifestStore(
			filesystem.Paths{
				OutputFile: config.DefaultKubeOutputsFile,
				BaseDir:    path.Join(argoBaseDir, config.DefaultExternalSecretsBaseDir),
			},
			a.Fs,
		),
		console.NewManifestReport(a.Out, nil, a.Spinner),
	)

	paramService := core.NewParameterService(
		rest.NewParameterAPI(a.Client),
		filesystem.NewParameterStore(
			filesystem.Paths{
				OutputFile: config.DefaultParameterOutputsFile,
				BaseDir:    path.Join(argoBaseDir, config.DefaultParameterBaseDir),
			},
			a.Fs,
		),
		console.NewParameterReport(a.Out, nil, a.Spinner),
	)

	argoService := core.NewArgoCDService(
		githubService,
		certService,
		manifestService,
		paramService,
		rest.NewArgoCDAPI(a.Client),
		filesystem.NewArgoCDStore(
			filesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(argoBaseDir, config.DefaultHelmBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultArgoOutputsFile,
				BaseDir:    argoBaseDir,
			},
			a.Fs,
		),
		console.NewArgoCDReport(a.Out, exit, a.Spinner),
		stateSaver.NewArgoCDState(a.State),
	)

	return argoService.CreateArgoCD(a.Ctx, client.CreateArgoCDOpts{
		ID:                 a.ID,
		Domain:             a.HostedZone.Domain,
		FQDN:               a.HostedZone.FQDN,
		HostedZoneID:       a.HostedZone.HostedZoneID,
		GithubOrganisation: a.Organisation,
		Repository:         a.GithubRepo,
	})
}
