package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type monitoringService struct {
	out io.Writer
}

func (m monitoringService) CreateKubePromStack(_ context.Context, _ client.CreateKubePromStackOpts) (*client.KubePromStack, error) {
	fmt.Fprintf(m.out, formatCreate("Prometheus & Grafana"))

	return &client.KubePromStack{}, nil
}

func (m monitoringService) DeleteKubePromStack(_ context.Context, _ client.DeleteKubePromStackOpts) error {
	fmt.Fprintf(m.out, formatDelete("Prometheus & Grafana"))

	return nil
}

func (m monitoringService) CreateLoki(_ context.Context, _ api.ID) (*client.Helm, error) {
	fmt.Fprintf(m.out, formatCreate("Loki"))

	return &client.Helm{}, nil
}

func (m monitoringService) DeleteLoki(_ context.Context, _ api.ID) error {
	fmt.Fprintf(m.out, formatDelete("Loki"))

	return nil
}

func (m monitoringService) CreatePromtail(_ context.Context, _ api.ID) (*client.Helm, error) {
	fmt.Fprintf(m.out, formatCreate("Promtail"))

	return &client.Helm{}, nil
}

func (m monitoringService) DeletePromtail(_ context.Context, _ api.ID) error {
	fmt.Fprintf(m.out, formatDelete("Promtail"))

	return nil
}

func (m monitoringService) CreateTempo(_ context.Context, _ api.ID) (*client.Helm, error) {
	fmt.Fprintf(m.out, formatCreate("Tempo"))

	return &client.Helm{}, nil
}

func (m monitoringService) DeleteTempo(_ context.Context, _ api.ID) error {
	fmt.Fprintf(m.out, formatDelete("Tempo"))

	return nil
}
