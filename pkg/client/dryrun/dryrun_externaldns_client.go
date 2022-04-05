package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type externalDNSService struct {
	out io.Writer
}

func (e externalDNSService) CreateExternalDNS(_ context.Context, _ client.CreateExternalDNSOpts) (*client.ExternalDNS, error) {
	fmt.Fprintf(e.out, formatCreate("ExternalDNS controller"))

	return &client.ExternalDNS{}, nil
}

func (e externalDNSService) DeleteExternalDNS(_ context.Context, _ api.ID) error {
	fmt.Fprintf(e.out, formatDelete("ExternalDNS controller"))

	return nil
}
