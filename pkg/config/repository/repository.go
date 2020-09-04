// Package repository knows how to interact with repository data
package repository

import (
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// Data stores the configured state of a repository
// as used by okctl
type Data struct {
	Name      string
	Region    string
	OutputDir string
	Clusters  []Cluster
}

// ClusterForEnv returns the cluster for the given environment
func (d *Data) ClusterForEnv(env string) *Cluster {
	for _, cluster := range d.Clusters {
		if cluster.Environment == env {
			return &cluster
		}
	}

	return nil
}

// SetClusterForEnv saves the cluster for the given environment
func (d *Data) SetClusterForEnv(cluster *Cluster, env string) {
	for i, c := range d.Clusters {
		if c.Environment == env {
			d.Clusters[i] = *cluster

			return
		}
	}

	d.Clusters = append(d.Clusters, *cluster)
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
	HostedZone   HostedZone
	AWS          AWS
	Certificates []Certificate
	Github       Github
	ArgoCD       ArgoCD
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
		validation.Field(&c.AWS),
		validation.Field(&c.Certificates),
	)
}

// ArgoCD contains information about the
// argocd setup
type ArgoCD struct {
	URL           string
	SecretKeyPath string
}

// Github contains information about the
// clusters configuration towards github
type Github struct {
	Organisation string
	Team         string
	DeployKey    DeployKey
	OauthApp     OauthApp
	Repository   Repository
}

// Repository contains github repository data
type Repository struct {
	Name   string
	GitURL string
}

// OauthApp contains github oauth application data
type OauthApp struct {
	Name             string
	ClientID         string
	ClientSecretPath string
}

// DeployKey contains github deploy key data
type DeployKey struct {
	Title string
	ID    int64
	Path  string
}

// HostedZone contains information about the
// clusters hostedzone delegation
type HostedZone struct {
	IsDelegated bool
	IsCreated   bool
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

// AWS represents the required information
type AWS struct {
	AccountID string
	Cidr      string
}

// Validate the data
func (a AWS) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.AccountID,
			validation.Required,
			validation.Match(regexp.MustCompile("^[0-9]{12}$")),
		),
	)
}

// Certificate represents a certificate
type Certificate struct {
	ARN    string
	Domain string
	FQDN   string
}

// Validate the certificate data
func (c Certificate) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ARN, validation.Required),
		validation.Field(&c.Domain, validation.Required),
		validation.Field(&c.FQDN, validation.Required),
	)
}

// New returns repository data with defaults set
func New() *Data {
	return &Data{
		Region:    v1alpha1.RegionEuWest1,
		OutputDir: "infrastructure",
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
