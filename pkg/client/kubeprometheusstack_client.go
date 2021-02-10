package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// KubePrometheusStack is the content of an external-secrets deployment
type KubePrometheusStack struct {
	Policy         *api.ManagedPolicy
	ServiceAccount *api.ServiceAccount
	Chart          *api.Helm
}

// CreateKubePrometheusStackOpts contains the required inputs
type CreateKubePrometheusStackOpts struct {
	ID api.ID
}

// KubePrometheusStackService is an implementation of the business logic
type KubePrometheusStackService interface {
	CreateKubePrometheusStack(ctx context.Context, opts CreateKubePrometheusStackOpts) (*KubePrometheusStack, error)
	//DeleteKubePrometheusStack(ctx context.Context, id api.ID) error
}

// KubePrometheusStackAPI invokes REST API endpoints
type KubePrometheusStackAPI interface {
	//CreateKubePrometheusStackPolicy(opts api.CreateKubePrometheusStackPolicyOpts) (*api.ManagedPolicy, error)
	//DeleteKubePrometheusStackPolicy(id api.ID) error
	//CreateKubePrometheusStackServiceAccount(opts api.CreateKubePrometheusStackServiceAccountOpts) (*api.ServiceAccount, error)
	//DeleteKubePrometheusStackServiceAccount(id api.ID) error
	CreateKubePrometheusStackHelmChart(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error)
}

// KubePrometheusStackStore is a storage layer implementation
type KubePrometheusStackStore interface {
	SaveKubePrometheusStack(externalSecrets *KubePrometheusStack) (*store.Report, error)
	RemoveKubePrometheusStack(id api.ID) (*store.Report, error)
}

// KubePrometheusStackReport is a report layer
type KubePrometheusStackReport interface {
	ReportCreateKubePrometheusStack(secret *KubePrometheusStack, report *store.Report) error
	ReportDeleteKubePrometheusStack(report *store.Report) error
}
