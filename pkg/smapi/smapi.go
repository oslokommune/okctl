// Package smapi provides some convenience functions for interacting
// with the AWS SecretsManager API
package smapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// SMAPI contains the required state for the secrets manager client
type SMAPI struct {
	provider v1alpha1.CloudProvider
}

const (
	rotateAfterDays int64 = 30
)

// New returns an initialised SecretsManager API Client
func New(provider v1alpha1.CloudProvider) *SMAPI {
	return &SMAPI{
		provider: provider,
	}
}

// RotateSecret enables a secret rotation
func (a *SMAPI) RotateSecret(lambdaARN, secretID string) error {
	_, err := a.provider.SecretsManager().RotateSecret(&secretsmanager.RotateSecretInput{
		RotationLambdaARN: aws.String(lambdaARN),
		RotationRules: &secretsmanager.RotationRulesType{
			AutomaticallyAfterDays: aws.Int64(rotateAfterDays),
		},
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return fmt.Errorf("creating secret rotation: %w", err)
	}

	return nil
}

// CancelRotateSecret removes a secret rotation
func (a *SMAPI) CancelRotateSecret(secretID string) error {
	_, err := a.provider.SecretsManager().CancelRotateSecret(&secretsmanager.CancelRotateSecretInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return fmt.Errorf("canceling secret rotation: %w", err)
	}

	return nil
}
