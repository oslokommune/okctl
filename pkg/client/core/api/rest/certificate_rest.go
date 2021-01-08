package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetCertificate matches the REST API route
const TargetCertificate = "certificates/"

type certificateAPI struct {
	client *HTTPClient
}

// nolint: godox
// TODO - Consider making own endpoint to delete certificate in separate
func (a *certificateAPI) DeleteCertificate(opts api.DeleteCertificateOpts) error {
	panic("Not implemented, called directly inside api when running delete identitypool")
}

func (a *certificateAPI) CreateCertificate(opts api.CreateCertificateOpts) (*api.Certificate, error) {
	into := &api.Certificate{}
	return into, a.client.DoPost(TargetCertificate, &opts, into)
}

// NewCertificateAPI returns an initialised REST API client
func NewCertificateAPI(client *HTTPClient) client.CertificateAPI {
	return &certificateAPI{
		client: client,
	}
}
