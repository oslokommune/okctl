package core

import (
	"bytes"
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/stretchr/testify/assert"
)

type clusterConfigRetrieverFn func(config *v1alpha5.ClusterConfig)

func createMonitoringService(retriever clusterConfigRetrieverFn) client.MonitoringService {
	var stdout bytes.Buffer

	spin, _ := spinner.New("test", &stdout)

	return NewMonitoringService(
		spin,
		mockAPI{},
		mockStore{},
		mockState{},
		mockReport{},
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
			Environment:  "test",
			Repository:   "something",
			ClusterName:  "something",
		},
		Domain: "test.oslo.systems",
	})
	assert.Nil(t, err)

	assert.NotEmpty(t, clusterConfig.IAM.ServiceAccounts[0].AttachPolicyARNs[0])
}

type mockAPI struct{}

func (m mockAPI) CreateKubePromStack(_ api.CreateKubePrometheusStackOpts) (*api.Helm, error) {
	panic("implement me")
}

func (m mockAPI) DeleteKubePromStack(_ api.DeleteHelmReleaseOpts) error {
	return nil
}

func (m mockAPI) CreateLoki(_ client.CreateLokiOpts) (*api.Helm, error) { panic("implement me") }

func (m mockAPI) DeleteLoki(_ api.DeleteHelmReleaseOpts) error { panic("implement me") }

func (m mockAPI) CreatePromtail(_ client.CreatePromtailOpts) (*api.Helm, error) {
	panic("implement me")
}

func (m mockAPI) DeletePromtail(_ api.DeleteHelmReleaseOpts) error { panic("implement me") }

func (m mockAPI) CreateTempo(_ api.CreateHelmReleaseOpts) (*api.Helm, error) { panic("implement me") }

func (m mockAPI) DeleteTempo(_ api.DeleteHelmReleaseOpts) error { panic("implement me") }

type mockStore struct{}

func (m mockStore) SaveKubePromStack(_ *client.KubePromStack) (*store.Report, error) {
	panic("implement me")
}

func (m mockStore) RemoveKubePromStack(_ api.ID) (*store.Report, error) {
	return nil, nil
}

func (m mockStore) SaveLoki(_ *client.Loki) (*store.Report, error) {
	panic("implement me")
}

func (m mockStore) RemoveLoki(_ api.ID) (*store.Report, error) {
	panic("implement me")
}

func (m mockStore) SavePromtail(_ *client.Promtail) (*store.Report, error) {
	panic("implement me")
}

func (m mockStore) RemovePromtail(_ api.ID) (*store.Report, error) {
	panic("implement me")
}

func (m mockStore) SaveTempo(_ *client.Tempo) (*store.Report, error) {
	panic("implement me")
}

func (m mockStore) RemoveTempo(_ api.ID) (*store.Report, error) {
	panic("implement me")
}

type mockState struct{}

func (m mockState) GetKubePromStack() (*client.KubePromStack, error) {
	return &client.KubePromStack{
		FargateCloudWatchPolicyARN:        "",
		FargateProfilePodExecutionRoleARN: mock.DefaultFargateProfilePodExecutionRoleARN,
	}, nil
}

func (m mockState) SaveKubePromStack(_ *client.KubePromStack) (*store.Report, error) {
	panic("implement me")
}

func (m mockState) RemoveKubePromStack(_ api.ID) (*store.Report, error) {
	return nil, nil
}

type mockReport struct{}

func (m mockReport) ReportSaveKubePromStack(_ *client.KubePromStack, _ []*store.Report) error {
	panic("implement me")
}

func (m mockReport) ReportRemoveKubePromStack(_ []*store.Report) error {
	return nil
}

func (m mockReport) ReportSaveLoki(_ *client.Loki, _ *store.Report) error {
	panic("implement me")
}

func (m mockReport) ReportRemoveLoki(_ *store.Report) error {
	panic("implement me")
}

func (m mockReport) ReportSavePromtail(_ *client.Promtail, _ *store.Report) error {
	panic("implement me")
}

func (m mockReport) ReportRemovePromtail(_ *store.Report) error {
	panic("implement me")
}

func (m mockReport) ReportSaveTempo(_ *client.Tempo, _ *store.Report) error {
	panic("implement me")
}

func (m mockReport) ReportRemoveTempo(_ *store.Report) error {
	panic("implement me")
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

func (m mockIdentityManagerService) CreateIdentityPool(_ context.Context, _ api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	panic("implement me")
}

func (m mockIdentityManagerService) CreateIdentityPoolClient(_ context.Context, _ api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
	panic("implement me")
}

func (m mockIdentityManagerService) CreateIdentityPoolUser(_ context.Context, _ client.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error) {
	panic("implement me")
}

func (m mockIdentityManagerService) DeleteIdentityPool(_ context.Context, _ api.ID) error {
	panic("implement me")
}

func (m mockIdentityManagerService) DeleteIdentityPoolClient(_ context.Context, _ api.DeleteIdentityPoolClientOpts) error {
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

func (m mockParameterService) CreateSecret(_ context.Context, _ api.CreateSecretOpts) (*api.SecretParameter, error) {
	panic("implement me")
}

func (m mockParameterService) DeleteSecret(_ context.Context, _ api.DeleteSecretOpts) error {
	return nil
}

func (m mockParameterService) DeleteAllsecrets(_ context.Context, _ state.Cluster) error {
	panic("implement me")
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
