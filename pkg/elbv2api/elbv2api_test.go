package elbv2api_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/oslokommune/okctl/pkg/elbv2api"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

func TestELBv2API_GetListenersForLoadBalancer(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		arn       string
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			arn:      mock.DefaultLoadBalancerARN,
			expect: []*elbv2.Listener{
				{
					Certificates: []*elbv2.Certificate{
						{
							CertificateArn: aws.String(mock.DefaultCertificateARN),
						},
					},
					ListenerArn: aws.String(mock.DefaultListenerARN),
				},
			},
			expectErr: false,
		},
		{
			name: "Should ignore lb not found",
			provider: mock.NewGoodCloudProvider().
				DescribeListenersResponse(nil, awserr.New(elbv2.ErrCodeLoadBalancerNotFoundException, "", errors.New("something"))),
			arn:       mock.DefaultLoadBalancerARN,
			expect:    []*elbv2.Listener(nil),
			expectErr: false,
		},
		{
			name:      "Should fail",
			provider:  mock.NewBadCloudProvider(),
			arn:       mock.DefaultLoadBalancerARN,
			expect:    "describing listeners: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := elbv2api.New(tc.provider).GetListenersForLoadBalancer(tc.arn)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}

func TestELBv2API_DeleteListenersWithCertificate(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		arn       string
		listeners []*elbv2.Listener
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			arn:      mock.DefaultCertificateARN,
			listeners: []*elbv2.Listener{
				{
					Certificates: []*elbv2.Certificate{
						{
							CertificateArn: aws.String(mock.DefaultCertificateARN),
						},
					},
					ListenerArn: aws.String(mock.DefaultListenerARN),
				},
			},
		},
		{
			name:     "Should fail",
			provider: mock.NewBadCloudProvider(),
			arn:      mock.DefaultCertificateARN,
			listeners: []*elbv2.Listener{
				{
					Certificates: []*elbv2.Certificate{
						{
							CertificateArn: aws.String(mock.DefaultCertificateARN),
						},
					},
					ListenerArn: aws.String(mock.DefaultListenerARN),
				},
			},
			expect:    "deleting listener: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := elbv2api.New(tc.provider).DeleteListenersWithCertificate(tc.arn, tc.listeners)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
