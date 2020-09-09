// Package repository knows how to interact with repository data
package repository

import (
	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
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

// Data stores the configured state of a repository
// as used by okctl
type Data struct {
	Name      string
	Region    string
	OutputDir string
	Clusters  map[string]*Cluster
}

// ClusterForEnv returns the cluster for the given environment
func (d *Data) ClusterForEnv(env string) *Cluster {
	if c, ok := d.Clusters[env]; ok {
		return c
	}

	return nil
}

// Validate the provided data
func (d *Data) Validate() error {
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
	OauthApp     map[string]*OauthApp
	Repositories map[string]*Repository
}

// Repository contains github repository data
type Repository struct {
	Name      string
	FullName  string
	Types     []string
	GitURL    string
	DeployKey *DeployKey
}

// OauthApp contains github oauth application data
type OauthApp struct {
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

// New returns repository data with defaults set
func New() *Data {
	return &Data{
		Name:      "",
		Region:    v1alpha1.RegionEuWest1,
		OutputDir: "infrastructure",
		Clusters:  map[string]*Cluster{},
	}
}

// Survey starts an interactive survey that queries
// the user for input
func (d *Data) Survey() (*Data, error) {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Name:",
				Help:    "A descriptive name, e.g., team or project, used among other things to prefix AWS resources",
			},
		},
		{
			Name: "region",
			Prompt: &survey.Select{
				Message: "Choose AWS region:",
				Options: v1alpha1.SupportedRegions(),
				Help:    "The AWS region resources will be created in",
			},
		},
		{
			Name: "basedir",
			Prompt: &survey.Input{
				Message: "Output directory:",
				Default: "infrastructure",
				Help:    "Path in the repository where generated files are stored",
			},
		},
	}

	answers := struct {
		Name    string
		Region  string
		BaseDir string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}

	d.Name = answers.Name
	d.Region = answers.Region
	d.OutputDir = answers.BaseDir

	return d, errors.Wrap(d.Validate(), "failed to validate repository data")
}

// YAML returns the state of the data object in YAML
func (d *Data) YAML() ([]byte, error) {
	return yaml.Marshal(d)
}
