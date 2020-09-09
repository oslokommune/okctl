package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"time"

	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"

	github2 "github.com/oslokommune/okctl/pkg/credentials/github"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/client/core/report/console"

	"github.com/oslokommune/okctl/pkg/config/repository"

	"github.com/spf13/afero"

	"github.com/hako/durafmt"

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
	Cidr           string
	RepositoryName string
	Region         string
	ClusterName    string
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
	- The infrastructure as code repository, must be %s

If you are uncertain about any of these things, please
go to our slack channel and ask for help:

	- %s

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

You can also install %s to your system from:
https://kubernetes.io/docs/tasks/tools/install-kubectl/

The installed version needs to be within 2 versions of the
kubernetes cluster version, which is: %s.

We have also setup %s for continuous deployment, you can access
the UI at this URL by logging in with Github:

%s

Happy, controlling!
`

const errMsg = `
Setup of cluster was aborted, because of:

%s

If you want to see more information about this error
you can %s of the command by running:

$ OKCTL_DEBUG=true okctl create cluster %s %s

These commands are idempotent and can be run as
many times as you like.

You can also inspect the logs, which contain more
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

			opts.Environment = args[0]
			opts.AWSAccountID = args[1]
			opts.RepositoryName = o.RepoData.Name
			opts.ClusterName = o.ClusterName(opts.Environment)
			opts.Region = o.Region()
			opts.Organisation = github.DefaultOrg

			err := opts.Validate()
			if err != nil {
				return errors.E(err, "failed to validate create cluster options", errors.Invalid)
			}

			return o.Initialise(opts.Environment, opts.AWSAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			appDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			prettyErr := func(err error) error {
				return fmt.Errorf(errMsg,
					err,
					aurora.Blue("enable debugging"),
					opts.Environment,
					opts.AWSAccountID,
					aurora.Bold("#kjøremiljø-support"),
					path.Join(appDir, config.DefaultLogDir, config.DefaultLogName),
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

			o.SetGithubOrganisationName(opts.Organisation, opts.Environment)
			err = o.WriteCurrentRepoData()
			if err != nil {
				return prettyErr(err)
			}

			_, err = fmt.Fprintf(o.Err, startMsg,
				aurora.Green("45 minutes"),
				aurora.Green("user input"),
				aurora.Blue(opts.Organisation),
				aurora.Blue("team"),
				aurora.Blue("infrastructure as code repository"),
				aurora.Blue("private"),
				aurora.Bold("#kjøremiljø-support"),
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

			cfg := yacspin.Config{
				Frequency:       100 * time.Millisecond,
				CharSet:         yacspin.CharSets[59],
				Suffix:          " creating",
				SuffixAutoColon: true,
				StopCharacter:   "✓",
				StopColors:      []string{"fgGreen"},
				Writer:          o.Out,
			}

			spinner, err := yacspin.New(cfg)
			if err != nil {
				return fmt.Errorf("failed to create spinner")
			}

			vpc, err := createVPC(&vpc{
				ctx:       o.Ctx,
				out:       o.Err,
				spinner:   spinner,
				cidr:      opts.Cidr,
				id:        id,
				c:         c,
				outputDir: outputDir,
				fs:        o.FileSystem,
			})
			if err != nil {
				spinner.StopFailMessage(err.Error())
				return prettyErr(err)
			}

			_, err = createCluster(&cluster{
				ctx:       o.Ctx,
				out:       o.Err,
				spinner:   spinner,
				vpc:       vpc,
				cidr:      opts.Cidr,
				id:        id,
				c:         c,
				repoData:  o.RepoData,
				outputDir: outputDir,
				repoDir:   repoDir,
				fs:        o.FileSystem,
			})
			if err != nil {
				return prettyErr(err)
			}

			_, err = createExternalSecrets(&externalSecrets{
				ctx:       o.Ctx,
				out:       o.Err,
				spinner:   spinner,
				c:         c,
				outputDir: outputDir,
				fs:        o.FileSystem,
				id:        id,
			})
			if err != nil {
				return prettyErr(err)
			}

			_, err = createAlbIngressController(&ingressController{
				ctx:       o.Ctx,
				out:       o.Err,
				spinner:   spinner,
				outputDir: outputDir,
				c:         c,
				fs:        o.FileSystem,
				id:        id,
				vpc:       vpc,
			})
			if err != nil {
				return prettyErr(err)
			}

			d, err := createPrimaryHostedZone(&hostedZone{
				ctx:       o.Ctx,
				out:       o.Err,
				spinner:   spinner,
				repoData:  o.RepoData,
				c:         c,
				outputDir: outputDir,
				repoDir:   repoDir,
				fs:        o.FileSystem,
				id:        id,
			})
			if err != nil {
				return prettyErr(err)
			}

			_, err = createExternalDNS(&externalDNS{
				ctx:       o.Ctx,
				out:       o.Err,
				spinner:   spinner,
				id:        id,
				outputDir: outputDir,
				c:         c,
				fs:        o.FileSystem,
				d:         d,
			})
			if err != nil {
				return prettyErr(err)
			}

			githubService, githubRepo, err := createGithubRepo(&githubClient{
				ctx:          o.Ctx,
				out:          o.Err,
				spinner:      spinner,
				repoDir:      repoDir,
				organisation: opts.Organisation,
				id:           id,
				repoData:     o.RepoData,
				c:            c,
				fs:           o.FileSystem,
				auth:         o.CredentialsProvider.Github(),
			})

			argoCD, err := createArgoCD(&argocdSetup{
				ctx:           o.Ctx,
				out:           o.Err,
				spinner:       spinner,
				id:            id,
				githubRepo:    githubRepo,
				githubService: githubService,
				d:             d,
				outputDir:     outputDir,
				repoDir:       repoDir,
				organisation:  opts.Organisation,
				fs:            o.FileSystem,
				c:             c,
				repoData:      o.RepoData,
			})
			if err != nil {
				return prettyErr(err)
			}

			kubeConfig := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterKubeConfig)
			awsConfig := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterAwsConfig)
			awsCredentials := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterAwsCredentials)

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

			_, err = fmt.Fprintf(o.Err, endMsg,
				aurora.Green("kubernetes cluster"),
				exports,
				opts.Environment,
				aurora.Green("kubectl"),
				k.BinaryPath,
				k.BinaryPath,
				aurora.Green("kubectl"),
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

	return cmd
}

func timer(spinner *yacspin.Spinner, component string) chan struct{} {
	exit := make(chan struct{})

	go func(ch chan struct{}, start time.Time) {
		tick := time.Tick(1 * time.Millisecond)

		for {
			select {
			case <-ch:
				return
			case <-tick:
				spinner.Message(component + " (elapsed: " + durafmt.Parse(time.Since(start)).LimitFirstN(2).String() + ")") // nolint: gomnd
			}
		}
	}(exit, time.Now())

	return exit
}

type vpc struct {
	ctx       context.Context
	out       io.Writer
	spinner   *yacspin.Spinner
	cidr      string
	id        api.ID
	c         *rest.HTTPClient
	outputDir string
	fs        *afero.Afero
}

func createVPC(v *vpc) (*api.Vpc, error) {
	_ = v.spinner.Start()
	exit := timer(v.spinner, "vpc")

	vpcService := core.NewVPCService(
		rest.NewVPCAPI(v.c),
		filesystem.NewVpcStore(
			config.DefaultVpcOutputs,
			config.DefaultVpcCloudFormationTemplate,
			path.Join(v.outputDir, config.DefaultVpcBaseDir),
			v.fs,
		),
		console.NewVPCReport(v.out, v.spinner, exit),
	)

	return vpcService.CreateVpc(v.ctx, api.CreateVpcOpts{
		ID:   v.id,
		Cidr: v.cidr,
	})
}

type cluster struct {
	ctx       context.Context
	out       io.Writer
	spinner   *yacspin.Spinner
	vpc       *api.Vpc
	cidr      string
	id        api.ID
	c         *rest.HTTPClient
	repoData  *repository.Data
	outputDir string
	repoDir   string
	fs        *afero.Afero
}

func createCluster(c *cluster) (*api.Cluster, error) {
	_ = c.spinner.Start()
	exit := timer(c.spinner, "cluster")

	clusterService := core.NewClusterService(
		rest.NewClusterAPI(c.c),
		filesystem.NewClusterStore(
			filesystem.Paths{
				ConfigFile: config.DefaultRepositoryConfig,
				BaseDir:    c.repoDir,
			},
			filesystem.Paths{
				ConfigFile: config.DefaultClusterConfig,
				BaseDir:    path.Join(c.outputDir, config.DefaultClusterBaseDir),
			},
			c.fs,
			c.repoData,
		),
		console.NewClusterReport(c.out, exit, c.spinner),
	)

	return clusterService.CreateCluster(c.ctx, api.ClusterCreateOpts{
		ID:                c.id,
		Cidr:              c.cidr,
		VpcID:             c.vpc.VpcID,
		VpcPrivateSubnets: c.vpc.PrivateSubnets,
		VpcPublicSubnets:  c.vpc.PublicSubnets,
	})
}

type externalSecrets struct {
	ctx       context.Context
	out       io.Writer
	spinner   *yacspin.Spinner
	c         *rest.HTTPClient
	outputDir string
	fs        *afero.Afero
	id        api.ID
}

func createExternalSecrets(e *externalSecrets) (*client.ExternalSecrets, error) {
	_ = e.spinner.Start()
	exit := timer(e.spinner, "external-secrets")

	externalSecretsService := core.NewExternalSecretsService(
		rest.NewExternalSecretsAPI(e.c),
		filesystem.NewExternalSecretsStore(
			filesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(e.outputDir, config.DefaultExternalSecretsBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(e.outputDir, config.DefaultExternalSecretsBaseDir),
			},
			filesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(e.outputDir, config.DefaultExternalSecretsBaseDir),
			},
			e.fs,
		),
		console.NewExternalSecretsReport(e.out, exit, e.spinner),
	)

	return externalSecretsService.CreateExternalSecrets(e.ctx, client.CreateExternalSecretsOpts{
		ID: e.id,
	})
}

