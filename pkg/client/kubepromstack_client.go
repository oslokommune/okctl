package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// KubePromStack is the content of a kubernetes prometheus stack deployment
type KubePromStack struct {
	ID                     api.ID
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

// CreateKubePromStackOpts are the required inputs
type CreateKubePromStackOpts struct {
	ID           api.ID
	Domain       string
	HostedZoneID string
	AuthDomain   string
	UserPoolID   string
}

// KubePromStackService is an implementation of the business logic
type KubePromStackService interface {
	CreateKubePromStack(ctx context.Context, opts CreateKubePromStackOpts) (*KubePromStack, error)
}

// KubePromStackAPI invokes REST API endpoints
type KubePromStackAPI interface {
	CreateKubePromStack(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error)
}

// KubePromStackStore is a storage layer implementation
type KubePromStackStore interface {
	SaveKubePromStack(stack *KubePromStack) (*store.Report, error)
}

// KubePromStackState is a state layer implementation
type KubePromStackState interface {
	SaveKubePromStack(stack *KubePromStack) (*store.Report, error)
}

// KubePromStackReport is a report layer
type KubePromStackReport interface {
	ReportKubePromStack(stack *KubePromStack, reports []*store.Report) error
}
