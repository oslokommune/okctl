package main

import (
	"fmt"
	"io/ioutil"
	"path"

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
	"github.com/oslokommune/okctl/pkg/domain"
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
	DomainName     string
	FQDN           string
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
		validation.Field(&o.DomainName, validation.Required),
		validation.Field(&o.FQDN, validation.Required),
		validation.Field(&o.Organisation, validation.Required),
	)
}

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
			opts.DomainName = o.PrimaryDomain(opts.Environment)
			opts.FQDN = o.PrimaryFQDN(opts.Environment)
			opts.Organisation = github.DefaultOrg

			// FIXME: Move this into the domain ask thingy
			if !o.HostedZoneIsCreated(opts.DomainName, opts.Environment) {
				d, err := domain.NewDefaultWithSurvey(opts.RepositoryName, opts.Environment)
				if err != nil {
					return fmt.Errorf("failed to get domain name: %w", err)
				}
				opts.DomainName = d.Domain
				opts.FQDN = d.FQDN
			}

			err := opts.Validate()
			if err != nil {
				return errors.E(err, "failed to validate create cluster options", errors.Invalid)
			}

			return o.Initialise(opts.Environment, opts.AWSAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			repoDir, err := o.GetRepoDir()
			if err != nil {
				return err
			}

			outputDir, err := o.GetRepoOutputDir(opts.Environment)
			if err != nil {
				return err
			}

			// Discarding the output for now, until we restructure
			// the API to return everything we need to write
			// the result ourselves
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
				return err
			}

			vpcService := core.NewVPCService(
				rest.NewVPCAPI(c),
				filesystem.NewVpcStore(
					config.DefaultVpcOutputs,
					config.DefaultVpcCloudFormationTemplate,
					path.Join(outputDir, config.DefaultVpcBaseDir),
					o.FileSystem,
				),
			)

			vpc, err := vpcService.CreateVpc(o.Ctx, api.CreateVpcOpts{
				ID:   id,
				Cidr: opts.Cidr,
			})
			if err != nil {
				return err
			}

			clusterService := core.NewClusterService(
				rest.NewClusterAPI(c),
				filesystem.NewClusterStore(
					filesystem.Paths{
						ConfigFile: config.DefaultRepositoryConfig,
						BaseDir:    repoDir,
					},
					filesystem.Paths{
						ConfigFile: config.DefaultClusterConfig,
						BaseDir:    path.Join(outputDir, config.DefaultClusterBaseDir),
					},
					o.FileSystem,
					o.RepoData,
				),
			)

			_, err = clusterService.CreateCluster(o.Ctx, api.ClusterCreateOpts{
				ID:                id,
				Cidr:              opts.Cidr,
				VpcID:             vpc.VpcID,
				VpcPrivateSubnets: vpc.PrivateSubnets,
				VpcPublicSubnets:  vpc.PublicSubnets,
			})
			if err != nil {
				return err
			}

			externalSecretsService := core.NewExternalSecretsService(
				rest.NewExternalSecretsAPI(c),
				filesystem.NewExternalSecretsStore(
					filesystem.Paths{
						OutputFile:         config.DefaultPolicyOutputFile,
						CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
						BaseDir:            path.Join(outputDir, config.DefaultExternalSecretsBaseDir),
					},
					filesystem.Paths{
						OutputFile: config.DefaultServiceAccountOutputsFile,
						ConfigFile: config.DefaultServiceAccountConfigFile,
						BaseDir:    path.Join(outputDir, config.DefaultExternalSecretsBaseDir),
					},
					filesystem.Paths{
						OutputFile:  config.DefaultHelmOutputsFile,
						ReleaseFile: config.DefaultHelmReleaseFile,
						ChartFile:   config.DefaultHelmChartFile,
						BaseDir:     path.Join(outputDir, config.DefaultExternalSecretsBaseDir),
					},
					o.FileSystem,
				),
			)

			_, err = externalSecretsService.CreateExternalSecrets(o.Ctx, client.CreateExternalSecretsOpts{
				ID: id,
			})

			// alb ingress
			albIngressControllerService := core.NewALBIngressControllerService(
				rest.NewALBIngressControllerAPI(c),
				filesystem.NewALBIngressControllerStore(
					filesystem.Paths{
						OutputFile:         config.DefaultPolicyOutputFile,
						CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
						BaseDir:            path.Join(outputDir, config.DefaultAlbIngressControllerBaseDir),
					},
					filesystem.Paths{
						OutputFile: config.DefaultServiceAccountOutputsFile,
						ConfigFile: config.DefaultServiceAccountConfigFile,
						BaseDir:    path.Join(outputDir, config.DefaultAlbIngressControllerBaseDir),
					},
					filesystem.Paths{
						OutputFile:  config.DefaultHelmOutputsFile,
						ReleaseFile: config.DefaultHelmReleaseFile,
						ChartFile:   config.DefaultHelmChartFile,
						BaseDir:     path.Join(outputDir, config.DefaultAlbIngressControllerBaseDir),
					},
					o.FileSystem,
				),
			)

			_, err = albIngressControllerService.CreateALBIngressController(o.Ctx, client.CreateALBIngressControllerOpts{
				ID:    id,
				VPCID: vpc.VpcID,
			})

			domainService := core.NewDomainService(
				rest.NewDomainAPI(c),
				filesystem.NewDomainStore(
					o.RepoData,
					filesystem.Paths{
						OutputFile:         config.DefaultDomainOutputsFile,
						CloudFormationFile: config.DefaultDomainCloudFormationTemplate,
						BaseDir:            path.Join(outputDir, config.DefaultDomainBaseDir),
					},
					filesystem.Paths{
						ConfigFile: config.DefaultRepositoryConfig,
						BaseDir:    repoDir,
					},
					o.FileSystem,
				),
			)

			d, err := domainService.CreateHostedZone(o.Ctx, api.CreateHostedZoneOpts{
				ID:     id,
				Domain: opts.DomainName,
				FQDN:   opts.FQDN,
			})

			// FIXME: Move this stuff into the domain create
			a := ask.New()

			if !o.HostedZoneIsDelegated(opts.DomainName, opts.Environment) {
				err = a.ConfirmPostingNameServers(o.Out, d.Domain, d.NameServers)
				if err != nil {
					return err
				}

				o.SetHostedZoneIsDelegated(true, opts.DomainName, opts.Environment)

				err := o.WriteCurrentRepoData()
				if err != nil {
					return err
				}
			}

			// external dns
			externalDNSService := core.NewExternalDNSService(
				rest.NewExternalDNSAPI(c),
				filesystem.NewExternalDNSStore(
					filesystem.Paths{
						OutputFile:         config.DefaultPolicyOutputFile,
						CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
						BaseDir:            path.Join(outputDir, config.DefaultExternalDNSBaseDir),
					},
					filesystem.Paths{
						OutputFile: config.DefaultServiceAccountOutputsFile,
						ConfigFile: config.DefaultServiceAccountConfigFile,
						BaseDir:    path.Join(outputDir, config.DefaultExternalDNSBaseDir),
					},
					filesystem.Paths{
						OutputFile: config.DefaultKubeOutputsFile,
						BaseDir:    path.Join(outputDir, config.DefaultExternalDNSBaseDir),
					},
					o.FileSystem,
				),
			)

			_, err = externalDNSService.CreateExternalDNS(o.Ctx, client.CreateExternalDNSOpts{
				ID:           id,
				HostedZoneID: d.HostedZoneID,
				Domain:       d.Domain,
			})

			g, err := github.New(o.Ctx, o.CredentialsProvider.Github())
			if err != nil {
				return err
			}

			repo, err := git.GithubRepoFullName(opts.Organisation, repoDir)
			if err != nil {
				return err
			}

			githubService := core.NewGithubService(
				rest.NewGithubAPI(
					rest.NewParameterAPI(c),
					g,
				),
				filesystem.NewGithubStore(
					filesystem.Paths{
						ConfigFile: config.DefaultRepositoryConfig,
						BaseDir:    repoDir,
					},
					o.RepoData,
					o.FileSystem,
				),
			)

			githubRepo, err := githubService.ReadyGithubInfrastructureRepository(o.Ctx, client.ReadyGithubInfrastructureRepositoryOpts{
				ID:           id,
				Organisation: opts.Organisation,
				Repository:   repo,
			})
			if err != nil {
				return err
			}

			// Everything below here goes into ArgoCD

			argoBaseDir := path.Join(outputDir, config.DefaultArgoCDBaseDir)

			certStore := filesystem.NewCertificateStore(
				o.RepoData,
				filesystem.Paths{
					OutputFile:         config.DefaultCertificateOutputsFile,
					CloudFormationFile: config.DefaultCertificateCloudFormationTemplate,
					BaseDir:            path.Join(argoBaseDir, config.DefaultCertificateBaseDir),
				},
				filesystem.Paths{
					ConfigFile: config.DefaultRepositoryConfig,
					BaseDir:    repoDir,
				},
				o.FileSystem,
			)

			manifestStore := filesystem.NewManifestStore(
				filesystem.Paths{
					OutputFile: config.DefaultKubeOutputsFile,
					BaseDir:    path.Join(argoBaseDir, config.DefaultExternalSecretsBaseDir),
				},
				o.FileSystem,
			)

			paramStore := filesystem.NewParameterStore(
				filesystem.Paths{
					OutputFile: config.DefaultParameterOutputsFile,
					BaseDir:    path.Join(argoBaseDir, config.DefaultParameterBaseDir),
				},
				o.FileSystem,
			)

			argoService := core.NewArgoCDService(
				rest.NewArgoCDAPI(
					githubService,
					rest.NewParameterAPI(c),
					rest.NewManifestAPI(c),
					rest.NewCertificateAPI(c),
					c,
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
					o.RepoData,
					filesystem.Paths{
						ConfigFile: config.DefaultRepositoryConfig,
						BaseDir:    repoDir,
					},
					paramStore,
					certStore,
					manifestStore,
					o.FileSystem,
				),
			)

			_, err = argoService.CreateArgoCD(o.Ctx, client.CreateArgoCDOpts{
				ID:                 id,
				Domain:             d.Domain,
				FQDN:               d.FQDN,
				HostedZoneID:       d.HostedZoneID,
				GithubOrganisation: opts.Organisation,
				Repository:         githubRepo,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.Cidr, "cidr", "c", defaultCidr,
		"CIDR block the AWS VPC and subnets are created within")

	return cmd
}
