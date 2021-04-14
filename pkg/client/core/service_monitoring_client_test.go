package core

import (
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/stretchr/testify/assert"
)

type clusterConfigRetrieverFn func(config *v1alpha5.ClusterConfig)

func createMonitoringService(retriever clusterConfigRetrieverFn) client.MonitoringService {
	return NewMonitoringService(
		mockState{},
		mockHelm{},
		mockCertService{},
		mockIdentityManagerService{},
		mockManifestService{},
		mockParameterService{},
		mockServiceAccountService{retrieverFn: retriever},
		mockManagedPolicyService{},
		mock.NewGoodCloudProvider(),
	)
}

func TestKM196(t *testing.T) {
	// Policy ARN was an invalid value and it wasnt caught until runtime
	var clusterConfig *v1alpha5.ClusterConfig

	monitoringService := createMonitoringService(func(config *v1alpha5.ClusterConfig) {
		clusterConfig = config
	})

	err := monitoringService.DeleteKubePromStack(context.Background(), client.DeleteKubePromStackOpts{
		ID: api.ID{
			Region:       "eu-west-1",
			AWSAccountID: "123456789012",
			ClusterName:  "something",
		},
		Domain: "test.oslo.systems",
	})
	assert.Nil(t, err)

	assert.NotEmpty(t, clusterConfig.IAM.ServiceAccounts[0].AttachPolicyARNs[0])
}

type mockHelm struct {
}

func (m mockHelm) CreateHelmRelease(context.Context, client.CreateHelmReleaseOpts) (*client.Helm, error) {
	return nil, nil
}

func (m mockHelm) DeleteHelmRelease(context.Context, client.DeleteHelmReleaseOpts) error {
	return nil
}

type mockState struct{}

func (m mockState) SaveKubePromStack(_ *client.KubePromStack) error {
	return nil
}

func (m mockState) RemoveKubePromStack() error {
	return nil
}

func (m mockState) GetKubePromStack() (*client.KubePromStack, error) {
	return &client.KubePromStack{
		FargateCloudWatchPolicyARN:        "",
		FargateProfilePodExecutionRoleARN: mock.DefaultFargateProfilePodExecutionRoleARN,
	}, nil
}

type mockCertService struct{}

func (m mockCertService) CreateCertificate(_ context.Context, _ client.CreateCertificateOpts) (*client.Certificate, error) {
	panic("implement me")
}

func (m mockCertService) DeleteCertificate(_ context.Context, _ client.DeleteCertificateOpts) error {
	return nil
}

func (m mockCertService) DeleteCognitoCertificate(_ context.Context, _ client.DeleteCognitoCertificateOpts) error {
	panic("implement me")
}

type mockIdentityManagerService struct{}

func (m mockIdentityManagerService) CreateIdentityPool(_ context.Context, _ client.CreateIdentityPoolOpts) (*client.IdentityPool, error) {
	panic("implement me")
}

func (m mockIdentityManagerService) CreateIdentityPoolClient(_ context.Context, _ client.CreateIdentityPoolClientOpts) (*client.IdentityPoolClient, error) {
	panic("implement me")
}

func (m mockIdentityManagerService) CreateIdentityPoolUser(_ context.Context, _ client.CreateIdentityPoolUserOpts) (*client.IdentityPoolUser, error) {
	panic("implement me")
}

func (m mockIdentityManagerService) DeleteIdentityPool(_ context.Context, _ api.ID) error {
	panic("implement me")
}

func (m mockIdentityManagerService) DeleteIdentityPoolClient(_ context.Context, _ client.DeleteIdentityPoolClientOpts) error {
	return nil
}

type mockManifestService struct{}

func (m mockManifestService) DeleteNamespace(_ context.Context, _ api.DeleteNamespaceOpts) error {
	return nil
}

func (m mockManifestService) CreateStorageClass(_ context.Context, _ api.CreateStorageClassOpts) (*client.KubernetesManifest, error) {
	panic("implement me")
}

func (m mockManifestService) CreateExternalSecret(_ context.Context, _ client.CreateExternalSecretOpts) (*client.KubernetesManifest, error) {
	panic("implement me")
}

func (m mockManifestService) DeleteExternalSecret(_ context.Context, _ client.DeleteExternalSecretOpts) error {
	return nil
}

func (m mockManifestService) CreateConfigMap(_ context.Context, _ client.CreateConfigMapOpts) (*client.KubernetesManifest, error) {
	panic("implement me")
}

func (m mockManifestService) DeleteConfigMap(_ context.Context, _ client.DeleteConfigMapOpts) error {
	return nil
}

func (m mockManifestService) ScaleDeployment(_ context.Context, _ api.ScaleDeploymentOpts) error {
	panic("implement me")
}

func (m mockManifestService) CreateNamespace(_ context.Context, _ api.CreateNamespaceOpts) (*client.KubernetesManifest, error) {
	panic("implement me")
}

type mockParameterService struct{}

func (m mockParameterService) CreateSecret(_ context.Context, _ client.CreateSecretOpts) (*client.SecretParameter, error) {
	panic("implement me")
}

func (m mockParameterService) DeleteSecret(_ context.Context, _ client.DeleteSecretOpts) error {
	return nil
}

type mockServiceAccountService struct {
	retrieverFn clusterConfigRetrieverFn
}

func (m mockServiceAccountService) CreateServiceAccount(_ context.Context, _ client.CreateServiceAccountOpts) (*client.ServiceAccount, error) {
	panic("implement me")
}

func (m mockServiceAccountService) DeleteServiceAccount(_ context.Context, opts client.DeleteServiceAccountOpts) error {
	m.retrieverFn(opts.Config)

	return nil
}

type mockManagedPolicyService struct{}

func (m mockManagedPolicyService) CreatePolicy(_ context.Context, _ client.CreatePolicyOpts) (*client.ManagedPolicy, error) {
	panic("implement me")
}

func (m mockManagedPolicyService) DeletePolicy(_ context.Context, _ client.DeletePolicyOpts) error {
	return nil
}
