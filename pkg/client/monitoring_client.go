package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// KubePromStack is the content of a kubernetes prometheus stack deployment
type KubePromStack struct {
	ID                                api.ID
	AuthHostname                      string
	CertificateARN                    string
	ClientID                          string
	FargateCloudWatchPolicyARN        string
	FargateProfilePodExecutionRoleARN string
	Hostname                          string
	SecretsAdminPassKey               string
	SecretsAdminUserKey               string
	SecretsClientSecretKey            string
	SecretsConfigName                 string
	SecretsCookieSecretKey            string
	Certificate                       *Certificate
	Chart                             *Helm
	ExternalSecret                    *KubernetesManifest
	IdentityPoolClient                *IdentityPoolClient
}

// CreateKubePromStackOpts are the required inputs
type CreateKubePromStackOpts struct {
	ID           api.ID
	Domain       string
	HostedZoneID string
	AuthDomain   string
	UserPoolID   string
}

// DeleteKubePromStackOpts are the required inputs
type DeleteKubePromStackOpts struct {
	ID     api.ID
	Domain string
}

// MonitoringService is an implementation of the business logic
type MonitoringService interface {
	CreateKubePromStack(ctx context.Context, opts CreateKubePromStackOpts) (*KubePromStack, error)
	DeleteKubePromStack(ctx context.Context, opts DeleteKubePromStackOpts) error
	CreateLoki(ctx context.Context, id api.ID) (*Helm, error)
	DeleteLoki(ctx context.Context, id api.ID) error
	CreatePromtail(ctx context.Context, id api.ID) (*Helm, error)
	DeletePromtail(ctx context.Context, id api.ID) error
	CreateTempo(ctx context.Context, id api.ID) (*Helm, error)
	DeleteTempo(ctx context.Context, id api.ID) error
}

// MonitoringState is a state layer implementation
type MonitoringState interface {
	SaveKubePromStack(stack *KubePromStack) error
	RemoveKubePromStack() error
	GetKubePromStack() (*KubePromStack, error)
	HasKubePromStack() (bool, error)
}
