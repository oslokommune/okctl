// Package cleaner knows how to clean things up
package cleaner

import (
	"errors"

	"github.com/oslokommune/okctl/pkg/acmapi"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	arnpkg "github.com/oslokommune/okctl/pkg/arn"
	"github.com/oslokommune/okctl/pkg/elbv2api"
)

// Cleaner contains state for cleaning things up
type Cleaner struct {
	elbv2  *elbv2api.ELBv2API
	acmapi *acmapi.ACMAPI
}

// New returns an initialised cleaner
func New(provider v1alpha1.CloudProvider) *Cleaner {
	return &Cleaner{
		elbv2:  elbv2api.New(provider),
		acmapi: acmapi.New(provider),
	}
}

// RemoveThingsUsingCertForDomain removes things using certificate after finding cert
// for domain
func (c *Cleaner) RemoveThingsUsingCertForDomain(domain string) error {
	certificateARN, err := c.acmapi.CertificateARNForDomain(domain)
	if err != nil {
		if errors.Is(err, acmapi.ErrNotFound) {
			return nil
		}

		return err
	}

	return c.RemoveThingsThatAreUsingCertificate(certificateARN)
}

// RemoveThingsThatAreUsingCertificate removes usages of a certificate
func (c *Cleaner) RemoveThingsThatAreUsingCertificate(certificateARN string) error {
	arns, err := c.acmapi.InUseBy(certificateARN)
	if err != nil {
		return err
	}

	for _, arn := range arns {
		isLoadBalancer, err := arnpkg.IsLoadBalancer(arn)
		if err != nil {
			return err
		}

		if isLoadBalancer {
			listeners, err := c.elbv2.GetListenersForLoadBalancer(arn)
			if err != nil {
				return err
			}

			err = c.elbv2.DeleteListenersWithCertificate(certificateARN, listeners)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
