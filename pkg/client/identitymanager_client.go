package client

import (
	"context"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// CreateIdentityPoolUserOpts contains the required inputs
type CreateIdentityPoolUserOpts struct {
	ID         api.ID
	Email      string
	UserPoolID string
}

// nolint: lll
const emailRx = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"

// Validate the inputs
func (o CreateIdentityPoolUserOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.UserPoolID, validation.Required),
		validation.Field(&o.Email, validation.Required, validation.Match(regexp.MustCompile(emailRx)).Error("must be valid email")),
	)
}

// IdentityManagerService orchestrates the creation of an identity pool
type IdentityManagerService interface {
	CreateIdentityPool(ctx context.Context, opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error)
	CreateIdentityPoolClient(ctx context.Context, opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error)
	CreateIdentityPoolUser(ctx context.Context, opts CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error)
	DeleteIdentityPool(ctx context.Context, opts api.ID) error
	DeleteIdentityPoolClient(ctx context.Context, opts api.DeleteIdentityPoolClientOpts) error
}

// IdentityManagerAPI invokes the API calls for creating an identity pool
type IdentityManagerAPI interface {
	CreateIdentityPool(opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error)
	CreateIdentityPoolClient(opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error)
	CreateIdentityPoolUser(opts api.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error)
	DeleteIdentityPool(opts api.DeleteIdentityPoolOpts) error
	DeleteIdentityPoolClient(opts api.DeleteIdentityPoolClientOpts) error
}

// IdentityManagerStore stores the data
type IdentityManagerStore interface {
	SaveIdentityPool(pool *api.IdentityPool) (*store.Report, error)
	SaveIdentityPoolClient(client *api.IdentityPoolClient) (*store.Report, error)
	SaveIdentityPoolUser(user *api.IdentityPoolUser) (*store.Report, error)
	RemoveIdentityPool(id api.ID) (*store.Report, error)
	RemoveIdentityPoolClient(opts api.DeleteIdentityPoolClientOpts) (*store.Report, error)
}

// IdentityManagerState implements the state layer
type IdentityManagerState interface {
	SaveIdentityPool(pool *api.IdentityPool) (*store.Report, error)
	SaveIdentityPoolClient(client *api.IdentityPoolClient) (*store.Report, error)
	SaveIdentityPoolUser(user *api.IdentityPoolUser) (*store.Report, error)
	GetIdentityPool() state.IdentityPool
	RemoveIdentityPool(id api.ID) (*store.Report, error)
	RemoveIdentityPoolClient(opts api.DeleteIdentityPoolClientOpts) (*store.Report, error)
}

// IdentityManagerReport provides output of the actions
type IdentityManagerReport interface {
	ReportIdentityPool(pool *api.IdentityPool, reports []*store.Report) error
	ReportIdentityPoolClient(client *api.IdentityPoolClient, reports []*store.Report) error
	ReportIdentityPoolUser(client *api.IdentityPoolUser, reports []*store.Report) error
	ReportDeleteIdentityPool(reports []*store.Report) error
	ReportDeleteIdentityPoolClient(reports []*store.Report) error
}