type ingressController struct {
	ctx       context.Context
	out       io.Writer
	spinner   *yacspin.Spinner
	outputDir string
	c         *rest.HTTPClient
	fs        *afero.Afero
	id        api.ID
	vpc       *api.Vpc
}

func createAlbIngressController(c *ingressController) (*client.ALBIngressController, error) {
	_ = c.spinner.Start()
	exit := timer(c.spinner, "alb-ingress-controller")

	albIngressControllerService := core.NewALBIngressControllerService(
		rest.NewALBIngressControllerAPI(c.c),
		filesystem.NewALBIngressControllerStore(
			filesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(c.outputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(c.outputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			filesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(c.outputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			c.fs,
		),
		console.NewAlbIngressControllerReport(c.out, exit, c.spinner),
	)

	return albIngressControllerService.CreateALBIngressController(c.ctx, client.CreateALBIngressControllerOpts{
		ID:    c.id,
		VPCID: c.vpc.VpcID,
	})
}

type hostedZone struct {
	ctx       context.Context
	out       io.Writer
	spinner   *yacspin.Spinner
	repoData  *repository.Data
	c         *rest.HTTPClient
	outputDir string
	repoDir   string
	fs        *afero.Afero
	id        api.ID
}

func createPrimaryHostedZone(h *hostedZone) (*api.HostedZone, error) {
	_ = h.spinner.Start()
	exit := timer(h.spinner, "primary-hosted-zone")

	a := ask.New()

	domainService := core.NewDomainService(
		h.out,
		h.repoData,
		a,
		rest.NewDomainAPI(h.c),
		filesystem.NewDomainStore(
			h.repoData,
			filesystem.Paths{
				OutputFile:         config.DefaultDomainOutputsFile,
				CloudFormationFile: config.DefaultDomainCloudFormationTemplate,
				BaseDir:            path.Join(h.outputDir, config.DefaultDomainBaseDir),
			},
			filesystem.Paths{
				ConfigFile: config.DefaultRepositoryConfig,
				BaseDir:    h.repoDir,
			},
			h.fs,
		),
		console.NewDomainReport(h.out, exit, h.spinner),
		h.spinner,
	)

	return domainService.CreatePrimaryHostedZone(h.ctx, client.CreatePrimaryHostedZoneOpts{
		ID: h.id,
	})
}

type externalDNS struct {
	ctx       context.Context
	out       io.Writer
	spinner   *yacspin.Spinner
	id        api.ID
	outputDir string
	c         *rest.HTTPClient
	fs        *afero.Afero
	d         *api.HostedZone
}

func createExternalDNS(e *externalDNS) (*client.ExternalDNS, error) {
	_ = e.spinner.Start()
	exit := timer(e.spinner, "external-dns")

	externalDNSService := core.NewExternalDNSService(
		rest.NewExternalDNSAPI(e.c),
		filesystem.NewExternalDNSStore(
			filesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(e.outputDir, config.DefaultExternalDNSBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(e.outputDir, config.DefaultExternalDNSBaseDir),
			},
			filesystem.Paths{
				OutputFile: config.DefaultKubeOutputsFile,
				BaseDir:    path.Join(e.outputDir, config.DefaultExternalDNSBaseDir),
			},
			e.fs,
		),
		console.NewExternalDNSReport(e.out, exit, e.spinner),
	)

	return externalDNSService.CreateExternalDNS(e.ctx, client.CreateExternalDNSOpts{
		ID:           e.id,
		HostedZoneID: e.d.HostedZoneID,
		Domain:       e.d.Domain,
	})
}

type githubClient struct {
	ctx                   context.Context
	out                   io.Writer
	spinner               *yacspin.Spinner
	repoDir, organisation string
	id                    api.ID
	repoData              *repository.Data
	c                     *rest.HTTPClient
	fs                    *afero.Afero
	auth                  github2.Authenticator
}

func createGithubRepo(c *githubClient) (client.GithubService, *client.GithubRepository, error) {
	g, err := github.New(c.ctx, c.auth)
	if err != nil {
		return nil, nil, err
	}

	repo, err := git.GithubRepoFullName(c.organisation, c.repoDir)
	if err != nil {
		return nil, nil, err
	}

	_ = c.spinner.Start()
	exit := timer(c.spinner, "github")

	githubService := core.NewGithubService(
		rest.NewGithubAPI(
			c.out,
			ask.New(),
			c.spinner,
			rest.NewParameterAPI(c.c),
			g,
		),
		filesystem.NewGithubStore(
			filesystem.Paths{
				ConfigFile: config.DefaultRepositoryConfig,
				BaseDir:    c.repoDir,
			},
			c.repoData,
			c.fs,
		),
		console.NewGithubReport(c.out, exit, c.spinner),
	)

	githubRepo, err := githubService.ReadyGithubInfrastructureRepository(c.ctx, client.ReadyGithubInfrastructureRepositoryOpts{
		ID:           c.id,
		Organisation: c.organisation,
		Repository:   repo,
	})

	return githubService, githubRepo, err
}

type argocdSetup struct {
	ctx           context.Context
	out           io.Writer
	spinner       *yacspin.Spinner
	id            api.ID
	githubRepo    *client.GithubRepository
	githubService client.GithubService
	d             *api.HostedZone
	outputDir     string
	repoDir       string
	organisation  string
	fs            *afero.Afero
	c             *rest.HTTPClient
	repoData      *repository.Data
}

// nolint: funlen
func createArgoCD(a *argocdSetup) (*client.ArgoCD, error) {
	argoBaseDir := path.Join(a.outputDir, config.DefaultArgoCDBaseDir)

	_ = a.spinner.Start()
	exit := timer(a.spinner, "argocd")

	certStore := filesystem.NewCertificateStore(
		a.repoData,
		filesystem.Paths{
			OutputFile:         config.DefaultCertificateOutputsFile,
			CloudFormationFile: config.DefaultCertificateCloudFormationTemplate,
			BaseDir:            path.Join(argoBaseDir, config.DefaultCertificateBaseDir),
		},
		filesystem.Paths{
			ConfigFile: config.DefaultRepositoryConfig,
			BaseDir:    a.repoDir,
		},
		a.fs,
	)

	manifestStore := filesystem.NewManifestStore(
		filesystem.Paths{
			OutputFile: config.DefaultKubeOutputsFile,
			BaseDir:    path.Join(argoBaseDir, config.DefaultExternalSecretsBaseDir),
		},
		a.fs,
	)

	paramStore := filesystem.NewParameterStore(
		filesystem.Paths{
			OutputFile: config.DefaultParameterOutputsFile,
			BaseDir:    path.Join(argoBaseDir, config.DefaultParameterBaseDir),
		},
		a.fs,
	)

	argoService := core.NewArgoCDService(
		rest.NewArgoCDAPI(
			a.githubService,
			rest.NewParameterAPI(a.c),
			rest.NewManifestAPI(a.c),
			rest.NewCertificateAPI(a.c),
			a.c,
		),
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
			a.repoData,
			filesystem.Paths{
				ConfigFile: config.DefaultRepositoryConfig,
				BaseDir:    a.repoDir,
			},
			paramStore,
			certStore,
			manifestStore,
			a.fs,
		),
		console.NewArgoCDReport(a.out, exit, a.spinner),
	)

	return argoService.CreateArgoCD(a.ctx, client.CreateArgoCDOpts{
		ID:                 a.id,
		Domain:             a.d.Domain,
		FQDN:               a.d.FQDN,
		HostedZoneID:       a.d.HostedZoneID,
		GithubOrganisation: a.organisation,
		Repository:         a.githubRepo,
	})
}
