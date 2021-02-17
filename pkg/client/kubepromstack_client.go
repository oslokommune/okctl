package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// KubePrometheusStack is the content of a kubernetes prometheus stack deployment
type KubePrometheusStack struct {
	CertificateARN         string
	Hostname               string
	AuthHostname           string
	ClientID               string
	SecretsConfigName      string
	SecretsCookieSecretKey string
	SecretsClientSecretKey string
	SecretsAdminUserKey    string
	SecretsAdminPassKey    string
	Chart                  *api.Helm
	Certificate            *api.Certificate
	IdentityPoolClient     *api.IdentityPoolClient
	ExternalSecret         *ExternalSecret
}

// CreateKubePrometheusStackOpts are the required inputs
type CreateKubePrometheusStackOpts struct {
	ID           api.ID
	Domain       string
	HostedZoneID string
	AuthDomain   string
	UserPoolID   string
}

// KubePrometheusStackService is an implementation of the business logic
type KubePrometheusStackService interface {
	CreateKubePrometheusStack(ctx context.Context, opts CreateKubePrometheusStackOpts) (*KubePrometheusStack, error)
}

// KubePrometheusStackAPI invokes REST API endpoints
type KubePrometheusStackAPI interface {
	CreateKubePrometheusStackHelmChart(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error)
}

// KubePrometheusStackStore is a storage layer implementation
type KubePrometheusStackStore interface {
	SaveKubePrometheusStack(stack *KubePrometheusStack) (*store.Report, error)
}

// KubePrometheusStackReport is a report layer
type KubePrometheusStackReport interface {
	ReportCreateKubePrometheusStack(stack *KubePrometheusStack, report *store.Report) error
}
