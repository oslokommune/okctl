// Package acmapi knows how to interact with AWS ACM API
package acmapi

import (
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// ErrNotFound is a not found error
var ErrNotFound = errors.New("not found")

// ACMAPI contains state for interacting with the API
type ACMAPI struct {
	provider v1alpha1.CloudProvider
}

// New returns an initialised ACM API client
func New(provider v1alpha1.CloudProvider) *ACMAPI {
	return &ACMAPI{
		provider: provider,
	}
}

// InUseBy returns the list of ARNs that are currently using
// a certificate
func (a *ACMAPI) InUseBy(certificateARN string) ([]string, error) {
	cert, err := a.provider.ACM().DescribeCertificate(&acm.DescribeCertificateInput{
		CertificateArn: aws.String(certificateARN),
	})
	if err != nil {
		return nil, fmt.Errorf(constant.ErrorDescribingCert, err)
	}

	inUseBy := make([]string, len(cert.Certificate.InUseBy))

	for i, arn := range cert.Certificate.InUseBy {
		inUseBy[i] = *arn
	}

	return inUseBy, nil
}

// CertificateARNForDomain returns the certificate arn for the domain
func (a *ACMAPI) CertificateARNForDomain(domain string) (string, error) {
	var nextToken *string

	for {
		certs, err := a.provider.ACM().ListCertificates(&acm.ListCertificatesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return "", err
		}

		for _, cert := range certs.CertificateSummaryList {
			if *cert.DomainName == domain {
				return *cert.CertificateArn, nil
			}
		}

		if certs.NextToken == nil {
			break
		}

		nextToken = certs.NextToken
	}

	return "", ErrNotFound
}
