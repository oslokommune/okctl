// Package api provides the domain model for okctl
package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Cluster contains the core state for a cluster
type Cluster struct {
	ID     ID
	Cidr   string
	Config *ClusterConfig
}

// ClusterCreateOpts specifies the required inputs for creating a cluster
type ClusterCreateOpts struct {
	ID                ID
	Cidr              string
	VpcID             string
	VpcPrivateSubnets []VpcSubnet
	VpcPublicSubnets  []VpcSubnet
}

// Validate the create inputs
func (o *ClusterCreateOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Cidr, validation.Required),
		validation.Field(&o.VpcID, validation.Required),
		validation.Field(&o.VpcPrivateSubnets, validation.Required),
		validation.Field(&o.VpcPublicSubnets, validation.Required),
	)
}

// ClusterDeleteOpts specifies the required inputs for deleting a cluster
type ClusterDeleteOpts struct {
	ID ID
}

// Validate the delete inputs
func (o *ClusterDeleteOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
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
