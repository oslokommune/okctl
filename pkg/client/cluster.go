package client

import (
	"github.com/oslokommune/okctl/pkg/api"
)

// We are shadowing the api interfaces for now, but
// this is probably not sustainable.

// ClusterService orchestrates the creation of a cluster
type ClusterService interface {
	api.ClusterService
}

// ClusterAPI invokes the API calls for creating a cluster
type ClusterAPI interface {
	api.ClusterRun
}

// ClusterStore stores the data
type ClusterStore interface {
	api.ClusterStore
}
