package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type identityManagerService struct {
	out io.Writer
}

func (i identityManagerService) CreateIdentityPool(_ context.Context, _ client.CreateIdentityPoolOpts) (*client.IdentityPool, error) {
	fmt.Fprintf(i.out, formatCreate("Identity Pool"))

	return &client.IdentityPool{}, nil
}

func (i identityManagerService) DeleteIdentityPool(_ context.Context, _ api.ID) error {
	fmt.Fprintf(i.out, formatDelete("Identity Pool"))

	return nil
}

func (i identityManagerService) CreateIdentityPoolClient(_ context.Context, _ client.CreateIdentityPoolClientOpts) (*client.IdentityPoolClient, error) {
	fmt.Fprintf(i.out, formatCreate("Identity Pool Client"))

	return &client.IdentityPoolClient{}, nil
}

func (i identityManagerService) DeleteIdentityPoolClient(_ context.Context, _ client.DeleteIdentityPoolClientOpts) error {
	fmt.Fprintf(i.out, formatDelete("Identity Pool Client"))

	return nil
}

func (i identityManagerService) CreateIdentityPoolUser(_ context.Context, opts client.CreateIdentityPoolUserOpts) (*client.IdentityPoolUser, error) {
	fmt.Fprintf(i.out, formatCreate(fmt.Sprintf("Identity Pool User %s", opts.Email)))

	return &client.IdentityPoolUser{}, nil
}

func (i identityManagerService) DeleteIdentityPoolUser(_ context.Context, opts client.DeleteIdentityPoolUserOpts) error {
	fmt.Fprintf(i.out, formatDelete(fmt.Sprintf("Identity Pool User %s", opts.UserEmail)))

	return nil
}
