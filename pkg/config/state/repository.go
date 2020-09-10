package state

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
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
	Name      string
	Region    string
	OutputDir string
	Clusters  map[string]*Cluster
}

// ClusterForEnv returns the cluster for the given environment
func (d *Repository) ClusterForEnv(env string) *Cluster {
	if c, ok := d.Clusters[env]; ok {
		return c
	}

	return nil
}

// Validate the provided data
func (d *Repository) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Name,
			validation.Required,
		),
		validation.Field(&d.Region,
			validation.Required,
			validation.In(func() []interface{} {
				var o []interface{}
				for _, r := range v1alpha1.SupportedRegions() {
					o = append(o, r)
				}
				return o
			}()...),
		),
		validation.Field(&d.OutputDir,
			validation.Required,
		),
		validation.Field(&d.Clusters),
	)
}

// Cluster represents an okctl created
// cluster
type Cluster struct {
	Name         string
	Environment  string
	AWSAccountID string
	HostedZone   map[string]*HostedZone
	VPC          *VPC
	Certificates map[string]string // domain:arn
	Github       *Github
	ArgoCD       *ArgoCD
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

// VPC contains state about the VPC
type VPC struct {
	VpcID   string
	CIDR    string
	Subnets map[string][]*VPCSubnet
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
	SecretKey *SecretKeySecret
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
	OauthApp     map[string]*GithubOauthApp
	Repositories map[string]*GithubRepository
}

// GithubRepository contains github repository data
type GithubRepository struct {
	Name      string
	FullName  string
	Types     []string
	GitURL    string
	DeployKey *DeployKey
}

// GithubOauthApp contains github oauth application data
type GithubOauthApp struct {
	Team         string
	Name         string
	SiteURL      string
	CallbackURL  string
	ClientID     string
	ClientSecret *ClientSecret
}

// ClientSecret contains state about
// an oauth app client secret
type ClientSecret struct {
	Name    string
	Path    string
	Version int64
}

// DeployKey contains github deploy key data
type DeployKey struct {
	Title            string
	ID               int64
	PublicKey        string
	PrivateKeySecret *PrivateKeySecret
}

// PrivateKeySecret contains information
// about a private key
type PrivateKeySecret struct {
	Name    string
	Path    string
	Version int64
}

// HostedZone contains information about the
// clusters hostedzone delegation
type HostedZone struct {
	IsDelegated bool
	IsCreated   bool
	Primary     bool
	Domain      string
	FQDN        string
	NameServers []string
}

// Validate the hostedzone
func (h *HostedZone) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Domain, validation.Required),
		validation.Field(&h.FQDN, validation.Required),
		validation.Field(&h.IsDelegated, validation.Required),
		validation.Field(&h.IsCreated, validation.Required),
	)
}

// NewUser returns repository data with defaults set
func NewRepository() *Repository {
	return &Repository{
		Name:      "",
		Region:    v1alpha1.RegionEuWest1,
		OutputDir: "infrastructure",
		Clusters:  map[string]*Cluster{},
	}
}
