package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// Autoscaler is the content of an autoscaler deployment
type Autoscaler struct {
	Policy         *api.ManagedPolicy
	ServiceAccount *api.ServiceAccount
	Chart          *api.Helm
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

// AutoscalerAPI invokes REST API endpoints
type AutoscalerAPI interface {
	CreateAutoscalerPolicy(opts api.CreateAutoscalerPolicy) (*api.ManagedPolicy, error)
	DeleteAutoscalerPolicy(id api.ID) error
	CreateAutoscalerServiceAccount(opts api.CreateAutoscalerServiceAccountOpts) (*api.ServiceAccount, error)
	DeleteAutoscalerServiceAccount(id api.ID) error
	CreateAutoscalerHelmChart(opts api.CreateAutoscalerHelmChartOpts) (*api.Helm, error)
}

// AutoscalerStore is a storage layer implementation
type AutoscalerStore interface {
	SaveAutoscaler(scaler *Autoscaler) (*store.Report, error)
	RemoveAutoscaler(id api.ID) (*store.Report, error)
}

// AutoscalerReport is a report layer
type AutoscalerReport interface {
	ReportCreateAutoscaler(scaler *Autoscaler, report *store.Report) error
	ReportDeleteAutoscaler(report *store.Report) error
}
