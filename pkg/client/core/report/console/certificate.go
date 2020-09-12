package console

import (
	"fmt"
	"io"

	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/theckman/yacspin"
)

type certificateReport struct {
	console *Console
}

func (r *certificateReport) SaveCertificate(certificate *api.Certificate, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s (arn: %s)",
		aurora.Blue(certificate.Domain),
		certificate.CertificateARN,
	)

	return r.console.Report(actions, "certificate", description)
}

// NewCertificateReport returns an initialised reporter
func NewCertificateReport(out io.Writer, exit chan struct{}, spinner *yacspin.Spinner) client.CertificateReport {
	return &certificateReport{
		console: New(out, exit, spinner),
	}
}
