// Package api provides the domain model for okctl
package api

import (
	"context"
	"regexp"

	val "github.com/go-ozzo/ozzo-validation/v4"
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
	Environment    string
	AWSAccountID   string
	Cidr           string
	ClusterName    string
	RepositoryName string
	Region         string
	Config         *ClusterConfig
}

// ClusterCreateOpts specifies the required inputs for creating a cluster
type ClusterCreateOpts struct {
	Environment    string
	AWSAccountID   string
	Cidr           string
	RepositoryName string
	Region         string
	ClusterName    string

	VpcID             string
	VpcPrivateSubnets []VpcSubnet
	VpcPublicSubnets  []VpcSubnet
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
		val.Field(&o.Cidr, val.Required),
		val.Field(&o.Region, val.Required),
		val.Field(&o.VpcID, val.Required),
		val.Field(&o.VpcPrivateSubnets, val.Required),
		val.Field(&o.VpcPublicSubnets, val.Required),
	)
}

// ClusterDeleteOpts specifies the required inputs for deleting a cluster
type ClusterDeleteOpts struct {
	Environment    string
	RepositoryName string
	ClusterName    string
}

// Validate the delete inputs
func (o *ClusterDeleteOpts) Validate() error {
	return val.ValidateStruct(o,
		val.Field(&o.Environment,
			val.Required,
			val.Length(envMinLength, envMaxLength),
		),
		val.Field(&o.RepositoryName, val.Required),
		val.Field(&o.ClusterName, val.Required),
	)
}

// ClusterService provides an interface for the business logic when working with clusters
type ClusterService interface {
	CreateCluster(context.Context, ClusterCreateOpts) (*Cluster, error)
	DeleteCluster(context.Context, ClusterDeleteOpts) error
}

// ClusterRun provides an interface for running CLIs
type ClusterRun interface {
	CreateCluster(string, *ClusterConfig) error
	DeleteCluster(string) error
}

// ClusterStore provides an interface for storage operations
type ClusterStore interface {
	SaveCluster(*Cluster) error
	DeleteCluster(string) error
	GetCluster(string) (*Cluster, error)
}
