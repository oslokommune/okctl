package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type certificateDirectClient struct {
	service api.CertificateService
}

func (c *certificateDirectClient) DeleteCognitoCertificate(opts api.DeleteCognitoCertificateOpts) error {
	return c.service.DeleteCognitoCertificate(context.Background(), opts)
}

func (c *certificateDirectClient) DeleteCertificate(opts api.DeleteCertificateOpts) error {
	return c.service.DeleteCertificate(context.Background(), opts)
}

func (c *certificateDirectClient) CreateCertificate(opts api.CreateCertificateOpts) (*api.Certificate, error) {
	return c.service.CreateCertificate(context.Background(), opts)
}

// NewCertificateAPI returns an initialised API with core service to call directly
func NewCertificateAPI(service api.CertificateService) client.CertificateAPI {
	return &certificateDirectClient{
		service: service,
	}
}
