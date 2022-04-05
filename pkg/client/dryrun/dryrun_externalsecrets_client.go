package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type externalSecretsService struct {
	out io.Writer
}

func (e externalSecretsService) CreateExternalSecrets(_ context.Context, _ client.CreateExternalSecretsOpts) (*client.ExternalSecrets, error) {
	fmt.Fprintf(e.out, formatCreate("ExternalSecrets controller"))

	return &client.ExternalSecrets{}, nil
}

func (e externalSecretsService) DeleteExternalSecrets(_ context.Context, _ api.ID) error {
	fmt.Fprintf(e.out, formatDelete("ExternalSecrets controller"))

	return nil
}
