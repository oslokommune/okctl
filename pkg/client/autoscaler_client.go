package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// Autoscaler is the content of an autoscaler deployment
type Autoscaler struct {
	Policy         *ManagedPolicy
	ServiceAccount *ServiceAccount
	Chart          *Helm
}

// CreateAutoscalerOpts contains the required inputs
type CreateAutoscalerOpts struct {
	ID api.ID
}

// AutoscalerService is an implementation of the business logic
type AutoscalerService interface {
	CreateAutoscaler(ctx context.Context, opts CreateAutoscalerOpts) (*Autoscaler, error)
	DeleteAutoscaler(ctx context.Context, id api.ID) error
}

// AutoscalerState knows how to store and retrieve information about the Autoscaler
type AutoscalerState interface {
	HasAutoscaler() (bool, error)
}
