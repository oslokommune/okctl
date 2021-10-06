package core

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestKM203(t *testing.T) {
	// Ensures the identity manager service returns an error of kind Timeout when create certificate times out
	service := NewIdentityManagerService(mockIdentityManagerCloudProvider{}, mockCertificateCloudProvider{})

	_, err := service.CreateIdentityPool(context.Background(), api.CreateIdentityPoolOpts{})

	assert.True(t, errors.IsKind(err, errors.Timeout))
}

type mockIdentityManagerCloudProvider struct{}

func (m mockIdentityManagerCloudProvider) DeleteIdentityPoolUser(_ api.DeleteIdentityPoolUserOpts) error {
	panic("implement me")
}

func (m mockIdentityManagerCloudProvider) CreateIdentityPool(_ context.Context, _ string, _ api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	panic("implement me")
}

func (m mockIdentityManagerCloudProvider) CreateIdentityPoolClient(_ context.Context, _ api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
	panic("implement me")
}

func (m mockIdentityManagerCloudProvider) CreateIdentityPoolUser(_ context.Context, _ api.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error) {
	panic("implement me")
}

func (m mockIdentityManagerCloudProvider) DeleteIdentityPool(_ api.DeleteIdentityPoolOpts) error {
	panic("implement me")
}

func (m mockIdentityManagerCloudProvider) DeleteIdentityPoolClient(_ api.DeleteIdentityPoolClientOpts) error {
	panic("implement me")
}

type mockCertificateCloudProvider struct{}

func (m mockCertificateCloudProvider) CreateCertificate(_ context.Context, _ api.CreateCertificateOpts) (*api.Certificate, error) {
	return nil, stderrors.New("stack creation time exceeded the specified timeout")
}

func (m mockCertificateCloudProvider) DeleteCertificate(_ api.DeleteCertificateOpts) error {
	panic("implement me")
}

func (m mockCertificateCloudProvider) DeleteCognitoCertificate(_ api.DeleteCognitoCertificateOpts) error {
	panic("implement me")
}
