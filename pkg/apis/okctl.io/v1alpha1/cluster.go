package v1alpha1

import (
	"fmt"

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
	AccountID int `json:"accountID"`
}

// String returns a unique identifier for a cluster
// Not sure about this..
func (c *ClusterMeta) String() string {
	return fmt.Sprintf("%s-%s.%s.okctl.io/%d", c.Name, c.Environment, c.Region, c.AccountID)
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

// ClusterIntegrations ...
type ClusterIntegrations struct {
	// ALBIngressController if set to true will install the AWS ALB Ingress Controller
	// into the cluster
	// +optional
	ALBIngressController bool `json:"albIngressController,omitempty"`

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

// ClusterTypeMeta returns an initialised TypeMeta object
// for a Cluster
func ClusterTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ClusterKind,
		APIVersion: ClusterAPIVersion,
	}
}

// NewDefaultCluster returns a cluster definition with sensible defaults
func NewDefaultCluster(name, env, org, repo, team string, accountID int) Cluster {
	return Cluster{
		TypeMeta: ClusterTypeMeta(),
		Metadata: ClusterMeta{
			Name:        name,
			Environment: env,
			Region:      "eu-west-1",
			AccountID:   accountID,
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
			ALBIngressController: true,
			ExternalDNS:          true,
			ExternalSecrets:      true,
			Cognito:              true,
			ArgoCD:               true,
		},
		DNSZones: []ClusterDNSZone{
			{
				ParentDomain:  fmt.Sprintf("%s-%s.oslo.systems", name, env),
				ReuseExisting: false,
			},
		},
	}
}
