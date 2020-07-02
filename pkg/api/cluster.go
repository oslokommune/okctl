// Package api provides the domain model for okctl
package api

import (
	"context"
	"regexp"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

const (
	envMinLength     = 3
	envMaxLength     = 100
	repoMinLength    = 3
	repoMaxLength    = 100
	clusterMinLength = 3
	clusterMaxLength = 100
)

// Cluster contains the core state for a cluster
type Cluster struct {
	Environment  string
	AWSAccountID string
	Cidr         string
	Config       *v1alpha1.ClusterConfig
}

// ClusterCreateOpts specifies the required inputs for creating a cluster
type ClusterCreateOpts struct {
	Environment    string
	AWSAccountID   string
	Cidr           string
	RepositoryName string
	Region         string
	ClusterName    string
}

// Validate the create inputs
func (o *ClusterCreateOpts) Validate() error {
	return val.ValidateStruct(o,
		val.Field(&o.AWSAccountID,
			val.Required,
			val.Match(regexp.MustCompile("^[0-9]{12}$")).
				Error("must consist of 12 digits"),
		),
		val.Field(&o.RepositoryName,
			val.Required,
			val.Length(repoMinLength, repoMaxLength),
		),
		val.Field(&o.ClusterName,
			val.Required,
			val.Length(clusterMinLength, clusterMaxLength),
		),
		val.Field(&o.Environment,
			val.Required,
			val.Length(envMinLength, envMaxLength),
		),
		val.Field(&o.Cidr,
			val.Required,
		),
		val.Field(&o.Region,
			val.Required,
		),
	)
}

// ClusterDeleteOpts specifies the required inputs for deleting a cluster
type ClusterDeleteOpts struct {
	Environment    string
	RepositoryName string
}

// Validate the delete inputs
func (o *ClusterDeleteOpts) Validate() error {
	return val.ValidateStruct(o,
		val.Field(&o.Environment,
			val.Required,
			val.Length(envMinLength, envMaxLength),
		),
		val.Field(&o.RepositoryName,
			val.Required,
		),
	)
}

// ClusterService provides an interface for the business logic when working with clusters
type ClusterService interface {
	CreateCluster(context.Context, ClusterCreateOpts) (*Cluster, error)
	DeleteCluster(context.Context, ClusterDeleteOpts) error
}

// ClusterExe provides an interface for running CLIs
type ClusterExe interface {
	CreateCluster(*v1alpha1.ClusterConfig) error
	DeleteCluster(*v1alpha1.ClusterConfig) error
}

// ClusterCloud provides an interface for running actions towards the cloud
type ClusterCloud interface {
	CreateCluster(awsAccountID, clusterName, env, repoName, cidr, region string) (*v1alpha1.ClusterConfig, error)
	DeleteCluster(env, repoName string) error
}

// ClusterStore provides an interface for storag operations
type ClusterStore interface {
	SaveCluster(*Cluster) error
	DeleteCluster(env string) error
	GetCluster(env string) (*Cluster, error)
}
