package v1alpha1

import (
	"errors"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ClusterKind is a string value that represents the resource type
	ClusterKind = "Cluster"
	// ClusterAPIVersion defines the versioned schema of this representation
	ClusterAPIVersion = "okctl.io/v1alpha1"
)

// Cluster is a unique Kubernetes cluster with a set of integrations that
// can be enabled or disabled.
type Cluster struct {
	metav1.TypeMeta `json:",inline"`

	// Metadata uniquely identifies a cluster.
	Metadata ClusterMeta `json:"metadata"`

	// Github defines what organisation, repository, team, etc. that
	// this cluster will integrate with.
	Github ClusterGithub `json:"github"`

	// PrimaryDNSZone defines the main primary zone to associate with this
	// cluster. This will be the zone that we will use to create domains
	// for auth, ArgoCD, etc.
	PrimaryDNSZone ClusterDNSZone `json:"primaryDNSZone"`

	// VPC defines how we configure the VPC for the cluster
	// +optional
	VPC *ClusterVPC `json:"vpc,omitempty"`

	// Integrations defines what cluster integrations we deploy to the
	// cluster
	// +optional
	Integrations *ClusterIntegrations `json:"integrations,omitempty"`

	// DNSZones is an optional list of DNS zones managed or associated with
	// this cluster.
	// +optional
	DNSZones []ClusterDNSZone `json:"dnsZones,omitempty"`
}

// Validate calls each members Validate function
func (c Cluster) Validate() error {
	result := validation.ValidateStruct(&c,
		validation.Field(&c.Metadata),
		validation.Field(&c.Github),
		validation.Field(&c.PrimaryDNSZone),
		validation.Field(&c.VPC),
		validation.Field(&c.Integrations),
	)

	return result
}

// ClusterMeta describes a unique cluster
type ClusterMeta struct {
	// Name is a descriptive value given to the cluster, e.g., the name
	// of the team, product, project, etc.
	Name string `json:"name"`

	// Environment defines the purpose of the cluster, e.g., testing,
	// staging, production.
	Environment string `json:"environment"`

	// Region specifies the AWS region the cluster should be created in
	// https://aws.amazon.com/about-aws/global-infrastructure/regions_az/
	Region string `json:"region"`

	// AccountID specifies the AWS Account ID
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/console_account-alias.html
	AccountID string `json:"accountID"`
}

// Validate ensures ClusterMeta contains the right information
func (receiver ClusterMeta) Validate() error {
	return validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Name, validation.Required),
		validation.Field(&receiver.Environment,
			validation.Required,
			validation.Match(regexp.MustCompile("^[a-zA-Z]{3,64}$")).Error("must consist of 3-64 characters (a-z, A-Z)")),
		validation.Field(&receiver.Region, validation.Required, validation.In("eu-west-1").Error("for now, only \"eu-west-1\" is supported")),
		validation.Field(&receiver.AccountID, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{12}$")).Error("must consist of 12 digits")),
	)
}

// String returns a unique identifier for a cluster
// Not sure about this..
func (receiver *ClusterMeta) String() string {
	return fmt.Sprintf("%s-%s.%s.okctl.io/%s", receiver.Name, receiver.Environment, receiver.Region, receiver.AccountID)
}

// ClusterVPC is a definition of the VPC we create for the EKS cluster
type ClusterVPC struct {
	// CIDR is the IP-address range to associate with the VPC
	// https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing.
	// The VPC CIDR must be compatible with EKS: https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html
	// +optional
	CIDR string `json:"cidr,omitempty"`

	// HighAvailability means we create redundancy in the network setup. If set to
	// true we will create a NAT gateway per public subnet, instead of routing
	// all traffic through one.
	// +optional
	HighAvailability bool `json:"highAvailability,omitempty"`
}

// Validate ensures ClusterVPC contains the right information
func (c ClusterVPC) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.CIDR, validation.Required),
		validation.Field(&c.HighAvailability, validation.Required),
	)
}

