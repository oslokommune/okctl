package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type autoscalerService struct {
	out io.Writer
}

func (a autoscalerService) CreateAutoscaler(_ context.Context, _ client.CreateAutoscalerOpts) (*client.Autoscaler, error) {
	fmt.Fprintf(a.out, formatCreate("Autoscaler"))

	return &client.Autoscaler{
		Policy:         nil,
		ServiceAccount: nil,
		Chart:          nil,
	}, nil
}

func (a autoscalerService) DeleteAutoscaler(_ context.Context, _ api.ID) error {
	fmt.Fprintf(a.out, formatDelete("Autoscaler"))

	return nil
}
