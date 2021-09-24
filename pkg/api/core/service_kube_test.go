package core

import (
	"context"
	"testing"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestEnsureTimeoutType(t *testing.T) {
	// Had a case where the error kind got overridden. The service layer is supposed to convert what ever the lower
	// levels return into meaningful errors for okctl
	ks := NewKubeService(mockKubeRun{})

	_, err := ks.CreateExternalDNSKubeDeployment(context.Background(), api.CreateExternalDNSKubeDeploymentOpts{
		ID:           createValidID(),
		HostedZoneID: "someid",
		DomainFilter: "somefilter",
	})

	assert.True(t, errors.IsKind(err, errors.Timeout))
}

func createValidID() api.ID {
	return api.ID{
		Region:       "eu-west-1",
		AWSAccountID: "123456789012",
		ClusterName:  "bugged-test",
	}
}

type mockKubeRun struct{}

func (m mockKubeRun) DisableEarlyDEMUX(_ context.Context, _ api.ID) error {
	panic("implement me")
}

func (m mockKubeRun) CreateExternalDNSKubeDeployment(_ api.CreateExternalDNSKubeDeploymentOpts) (*api.ExternalDNSKube, error) {
	return nil, wait.ErrWaitTimeout
}

func (m mockKubeRun) CreateStorageClass(_ api.CreateStorageClassOpts) (*api.StorageClassKube, error) {
	panic("implement me")
}

func (m mockKubeRun) CreateExternalSecrets(_ api.CreateExternalSecretsOpts) (*api.ExternalSecretsKube, error) {
	panic("implement me")
}

func (m mockKubeRun) DeleteExternalSecrets(_ api.DeleteExternalSecretsOpts) error {
	panic("implement me")
}

func (m mockKubeRun) CreateConfigMap(_ api.CreateConfigMapOpts) (*api.ConfigMap, error) {
	panic("implement me")
}

func (m mockKubeRun) CreateNamespace(_ api.CreateNamespaceOpts) (*api.Namespace, error) {
	panic("implement me")
}

func (m mockKubeRun) DeleteNamespace(_ api.DeleteNamespaceOpts) error { panic("implement me") }
func (m mockKubeRun) DeleteConfigMap(_ api.DeleteConfigMapOpts) error { panic("implement me") }
func (m mockKubeRun) ScaleDeployment(_ api.ScaleDeploymentOpts) error { panic("implement me") }
