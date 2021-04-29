// Package acmapi knows how to interact with AWS ACM API
package acmapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

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
		return nil, fmt.Errorf("describing certificate: %w", err)
	}

	inUseBy := make([]string, len(cert.Certificate.InUseBy))

	for i, arn := range cert.Certificate.InUseBy {
		inUseBy[i] = *arn
	}

	return inUseBy, nil
}
