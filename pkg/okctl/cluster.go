// Package okctl provides a high level interface for all okctl commands
package okctl

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/oslokommune/okctl/pkg/api/core/cloud"
	"github.com/oslokommune/okctl/pkg/api/core/exe"
	"github.com/oslokommune/okctl/pkg/api/core/store"
)

// DeleteCluster removes an EKS cluster and all related resources
func (o *Okctl) DeleteCluster(opts *api.ClusterDeleteOpts) error {
	opts.RepositoryName = o.RepoData.Name

	eks := exe.NewClusterExe(o.BinariesProvider)
	clo := cloud.NewCluster(o.CloudProvider)
	sto := store.NewClusterStore(o.PersisterProvider)

	service := core.NewClusterService(sto, clo, eks)

	err := service.DeleteCluster(opts)
	if err != nil {
		return err
	}

	return nil
}

// CreateCluster starts the creation of all resources related to an EKS cluster
// such as a VPC, etc.
func (o *Okctl) CreateCluster(opts *api.ClusterCreateOpts) error {
	opts.RepositoryName = o.RepoData.Name
	opts.ClusterName = o.ClusterName(opts.Environment)

	eks := exe.NewClusterExe(o.BinariesProvider)
	clo := cloud.NewCluster(o.CloudProvider)
	sto := store.NewClusterStore(o.PersisterProvider)

	service := core.NewClusterService(sto, clo, eks)

	_, err := service.CreateCluster(opts)
	if err != nil {
		return nil
	}

	return err
}
