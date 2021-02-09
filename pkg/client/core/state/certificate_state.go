package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type certificateState struct {
	state state.Certificater
}

func (s *certificateState) RemoveCertificate(domain string) (*store.Report, error) {
	return s.state.DeleteCertificate(domain)
}

func (s *certificateState) GetCertificate(domain string) state.Certificate {
	return s.state.GetCertificate(domain)
}

func (s *certificateState) SaveCertificate(c *api.Certificate) (*store.Report, error) {
	cert := s.state.GetCertificate(c.Domain)

	cert.Domain = c.Domain
	cert.ARN = c.CertificateARN

	report, err := s.state.SaveCertificate(cert)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "Certificate",
			Path: fmt.Sprintf("domain=%s, clusterName=%s", c.Domain, c.ID.ClusterName),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

// NewCertificateState returns an initialised state handler
func NewCertificateState(state state.Certificater) client.CertificateState {
	return &certificateState{
		state: state,
	}
}
