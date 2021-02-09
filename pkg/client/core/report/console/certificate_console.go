package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type certificateReport struct {
	console *Console
}

func (r *certificateReport) RemoveCertificate(domain string, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s", aurora.Green(domain))

	return r.console.Report(actions, "certificate", description)
}

func (r *certificateReport) SaveCertificate(certificate *api.Certificate, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s (arn: %s)",
		aurora.Green(certificate.Domain),
		certificate.CertificateARN,
	)

	return r.console.Report(actions, "certificate", description)
}

// NewCertificateReport returns an initialised reporter
func NewCertificateReport(out io.Writer, spinner spinner.Spinner) client.CertificateReport {
	return &certificateReport{
		console: New(out, spinner),
	}
}
