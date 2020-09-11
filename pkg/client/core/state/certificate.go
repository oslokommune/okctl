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

func (s *certificateState) GetCertificate(domain string) *api.Certificate {
	cert := s.state.GetCertificate(domain)

	return &api.Certificate{
		Domain:         cert.Domain,
		CertificateARN: cert.ARN,
	}
}

func (s *certificateState) SaveCertificate(c *api.Certificate) (*store.Report, error) {
	report, err := s.state.SaveCertificate(&state.Certificate{
		Domain: c.Domain,
		ARN:    c.CertificateARN,
	})
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
