package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	// ClusterConfigKind maps up the resource Kind defined by eksctl
	ClusterConfigKind = "ClusterConfig"
	// ClusterConfigAPIVersion maps up the API Version we currently use towards eksctl
	ClusterConfigAPIVersion = "eksctl.io/v1alpha5"
)

// ClusterConfig is a recreation of:
// https://github.com/weaveworks/eksctl/blob/master/pkg/apis/eksctl.io/v1alpha5/types.go#L451
// where we have extract the parts that we are interested in for managing a eksctl cluster
type ClusterConfig struct {
	metav1.TypeMeta `json:",inline"`

	Metadata   ClusterMeta `json:"metadata"`
	IAM        ClusterIAM  `json:"iam"`
	VPC        ClusterVPC  `json:"vpc"`
	NodeGroups []NodeGroup `json:"nodeGroups"`
}

// ClusterMeta comes from eksctl and maps up what we need
type ClusterMeta struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

// ClusterIAM comes from eksctl and maps up what we need
type ClusterIAM struct {
	ServiceRolePermissionsBoundary             string `json:"serviceRolePermissionsBoundary"`
	FargatePodExecutionRolePermissionsBoundary string `json:"fargatePodExecutionRolePermissionsBoundary"`
	WithOIDC                                   bool   `json:"withOIDC"`
}

// ClusterVPC comes from eksctl and maps up what we need
type ClusterVPC struct {
	ID               string           `json:"id"`
	CIDR             string           `json:"cidr"`
	ClusterEndpoints ClusterEndpoints `json:"clusterEndpoints"`
	Subnets          ClusterSubnets   `json:"subnets"`
}

// ClusterEndpoints comes from eksctl and maps up what we need
type ClusterEndpoints struct {
	PrivateAccess bool `json:"privateAccess"`
	PublicAccess  bool `json:"publicAccess"`
}

// ClusterSubnets comes from eksctl and maps up what we need
type ClusterSubnets struct {
	Private map[string]ClusterNetwork `json:"private"`
	Public  map[string]ClusterNetwork `json:"public"`
}

// ClusterNetwork comes from eksctl and maps up what we need
type ClusterNetwork struct {
	ID   string `json:"id"`
	CIDR string `json:"cidr"`
}

// NodeGroup comes from eksctl and maps up what we need
type NodeGroup struct {
}

// ClusterConfigTypeMeta returns the defaults
func ClusterConfigTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ClusterConfigKind,
		APIVersion: ClusterConfigAPIVersion,
	}
}

// NewClusterConfig fills in all known default
// values
func NewClusterConfig() *ClusterConfig {
	return &ClusterConfig{
		TypeMeta: ClusterConfigTypeMeta(),
		IAM: ClusterIAM{
			WithOIDC: true,
		},
		VPC: ClusterVPC{
			ClusterEndpoints: ClusterEndpoints{
				PrivateAccess: true,
				PublicAccess:  true,
			},
			Subnets: ClusterSubnets{
				Private: map[string]ClusterNetwork{},
				Public:  map[string]ClusterNetwork{},
			},
		},
	}
}

// YAML returns a serializes version of the config
func (c *ClusterConfig) YAML() ([]byte, error) {
	return yaml.Marshal(c)
}

// CreateClusterConfigOpts defines the required inputs
// for creating a new cluster config
type CreateClusterConfigOpts struct {
	ClusterName  string
	Region       string
	Cidr         string
	AwsAccountID string
}

// Validate the cluster config options
func (o CreateClusterConfigOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.Cidr, validation.Required),
		validation.Field(&o.AwsAccountID, validation.Required),
	)
}

// ClusterConfigService defines the service interface for the cluster config
type ClusterConfigService interface {
	CreateClusterConfig(ctx context.Context, opts CreateClusterConfigOpts) (*ClusterConfig, error)
}

// ClusterConfigStore defines the storage operations
// for a cluster config
type ClusterConfigStore interface {
	SaveClusterConfig(*ClusterConfig) error
	DeleteClusterConfig(env string) error
	GetClusterConfig(env string) (*ClusterConfig, error)
}
