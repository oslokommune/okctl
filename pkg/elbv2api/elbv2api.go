// Package elbv2api knows how to interact with the ELBv2API
package elbv2api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// ELBv2API contains state for interacting with the API
type ELBv2API struct {
	provider v1alpha1.CloudProvider
}

// New returns an initialised client for interacting with AWS ELBv2 API
func New(provider v1alpha1.CloudProvider) *ELBv2API {
	return &ELBv2API{
		provider: provider,
	}
}

// GetListenersForLoadBalancer returns the listeners associated with a load balancer
func (a *ELBv2API) GetListenersForLoadBalancer(loadbalancerARN string) ([]*elbv2.Listener, error) {
	var marker *string = nil

	var all []*elbv2.Listener

	for {
		listeners, err := a.provider.ELBV2().DescribeListeners(&elbv2.DescribeListenersInput{
			LoadBalancerArn: aws.String(loadbalancerARN),
			Marker:          marker,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case elbv2.ErrCodeLoadBalancerNotFoundException:
					// If the load balancer doesn't exist, not much
					// point in continuing..
					break
				}
			}

			return nil, fmt.Errorf("describing listeners: %w", err)
		}

		all = append(all, listeners.Listeners...)

		if listeners.NextMarker == nil {
			break
		}

		marker = listeners.NextMarker
	}

	return all, nil
}

// DeleteListenersWithCertificate removes the listener if it is using a given certificate
func (a *ELBv2API) DeleteListenersWithCertificate(certificateARN string, listeners []*elbv2.Listener) error {
	for _, listener := range listeners {
		for _, cert := range listener.Certificates {
			if *cert.CertificateArn == certificateARN {
				_, err := a.provider.ELBV2().DeleteListener(&elbv2.DeleteListenerInput{
					ListenerArn: listener.ListenerArn,
				})
				if err != nil {
					return fmt.Errorf("deleting listener: %w", err)
				}
			}
		}
	}

	return nil
}