// ClusterDNSZone is analogous to a DNS Zone file (https://en.wikipedia.org/wiki/Zone_file).
// A DNS Zone represents a subset, in form of a single parent domain, of the hierarchical
// domain name structure. In AWS, we map this data to a Route53 HostedZone.
type ClusterDNSZone struct {
	// ParentDomain is the root domain for all DNS records of this
	// DNS zone delegation, e.g., `{team-name}.oslo.systems`
	ParentDomain string `json:"parentDomain"`

	// ReuseExisting determines if we should look for an existing DNS zone
	// or create a new one. If set to true, we will not attempt to create a
	// new DNS zone.
	ReuseExisting bool `json:"managedZone"`
}

// Validate ensures ClusterDNSZone contains necessary and correct information
func (c ClusterDNSZone) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ParentDomain, validation.Required, is.Domain),
	)
}

// ClusterGithub identifies a repository and path on github.com where
// we can set up an integration with Argo CD, among other things.
type ClusterGithub struct {
	// Organisation name on github.com, e.g., "oslokommune"
	Organisation string `json:"organisation"`

	// Repository name on github.com, e.g., "okctl". The repository
	// you specify here must be owned by the organisation specified above.
	Repository string `json:"repository"`

	// OutputPath is a path from the root of the org/repository where
	// we can store generated output files
	OutputPath string `json:"outputPath"`

	// Team name on github.com, e.g., "kjøremiljø". The team you
	// specify here must be owned by the organisation specified above.
	Team string `json:"team"`
}

// Validate returns an error if ClusterGithub is missing required information
func (c ClusterGithub) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Organisation, validation.Required),
		validation.Field(&c.Repository, validation.Required),
		validation.Field(&c.OutputPath, validation.Required),
		validation.Field(&c.Team, validation.Required),
	)
}

// ClusterIntegrations ...
type ClusterIntegrations struct {
	// ALBIngressController if set to true will install the AWS ALB Ingress Controller
	// into the cluster
	// +optional
	ALBIngressController bool `json:"albIngressController,omitempty"`

	// AWSLoadBalancerController if set to true will install the AWS load balancer controller
	// +optional
	AWSLoadBalancerController bool `json:"awsLoadBalancerController"`

	// ExternalDNS if set to true will install the external-dns controller into the cluster
	// +optional
	ExternalDNS bool `json:"externalDNS,omitempty"`

	// ExternalSecrets if set to true will install the external-secrets controller into the cluster
	// +optional
	ExternalSecrets bool `json:"externalSecrets,omitempty"`

	// Cognito if set to true will install the Cognito user pool into the cluster.
	// Might want to make this one more fine-grained, so that the teams can more easily
	// give access to their admin APIs or whatever. Might not be required for now.
	// +optional
	Cognito bool `json:"cognito,omitempty"`

	// ArgoCD if set to true will install the ArgoCD deployment setup into the cluster. This
	// integration requires ALBIngressController, ExternalDNS and Cognito.
	// +optional
	ArgoCD bool `json:"argoCD,omitempty"`
}

// Validate ensures there is no conflicting options
func (c ClusterIntegrations) Validate() error {
	if c.ArgoCD && !c.Cognito {
		return errors.New("the identity provider cognito is required when using ArgoCD")
	}

	return nil
}

// ClusterTypeMeta returns an initialised TypeMeta object
// for a Cluster
func ClusterTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ClusterKind,
		APIVersion: ClusterAPIVersion,
	}
}

// NewDefaultCluster returns a cluster definition with sensible defaults
func NewDefaultCluster(name, env, org, repo, team, accountID string) Cluster {
	return Cluster{
		TypeMeta: ClusterTypeMeta(),
		Metadata: ClusterMeta{
			Name:        name,
			Environment: env,
			Region:      "eu-west-1",
			AccountID:   accountID,
		},
		PrimaryDNSZone: ClusterDNSZone{
			ParentDomain:  fmt.Sprintf("%s-%s.oslo.systems", name, env),
			ReuseExisting: false,
		},
		Github: ClusterGithub{
			Organisation: org,
			Repository:   repo,
			OutputPath:   "infrastructure",
			Team:         team,
		},
		VPC: &ClusterVPC{
			CIDR:             "192.168.0.0/20",
			HighAvailability: true,
		},
		Integrations: &ClusterIntegrations{
			ALBIngressController:      false,
			AWSLoadBalancerController: true,
			ExternalDNS:               true,
			ExternalSecrets:           true,
			Cognito:                   true,
			ArgoCD:                    true,
		},
	}
}
