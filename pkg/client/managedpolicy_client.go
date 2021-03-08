package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// ManagedPolicyService implements the business logic
type ManagedPolicyService interface {
	CreatePolicy(ctx context.Context, opts api.CreatePolicyOpts) (*api.ManagedPolicy, error)
	DeletePolicy(ctx context.Context, opts api.DeletePolicyOpts) error
}

// ManagedPolicyAPI invokes the remote API
type ManagedPolicyAPI interface {
	CreatePolicy(opts api.CreatePolicyOpts) (*api.ManagedPolicy, error)
	DeletePolicy(opts api.DeletePolicyOpts) error
}

// ManagedPolicyStore provides a persistence layer
type ManagedPolicyStore interface {
	SaveCreatePolicy(policy *api.ManagedPolicy) (*store.Report, error)
	RemoveDeletePolicy(stackName string) (*store.Report, error)
}

// ManagedPolicyReports provides output on the result
type ManagedPolicyReport interface {
	ReportCreatePolicy(policy *api.ManagedPolicy, report *store.Report) error
	ReportDeletePolicy(stackName string, report *store.Report) error
}
