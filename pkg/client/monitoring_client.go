package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// KubePromStack is the content of a kubernetes prometheus stack deployment
type KubePromStack struct {
	ID                                api.ID
	CertificateARN                    string
	Hostname                          string
	AuthHostname                      string
	ClientID                          string
	SecretsConfigName                 string
	SecretsCookieSecretKey            string
	SecretsClientSecretKey            string
	SecretsAdminUserKey               string
	SecretsAdminPassKey               string
	Chart                             *api.Helm
	Certificate                       *Certificate
	IdentityPoolClient                *IdentityPoolClient
	ExternalSecret                    *KubernetesManifest
	FargateCloudWatchPolicyARN        string
	FargateProfilePodExecutionRoleARN string
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

// Loki is the content of a kubernetes loki deployment
type Loki struct {
	ID    api.ID
	Chart *api.Helm
}

// CreateLokiOpts are the required inputs
type CreateLokiOpts struct {
	ID api.ID
}

// DeleteLokiOpts are the required inputs
type DeleteLokiOpts struct {
	ID api.ID
}

// Promtail is the content of a kubernetes promtail deployment
type Promtail struct {
	ID    api.ID
	Chart *api.Helm
}

// CreatePromtailOpts are the required inputs
type CreatePromtailOpts struct {
	ID api.ID
}

// DeletePromtailOpts are the required inputs
type DeletePromtailOpts struct {
	ID api.ID
}

// Tempo is the content of a kubernetes tempo deployment
type Tempo struct {
	ID    api.ID
	Chart *api.Helm
}

// CreateTempoOpts are the required inputs
type CreateTempoOpts struct {
	ID api.ID
}

// DeleteTempoOpts are the required inputs
type DeleteTempoOpts struct {
	ID api.ID
}

// MonitoringService is an implementation of the business logic
type MonitoringService interface {
	CreateKubePromStack(ctx context.Context, opts CreateKubePromStackOpts) (*KubePromStack, error)
	DeleteKubePromStack(ctx context.Context, opts DeleteKubePromStackOpts) error
	CreateLoki(ctx context.Context, opts CreateLokiOpts) (*Loki, error)
	DeleteLoki(ctx context.Context, opts DeleteLokiOpts) error
	CreatePromtail(ctx context.Context, opts CreatePromtailOpts) (*Promtail, error)
	DeletePromtail(ctx context.Context, opts DeletePromtailOpts) error
	CreateTempo(ctx context.Context, opts CreateTempoOpts) (*Tempo, error)
	DeleteTempo(ctx context.Context, opts DeleteTempoOpts) error
}

// MonitoringAPI invokes REST API endpoints
type MonitoringAPI interface {
	CreateKubePromStack(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error)
	DeleteKubePromStack(opts api.DeleteHelmReleaseOpts) error
	CreateLoki(opts CreateLokiOpts) (*api.Helm, error)
	DeleteLoki(opts api.DeleteHelmReleaseOpts) error
	CreatePromtail(opts CreatePromtailOpts) (*api.Helm, error)
	DeletePromtail(opts api.DeleteHelmReleaseOpts) error
	CreateTempo(opts api.CreateHelmReleaseOpts) (*api.Helm, error)
	DeleteTempo(opts api.DeleteHelmReleaseOpts) error
}

// MonitoringStore is a storage layer implementation
type MonitoringStore interface {
	SaveKubePromStack(stack *KubePromStack) (*store.Report, error)
	RemoveKubePromStack(id api.ID) (*store.Report, error)
	SaveLoki(loki *Loki) (*store.Report, error)
	RemoveLoki(id api.ID) (*store.Report, error)
	SavePromtail(Promtail *Promtail) (*store.Report, error)
	RemovePromtail(id api.ID) (*store.Report, error)
	SaveTempo(Tempo *Tempo) (*store.Report, error)
	RemoveTempo(id api.ID) (*store.Report, error)
}

// MonitoringState is a state layer implementation
type MonitoringState interface {
	SaveKubePromStack(stack *KubePromStack) (*store.Report, error)
	RemoveKubePromStack(id api.ID) (*store.Report, error)
	GetKubePromStack() (*KubePromStack, error)
}

// MonitoringReport is a report layer
type MonitoringReport interface {
	ReportSaveKubePromStack(stack *KubePromStack, reports []*store.Report) error
	ReportRemoveKubePromStack(reports []*store.Report) error
	ReportSaveLoki(loki *Loki, report *store.Report) error
	ReportRemoveLoki(report *store.Report) error
	ReportSavePromtail(Promtail *Promtail, report *store.Report) error
	ReportRemovePromtail(report *store.Report) error
	ReportSaveTempo(Tempo *Tempo, report *store.Report) error
	ReportRemoveTempo(report *store.Report) error
}
