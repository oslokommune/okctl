package state

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/constant"
)

const (
	// TypeInfrastructure identifies a repository infrastructure as code
	TypeInfrastructure = "infrastructure"
	// TypeApplication identifies a repository for applications
	TypeApplication = "application"

	// SubnetTypePublic is a public subnet
	SubnetTypePublic = "public"
	// SubnetTypePrivate is a private subnet
	SubnetTypePrivate = "private"
)

// Repository stores the configured state of a repository
// as used by okctl
type Repository struct {
	Metadata Metadata
	Clusters map[string]Cluster
}

// HasEnvironment tests if the environment exists in this repository
func (r Repository) HasEnvironment(environment string) (hasEnvironment bool) {
	for cluster := range r.Clusters {
		if cluster == environment {
			return true
		}
	}

	return hasEnvironment
}

// Validate the repository
func (r Repository) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(r.Metadata),
		validation.Field(r.Clusters),
	)
}

// Metadata contains repository metadata
type Metadata struct {
	Name      string
	Region    string
	OutputDir string
}

// Validate the metadata
func (m Metadata) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name,
			validation.Required,
		),
		validation.Field(&m.Region,
			validation.Required,
			validation.In(func() []interface{} {
				var o []interface{}
				for _, r := range v1alpha1.SupportedRegions() {
					o = append(o, r)
				}
				return o
			}()...),
		),
		validation.Field(&m.OutputDir,
			validation.Required,
		),
	)
}

// Cluster represents an okctl created
// cluster
type Cluster struct {
	Name         string
	Environment  string
	AWSAccountID string
	HostedZone   map[string]HostedZone
	VPC          VPC
	Certificates map[string]Certificate
	Github       Github
	ArgoCD       ArgoCD
	IdentityPool IdentityPool
	Monitoring   Monitoring
}

const (
	envMinLength = 3
	envMaxLength = 10
)

// Validate the cluster data
func (c Cluster) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Environment,
			validation.Required,
			validation.Length(envMinLength, envMaxLength),
		),
		validation.Field(&c.HostedZone, validation.Required),
		validation.Field(&c.VPC),
		validation.Field(&c.Certificates),
	)
}

// Monitoring contains state information about
// the monitoring setup
type Monitoring struct {
	DashboardURL string
}

// Validate the monitoring struct
func (m Monitoring) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DashboardURL, validation.Required),
	)
}

// IdentityPool contains the state about an identity
// user pool
type IdentityPool struct {
	UserPoolID string
	AuthDomain string
	Alias      RecordSetAlias
	Clients    map[string]IdentityPoolClient
	Users      map[string]IdentityPoolUser
}

// Validate the identity pool
func (p IdentityPool) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserPoolID, validation.Required),
		validation.Field(&p.AuthDomain, validation.Required),
		validation.Field(&p.Alias, validation.Required),
		validation.Field(&p.Clients, validation.Required),
		validation.Field(&p.Users),
	)
}

// Validate pool
func (p IdentityPoolUser) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required),
	)
}

// IdentityPoolClient contains the state about an
// identity client
type IdentityPoolClient struct {
	Purpose      string
	CallbackURL  string
	ClientID     string
	ClientSecret ClientSecret
}

// IdentityPoolUser output
type IdentityPoolUser struct {
	Email string
}

// Validate the identity pool client
func (c IdentityPoolClient) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Purpose, validation.Required),
		validation.Field(&c.CallbackURL, validation.Required),
		validation.Field(&c.ClientID, validation.Required),
	)
}

// RecordSetAlias contains state about an alias
// record set
type RecordSetAlias struct {
	AliasDomain     string
	AliasHostedZone string
}

// Validate the record set alias
func (a RecordSetAlias) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.AliasDomain, validation.Required),
		validation.Field(&a.AliasHostedZone, validation.Required),
	)
}

// Certificate contains state about a certificate
type Certificate struct {
	Domain string
	ARN    string
}

// Validate the data
func (c Certificate) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Domain, validation.Required),
		validation.Field(&c.ARN, validation.Required),
	)
}

// VPC contains state about the VPC
type VPC struct {
	VpcID   string
	CIDR    string
	Subnets map[string][]VPCSubnet
}

// VPCSubnet is a vpc subnet
type VPCSubnet struct {
	CIDR             string
	AvailabilityZone string
}

// ArgoCD contains information about the
// argocd setup
type ArgoCD struct {
	SiteURL   string
	Domain    string
	SecretKey SecretKeySecret
}

// SecretKeySecret contains state about
// an argo cd secret key
type SecretKeySecret struct {
	Name    string
	Path    string
	Version int64
}

// Github contains information about the
// clusters configuration towards github
type Github struct {
	Organisation string
	OauthApp     map[string]GithubOauthApp
	Repositories map[string]GithubRepository
}

// GithubRepository contains github repository data
type GithubRepository struct {
	Name         string
	FullName     string
	Organization string
	Types        []string
	GitURL       string
	DeployKey    DeployKey
}

// Validate the data
func (r GithubRepository) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.FullName, validation.Required),
		validation.Field(&r.Types, validation.Required),
		validation.Field(&r.GitURL, validation.Required),
		validation.Field(&r.Organization, validation.Required),
	)
}

// GithubOauthApp contains github oauth application data
type GithubOauthApp struct {
	Team         string
	Name         string
	SiteURL      string
	CallbackURL  string
	ClientID     string
	ClientSecret ClientSecret
}

// Validate the data
func (a GithubOauthApp) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Team, validation.Required),
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.SiteURL, validation.Required),
		validation.Field(&a.CallbackURL, validation.Required),
		validation.Field(&a.ClientSecret, validation.Required),
		validation.Field(&a.ClientID, validation.Required),
	)
}

// ClientSecret contains state about
// an oauth app client secret
type ClientSecret struct {
	Name    string
	Path    string
	Version int64
}

// Validate the data
func (s ClientSecret) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Path, validation.Required),
		validation.Field(&s.Version, validation.Required),
	)
}

// DeployKey contains github deploy key data
type DeployKey struct {
	Title            string
	ID               int64
	PublicKey        string
	PrivateKeySecret PrivateKeySecret
}

// Validate the data
func (k DeployKey) Validate() error {
	return validation.ValidateStruct(&k,
		validation.Field(&k.Title, validation.Required),
		validation.Field(&k.ID, validation.Required),
		validation.Field(&k.PublicKey, validation.Required),
		validation.Field(&k.PrivateKeySecret, validation.Required),
	)
}

// PrivateKeySecret contains information
// about a private key
type PrivateKeySecret struct {
	Name    string
	Path    string
	Version int64
}

// Validate the data
func (s PrivateKeySecret) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Path, validation.Required),
	)
}

// HostedZone contains information about the
// clusters hostedzone delegation
type HostedZone struct {
	IsDelegated bool
	Primary     bool
	Managed     bool
	ID          string
	Domain      string
	FQDN        string
	NameServers []string
}

// Validate the hostedzone
func (h HostedZone) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Domain, validation.Required),
		validation.Field(&h.FQDN, validation.Required),
		validation.Field(h.NameServers, validation.Each(validation.Required)),
	)
}

// NewRepository returns repository data with defaults set
func NewRepository() *Repository {
	return &Repository{
		Metadata: Metadata{
			Region:    v1alpha1.RegionEuWest1,
			OutputDir: constant.DefaultOutputDirectory,
		},
		Clusters: map[string]Cluster{},
	}
}
