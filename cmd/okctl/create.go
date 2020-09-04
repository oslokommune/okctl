package main

import (
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"

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
			opts.Organisation = github.DefaultOrg

			if !o.HostedZoneIsCreated(opts.Environment) {
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

			o.SetGithubOrganisationName(opts.Organisation, opts.Environment)
			err := o.WriteCurrentRepoData()
			if err != nil {
				return err
			}

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
				ClusterName: opts.ClusterName,
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

			o.SetHostedZoneIsCreated(true, opts.Environment)

			err = o.WriteCurrentRepoData()
			if err != nil {
				return err
			}

			a := ask.New()

			if !o.HostedZoneIsDelegated(opts.Environment) {
				err = a.ConfirmPostingNameServers(o.Out, d.Domain, d.NameServers)
				if err != nil {
					return err
				}

				o.SetHostedZoneIsDelegated(true, opts.Environment)

				err := o.WriteCurrentRepoData()
				if err != nil {
					return err
				}
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

			g, err := github.New(opts.Organisation, o.CredentialsProvider.Github())
			if err != nil {
				return err
			}

			githubRepo := o.GithubRepository(opts.Environment)
			if len(githubRepo.Name) == 0 || len(githubRepo.GitURL) == 0 {
				repo, err := git.GithubRepoFullName(opts.Organisation, repoDir)
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

				githubRepo.Name = selectedRepo.GetName()
				githubRepo.GitURL = selectedRepo.GetGitURL()

				o.SetGithubRepository(githubRepo, opts.Environment)
				err = o.WriteCurrentRepoData()
				if err != nil {
					return err
				}
			}

			team := o.GithubTeamName(opts.Environment)
			if len(team) == 0 {
				teams, err := g.Teams()
				if err != nil {
					return err
				}

				selectedTeam, err := a.SelectTeam(teams)
				if err != nil {
					return err
				}

				o.SetGithubTeamName(selectedTeam.GetName(), opts.Environment)
				err = o.WriteCurrentRepoData()
				if err != nil {
					return err
				}

				team = selectedTeam.GetName()
			}

			oauthApp := o.GithubOauthApp(opts.Environment)
			if len(oauthApp.Name) == 0 || len(oauthApp.ClientID) == 0 || len(oauthApp.ClientSecretPath) == 0 {
				app, err := a.CreateOauthApp(o.Out, ask.OauthAppOpts{
					Organisation: opts.Organisation,
					Name:         fmt.Sprintf("okctl-argocd-%s-%s", opts.RepositoryName, opts.Environment),
					URL:          fmt.Sprintf("https://%s", cert.Domain),
					CallbackURL:  fmt.Sprintf("https://%s/api/dex/callback", cert.Domain),
				})
				if err != nil {
					return err
				}

				clientSecret, err := c.CreateSecret(&api.CreateSecretOpts{
					AWSAccountID:   opts.AWSAccountID,
					RepositoryName: opts.RepositoryName,
					Environment:    opts.Environment,
					Name:           "client_secret",
					Secret:         app.ClientSecret,
				})
				if err != nil {
					return err
				}

				oauthApp.Name = app.Name
				oauthApp.ClientID = app.ClientID
				oauthApp.ClientSecretPath = clientSecret.Path

				o.SetGithubOauthApp(oauthApp, opts.Environment)
				err = o.WriteCurrentRepoData()
				if err != nil {
					return err
				}
			}

			deployKey := o.GithubDeployKey(opts.Environment)
			if len(deployKey.Title) == 0 || deployKey.ID == 0 || len(deployKey.Path) == 0 {
				key, err := keypair.New(keypair.DefaultRandReader(), keypair.DefaultBitSize).Generate()
				if err != nil {
					return err
				}

				dk, err := g.CreateDeployKey(githubRepo.Name, "okctl-argocd-read-key", string(key.PublicKey))
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

				deployKey.Title = dk.GetTitle()
				deployKey.ID = dk.GetID()
				deployKey.Path = privateKey.Path

				o.SetGithubDeployKey(deployKey, opts.Environment)
				err = o.WriteCurrentRepoData()
				if err != nil {
					return err
				}
			}

			argocd := o.ArgoCD(opts.Environment)
			if len(argocd.SecretKeyPath) == 0 {
				secretKey, err := c.CreateSecret(&api.CreateSecretOpts{
					AWSAccountID:   opts.AWSAccountID,
					RepositoryName: opts.RepositoryName,
					Environment:    opts.Environment,
					Name:           "argocd_secret_key",
					Secret:         uuid.New().String(),
				})
				if err != nil {
					return err
				}

				argocd.SecretKeyPath = secretKey.Path

				o.SetArgoCD(argocd, opts.Environment)
				err = o.WriteCurrentRepoData()
				if err != nil {
					return err
				}
			}

			_, err = c.CreateExternalSecrets(&api.CreateExternalSecretsOpts{
				Manifests: []api.Manifest{
					{
						Name:      "argocd-privatekey",
						Namespace: "argocd",
						Data: []api.Data{
							{
								Name: "ssh-private-key",
								Key:  deployKey.Path,
							},
						},
					},
					{
						Name:      "argocd-secret",
						Namespace: "argocd",
						Data: []api.Data{
							{
								Name: "dex.github.clientSecret",
								Key:  oauthApp.ClientSecretPath,
							},
							{
								Name: "server.secretkey",
								Key:  argocd.SecretKeyPath,
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}

			_, err = c.CreateArgoCD(&api.CreateArgoCDOpts{
				ClusterName:         opts.ClusterName,
				Repository:          opts.RepositoryName,
				Environment:         opts.Environment,
				ArgoDomain:          cert.Domain,
				ArgoCertificateARN:  cert.CertificateARN,
				GithubOrganisation:  opts.Organisation,
				GithubTeam:          team,
				GithubRepoURL:       githubRepo.GitURL,
				GithubRepoName:      githubRepo.Name,
				GithubOauthClientID: oauthApp.ClientID,
				PrivateKeyName:      "argocd-privatekey",
				PrivateKeyKey:       "ssh-private-key",
			})

			argocd = o.ArgoCD(opts.Environment)
			argocd.URL = fmt.Sprintf("https://%s", cert.Domain)
			o.SetArgoCD(argocd, opts.Environment)
			err = o.WriteCurrentRepoData()
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
