package main

import (
	"fmt"
	"io/ioutil"

	"github.com/oslokommune/okctl/pkg/keypair"

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
	)
}

// nolint: funlen gocyclo
func buildCreateClusterCommand(o *okctl.Okctl) *cobra.Command {
	// This should probably be a local struct, since we do much
	// more now then before
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
			opts.DomainName = o.Domain(opts.Environment)
			opts.FQDN = o.FQDN(opts.Environment)

			if len(o.Domain(opts.Environment)) == 0 {
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
			// Discarding the output for now, until we restructure
			// the API to return everything we need to write
			// the result ourselves
			c := client.New(o.Debug, ioutil.Discard, o.ServerURL)

			vpc, err := c.CreateVpc(&api.CreateVpcOpts{
				AwsAccountID: opts.AWSAccountID,
				ClusterName:  opts.ClusterName,
				Env:          opts.Environment,
				RepoName:     opts.RepositoryName,
				Cidr:         opts.Cidr,
				Region:       opts.Region,
			})
			if err != nil {
				return err
			}

			err = c.CreateCluster(&api.ClusterCreateOpts{
				Environment:       opts.Environment,
				AWSAccountID:      opts.AWSAccountID,
				Cidr:              opts.Cidr,
				RepositoryName:    opts.RepositoryName,
				Region:            opts.Region,
				ClusterName:       opts.ClusterName,
				VpcID:             vpc.ID,
				VpcPrivateSubnets: vpc.PrivateSubnets,
				VpcPublicSubnets:  vpc.PublicSubnets,
			})
			if err != nil {
				return err
			}

			policy, err := c.CreateExternalSecretsPolicy(&api.CreateExternalSecretsPolicyOpts{
				Repository:  opts.RepositoryName,
				Environment: opts.Environment,
			})
			if err != nil {
				return err
			}

			err = c.CreateExternalSecretsServiceAccount(&api.CreateExternalSecretsServiceAccountOpts{
				CreateServiceAccountOpts: api.CreateServiceAccountOpts{
					ClusterName:  opts.ClusterName,
					Environment:  opts.Environment,
					Region:       opts.Region,
					AWSAccountID: opts.AWSAccountID,
					PolicyArn:    policy.PolicyARN,
				},
			})
			if err != nil {
				return err
			}

			_, err = c.CreateExternalSecretsHelmChart(&api.CreateExternalSecretsHelmChartOpts{
				Repository:  opts.RepositoryName,
				Environment: opts.Environment,
			})
			if err != nil {
				return err
			}

			policy, err = c.CreateAlbIngressControllerPolicy(&api.CreateAlbIngressControllerPolicyOpts{
				Repository:  opts.RepositoryName,
				Environment: opts.Environment,
			})
			if err != nil {
				return err
			}

			err = c.CreateAlbIngressControllerServiceAccount(&api.CreateAlbIngressControllerServiceAccountOpts{
				CreateServiceAccountOpts: api.CreateServiceAccountOpts{
					ClusterName:  opts.ClusterName,
					Environment:  opts.Environment,
					Region:       opts.Region,
					AWSAccountID: opts.AWSAccountID,
					PolicyArn:    policy.PolicyARN,
				},
			})
			if err != nil {
				return err
			}

			_, err = c.CreateAlbIngressControllerHelmChart(&api.CreateAlbIngressControllerHelmChartOpts{
				ClusterName: opts.ClusterName,
				Repository:  opts.RepositoryName,
				Environment: opts.Environment,
				VpcID:       vpc.ID,
				Region:      opts.Region,
			})
			if err != nil {
				return err
			}

			d, err := c.CreateDomain(&api.CreateDomainOpts{
				Repository:  opts.RepositoryName,
				Environment: opts.Environment,
				FQDN:        opts.FQDN,
				Domain:      opts.DomainName,
			})
			if err != nil {
				return err
			}

			a := ask.New()

			err = a.ConfirmPostingNameServers(o.Out, d.Domain, d.NameServers)
			if err != nil {
				return err
			}

			policy, err = c.CreateExternalDNSPolicy(&api.CreateExternalDNSPolicyOpts{
				Repository:  opts.RepositoryName,
				Environment: opts.Environment,
			})
			if err != nil {
				return err
			}

			err = c.CreateExternalDNSServiceAccount(&api.CreateExternalDNSServiceAccountOpts{
				CreateServiceAccountOpts: api.CreateServiceAccountOpts{
					ClusterName:  opts.ClusterName,
					Environment:  opts.Environment,
					Region:       opts.Region,
					AWSAccountID: opts.AWSAccountID,
					PolicyArn:    policy.PolicyARN,
				},
			})
			if err != nil {
				return err
			}

			_, err = c.CreateExternalDNSKubeDeployment(&api.CreateExternalDNSKubeDeploymentOpts{
				HostedZoneID: d.HostedZoneID,
				DomainFilter: d.Domain,
			})
			if err != nil {
				return err
			}

			cert, err := c.CreateCertificate(&api.CreateCertificateOpts{
				Repository:   opts.RepositoryName,
				Environment:  opts.Environment,
				FQDN:         fmt.Sprintf("argocd.%s", opts.FQDN),
				Domain:       fmt.Sprintf("argocd.%s", opts.DomainName),
				HostedZoneID: d.HostedZoneID,
			})
			if err != nil {
				return err
			}

			repoDir, err := o.GetRepoDir()
			if err != nil {
				return err
			}

			repo, err := git.GithubRepoFullName("oslokommune", repoDir)
			if err != nil {
				return err
			}

			g, err := github.New("oslokommune", o.CredentialsProvider.Github())
			if err != nil {
				return err
			}

			repos, err := g.Repositories()
			if err != nil {
				return err
			}

			selectedRepo, err := a.SelectInfrastructureRepository(repo, repos)
			if err != nil {
				return err
			}

			teams, err := g.Teams()
			if err != nil {
				return err
			}

			selectedTeam, err := a.SelectTeam(teams)
			if err != nil {
				return err
			}

			oauthApp, err := a.CreateOauthApp(o.Out, ask.OauthAppOpts{
				Organisation: "oslokommune",
				Name:         fmt.Sprintf("okctl-argocd-%s-%s", opts.RepositoryName, opts.Environment),
				URL:          fmt.Sprintf("https://%s", cert.Domain),
				CallbackURL:  fmt.Sprintf("https://%s/api/dex/callback", cert.Domain),
			})
			if err != nil {
				return err
			}

			key, err := keypair.New(keypair.DefaultRandReader(), keypair.DefaultBitSize).Generate()
			if err != nil {
				return err
			}

			_, err = g.CreateDeployKey(selectedRepo.GetName(), "okctl-argocd-read-key", string(key.PublicKey))
			if err != nil {
				return err
			}

			clientSecret, err := c.CreateSecret(&api.CreateSecretOpts{
				AWSAccountID:   opts.AWSAccountID,
				RepositoryName: opts.RepositoryName,
				Environment:    opts.Environment,
				Name:           "client_secret",
				Secret:         oauthApp.ClientSecret,
			})
			if err != nil {
				return err
			}

			privateKey, err := c.CreateSecret(&api.CreateSecretOpts{
				AWSAccountID:   opts.AWSAccountID,
				RepositoryName: opts.RepositoryName,
				Environment:    opts.Environment,
				Name:           "private_key",
				Secret:         string(key.PrivateKey),
			})
			if err != nil {
				return err
			}

			_, err = c.CreateExternalSecrets(&api.CreateExternalSecretsOpts{
				Manifests: []api.Manifest{
					{
						Name:      "argocd-privatekey",
						Namespace: "argocd",
						Data: []api.Data{
							{
								Key:  "ssh-private-key",
								Name: privateKey.Path,
							},
						},
					},
					{
						Name:      "argocd-secret",
						Namespace: "argocd",
						Data: []api.Data{
							{
								Key:  "dex.github.clientSecret",
								Name: clientSecret.Path,
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}

			_, err = c.CreateArgoCD(&api.CreateArgoCDOpts{
				ClusterName:                 opts.ClusterName,
				Repository:                  opts.RepositoryName,
				Environment:                 opts.Environment,
				ArgoDomain:                  cert.Domain,
				ArgoCertificateARN:          cert.CertificateARN,
				GithubOrganisation:          "oslokommune",
				GithubTeam:                  selectedTeam.GetName(),
				GithubRepoURL:               selectedRepo.GetGitURL(),
				GithubRepoName:              selectedRepo.GetName(),
				GithubOauthClientID:         oauthApp.ClientID,
				GithubOauthClientSecretPath: clientSecret.Path,
				GithubDeployKeySecretPath:   privateKey.Path,
				PrivateKeyName:              "argocd-privatekey",
				PrivateKeyKey:               "ssh-private-key",
			})

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.Cidr, "cidr", "c", defaultCidr,
		"CIDR block the AWS VPC and subnets are created within")

	return cmd
}
