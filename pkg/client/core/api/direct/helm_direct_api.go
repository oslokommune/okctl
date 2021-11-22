package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type helmDirectClient struct {
	service api.HelmService
}

func (d *helmDirectClient) GetHelmRelease(opts api.GetHelmReleaseOpts) (*api.Helm, error) {
	return d.service.GetHelmRelease(context.Background(), opts)
}

func (d *helmDirectClient) CreateHelmRelease(opts api.CreateHelmReleaseOpts) (*api.Helm, error) {
	return d.service.CreateHelmRelease(context.Background(), opts)
}

func (d *helmDirectClient) DeleteHelmRelease(opts api.DeleteHelmReleaseOpts) error {
	return d.service.DeleteHelmRelease(context.Background(), opts)
}

func NewHelmAPI(service api.HelmService) client.HelmAPI {
	return &helmDirectClient{
		service:service,
	}
}