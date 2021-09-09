package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// AddPostgresToApplicationOpts defines required data to add a Postgres database to an application
type AddPostgresToApplicationOpts struct {
	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application

	DatabaseName string
}

// RemovePostgresFromApplicationOpts defines required data to remove a Postgres database from an application
type RemovePostgresFromApplicationOpts struct {
	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application

	DatabaseName string
}

// HasPostgresIntegrationOpts defines required data for checking if an integration is in place
type HasPostgresIntegrationOpts struct {
	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application

	DatabaseName string
}

// ApplicationPostgresService knows how to use the ApplicationPostgresService API
type ApplicationPostgresService interface {
	AddPostgresToApplication(ctx context.Context, opts AddPostgresToApplicationOpts) error
	RemovePostgresFromApplication(ctx context.Context, opts RemovePostgresFromApplicationOpts) error
	HasPostgresIntegration(ctx context.Context, opts HasPostgresIntegrationOpts) (bool, error)
}
