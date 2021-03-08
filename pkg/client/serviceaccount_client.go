package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// ServiceAccountService implements the business logic
type ServiceAccountService interface {
	CreateServiceAccount(ctx context.Context, opts api.CreateServiceAccountOpts) (*api.ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, opts api.DeleteServiceAccountOpts) error
}

// ServiceAccountAPI invokes the remote API
type ServiceAccountAPI interface {
	CreateServiceAccount(opts api.CreateServiceAccountOpts) (*api.ServiceAccount, error)
	DeleteServiceAccount(opts api.DeleteServiceAccountOpts) error
}

// ServiceAccountStore provides a persistence layer
type ServiceAccountStore interface {
	SaveCreateServiceAccount(account *api.ServiceAccount) (*store.Report, error)
	RemoveDeleteServiceAccount(name string) (*store.Report, error)
}

// ServiceAccountReports provides output on the result
type ServiceAccountReport interface {
	ReportCreateServiceAccount(account *api.ServiceAccount, report *store.Report) error
	ReportDeleteServiceAccount(name string, report *store.Report) error
}
