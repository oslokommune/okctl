package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	merrors "github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/iamapi"

	"github.com/oslokommune/okctl/pkg/eksapi"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/helm/charts/tempo"

	"github.com/oslokommune/okctl/pkg/datasource"
	"sigs.k8s.io/yaml"

	"github.com/oslokommune/okctl/pkg/helm/charts/promtail"

	"github.com/oslokommune/okctl/pkg/helm/charts/loki"

	"github.com/oslokommune/okctl/pkg/helm/charts/kubepromstack"

	"github.com/miekg/dns"

	"github.com/google/uuid"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type monitoringService struct {
	state                client.MonitoringState
	helm                 client.HelmService
	cert                 client.CertificateService
	ident                client.IdentityManagerService
	param                client.ParameterService
	manifest             client.ManifestService
	service              client.ServiceAccountService
	policy               client.ManagedPolicyService
	provider             v1alpha1.CloudProvider
	objectStorageService api.ObjectStorageService
	keyValueStoreService api.KeyValueStoreService
}

const (
	grafanaSubDomain = "grafana"
	grafanaPurpose   = "grafana"
	clientSecretName = "client-secret"
	secretKeyName    = "secret-key"
	adminUserName    = "admin-user"
	adminPassName    = "admin-pass"
	secretsCfgName   = "grafana-secrets-cm"

	lokiDatasourceConfigMapName        = "loki-datasource"
	tempoDatasourceConfigMapName       = "tempo-datasource"
	cloudwatchDatasourceConfigMapName  = "cloudwatch-datasource"
	notifiersProvisioningConfigMapName = "kube-prometheus-stack-grafana-notifiers"

	kubepromChartTimeout = 15 * time.Minute
)

func grafanaDomain(baseDomain string) string {
	return fmt.Sprintf("%s.%s", grafanaSubDomain, baseDomain)
}

func (s *monitoringService) DeleteTempo(ctx context.Context, id api.ID) error {
	err := s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        id,
		Name:      tempoDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	err = s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          id,
		ReleaseName: tempo.ReleaseName,
		Namespace:   tempo.Namespace,
	})
	if err != nil {
		return err
	}

	return nil
}

// nolint: funlen
func (s *monitoringService) CreateTempo(ctx context.Context, id api.ID) (*client.Helm, error) {
	chart := tempo.New(tempo.NewDefaultValues(), constant.DefaultChartApplyTimeout)

	values, err := chart.ValuesYAML()
	if err != nil {
		return nil, err
	}

	c, err := s.helm.CreateHelmRelease(ctx, client.CreateHelmReleaseOpts{
		ID:             id,
		RepositoryName: chart.RepositoryName,
		RepositoryURL:  chart.RepositoryURL,
		ReleaseName:    chart.ReleaseName,
		Version:        chart.Version,
		Chart:          chart.Chart,
		Namespace:      chart.Namespace,
		Values:         values,
	})
	if err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(datasource.NewTempo())
	if err != nil {
		return nil, err
	}

	_, err = s.manifest.CreateConfigMap(ctx, client.CreateConfigMapOpts{
		ID:        id,
		Name:      tempoDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
		Data: map[string]string{
			"tempo-datasource.yaml": string(data),
		},
		Labels: map[string]string{
			"grafana_datasource": "1",
		},
	})
	if err != nil {
		return nil, err
	}

	// The datasources are only loaded during grafana startup, so we need
	// to cycle grafana to have it pick up the changes
	for _, replicas := range []int32{0, 1} {
		err = s.manifest.ScaleDeployment(ctx, api.ScaleDeploymentOpts{
			ID:        id,
			Name:      constant.DefaultKubePrometheusStackGrafanaName,
			Namespace: constant.DefaultMonitoringNamespace,
			Replicas:  replicas,
		})
		if err != nil {
			if merrors.IsKind(err, merrors.NotExist) {
				break
			}

			return nil, err
		}
	}

	return c, nil
}

func (s *monitoringService) DeletePromtail(ctx context.Context, id api.ID) error {
	err := s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          id,
		ReleaseName: promtail.ReleaseName,
		Namespace:   promtail.Namespace,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *monitoringService) CreatePromtail(ctx context.Context, id api.ID) (*client.Helm, error) {
	chart := promtail.New(promtail.NewDefaultValues(), constant.DefaultChartApplyTimeout)

	values, err := chart.ValuesYAML()
	if err != nil {
		return nil, err
	}

	c, err := s.helm.CreateHelmRelease(ctx, client.CreateHelmReleaseOpts{
		ID:             id,
		RepositoryName: chart.RepositoryName,
		RepositoryURL:  chart.RepositoryURL,
		ReleaseName:    chart.ReleaseName,
		Version:        chart.Version,
		Chart:          chart.Chart,
		Namespace:      chart.Namespace,
		Values:         values,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *monitoringService) DeleteLoki(ctx context.Context, id api.ID) error {
	err := s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        id,
		Name:      lokiDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	err = s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          id,
		ReleaseName: loki.ReleaseName,
		Namespace:   loki.Namespace,
	})
	if err != nil && !merrors.IsKind(err, merrors.NotExist) {
		return err
	}

	config, err := clusterconfig.NewLokiServiceAccount(
		id.ClusterName,
		id.Region,
		constant.DefaultMonitoringNamespace,
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
		[]string{"N/A"},
	)
	if err != nil {
		return fmt.Errorf("generating service account config: %w", err)
	}

	err = s.service.DeleteServiceAccount(ctx, client.DeleteServiceAccountOpts{
		ID:     id,
		Name:   defaultLokiServiceAccountName,
		Config: config,
	})
	if err != nil {
		return fmt.Errorf("deleting service account: %w", err)
	}

	err = deleteDynamoDBIntegration(deleteDynamoDBOpts{
		ctx:                  ctx,
		policyService:        s.policy,
		keyValueStoreService: s.keyValueStoreService,
		provider:             s.provider,
		clusterID:            id,
	})
	if err != nil {
		return fmt.Errorf("deleting DynamoDB integration: %w", err)
	}

	err = deleteS3Integration(deleteS3IntegrationOpts{
		ctx:                  ctx,
		policyService:        s.policy,
		objectStorageService: s.objectStorageService,
		clusterID:            id,
	})
	if err != nil {
		return fmt.Errorf("deleting S3 integration: %w", err)
	}

	return nil
}

// nolint: funlen
func (s *monitoringService) CreateLoki(ctx context.Context, id api.ID) (*client.Helm, error) {
	absoluteBucketName := bucketNameGenerator(id.ClusterName)
	tablePrefix := tablePrefixGenerator(id.ClusterName)

	s3PolicyARN, err := createBucket(createBucketOpts{
		ctx:                  ctx,
		id:                   id,
		objectStorageService: s.objectStorageService,
		policyService:        s.policy,
		bucketName:           absoluteBucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("creating bucket: %w", err)
	}

	dynamoDBPolicy, err := createDynamoDBPolicy(createDynamoDBPolicyOpts{
		ctx:           ctx,
		id:            id,
		policyService: s.policy,
		tablePrefix:   tablePrefix,
	})
	if err != nil {
		return nil, fmt.Errorf("creating DynamoDB policy: %w", err)
	}

	clusterConfig, err := clusterconfig.NewLokiServiceAccount(
		id.ClusterName,
		id.Region,
		constant.DefaultMonitoringNamespace,
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
		[]string{s3PolicyARN, dynamoDBPolicy},
	)
	if err != nil {
		return nil, fmt.Errorf("preparing service account: %w", err)
	}

	_, err = s.service.CreateServiceAccount(ctx, client.CreateServiceAccountOpts{
		ID:        id,
		Name:      defaultLokiServiceAccountName,
		PolicyArn: dynamoDBPolicy,
		Config:    clusterConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("creating service account: %w", err)
	}

	chart := loki.New(loki.NewDefaultValues(absoluteBucketName, tablePrefix), constant.DefaultChartApplyTimeout)

	values, err := chart.ValuesYAML()
	if err != nil {
		return nil, err
	}

	c, err := s.helm.CreateHelmRelease(ctx, client.CreateHelmReleaseOpts{
		ID:             id,
		RepositoryName: chart.RepositoryName,
		RepositoryURL:  chart.RepositoryURL,
		ReleaseName:    chart.ReleaseName,
		Version:        chart.Version,
		Chart:          chart.Chart,
		Namespace:      chart.Namespace,
		Values:         values,
	})
	if err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(datasource.NewLoki())
	if err != nil {
		return nil, err
	}

	_, err = s.manifest.CreateConfigMap(ctx, client.CreateConfigMapOpts{
		ID:        id,
		Name:      lokiDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
		Data: map[string]string{
			"loki-datasource.yaml": string(data),
		},
		Labels: map[string]string{
			"grafana_datasource": "1",
		},
	})
	if err != nil {
		return nil, err
	}

	// The datasources are only loaded during grafana startup, so we need
	// to cycle grafana to have it pick up the changes
	for _, replicas := range []int32{0, 1} {
		err = s.manifest.ScaleDeployment(ctx, api.ScaleDeploymentOpts{
			ID:        id,
			Name:      constant.DefaultKubePrometheusStackGrafanaName,
			Namespace: constant.DefaultMonitoringNamespace,
			Replicas:  replicas,
		})
		if err != nil {
			if merrors.IsKind(err, merrors.NotExist) {
				break
			}

			return nil, err
		}
	}

	return c, nil
}

//nolint:funlen,gocyclo
func (s *monitoringService) DeleteKubePromStack(ctx context.Context, opts client.DeleteKubePromStackOpts) error {
	stack, err := s.state.GetKubePromStack()
	if err != nil {
		return err
	}

	err = iamapi.New(s.provider).DetachRolePolicy(stack.FargateCloudWatchPolicyARN, stack.FargateProfilePodExecutionRoleARN)
	if err != nil {
		return err
	}

	err = s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      cloudwatchDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	err = s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          opts.ID,
		ReleaseName: kubepromstack.ReleaseName,
		Namespace:   kubepromstack.Namespace,
	})
	if err != nil {
		return err
	}

	err = s.manifest.DeleteExternalSecret(ctx, client.DeleteExternalSecretOpts{
		ID:   opts.ID,
		Name: secretsCfgName,
		Secrets: map[string]string{
			secretsCfgName: constant.DefaultMonitoringNamespace,
		},
	})
	if err != nil {
		return err
	}

	err = s.ident.DeleteIdentityPoolClient(ctx, client.DeleteIdentityPoolClientOpts{
		ID:      opts.ID,
		Purpose: grafanaPurpose,
	})
	if err != nil {
		return err
	}

	for _, secretName := range []string{clientSecretName, secretKeyName, adminUserName, adminPassName} {
		err = s.param.DeleteSecret(ctx, client.DeleteSecretOpts{
			ID:   opts.ID,
			Name: secretName,
		})
		if err != nil {
			return fmt.Errorf("deleting parameter secret: %w", err)
		}
	}

	err = s.cert.DeleteCertificate(ctx, client.DeleteCertificateOpts{
		ID:     opts.ID,
		Domain: grafanaDomain(opts.Domain),
	})
	if err != nil {
		return err
	}

	cc, err := clusterconfig.NewCloudwatchDatasourceServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		constant.DefaultMonitoringNamespace,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
		[]string{"N/A"},
	)
	if err != nil {
		return err
	}

	err = s.service.DeleteServiceAccount(ctx, client.DeleteServiceAccountOpts{
		ID:     opts.ID,
		Name:   "cloudwatch-datasource", // Make this configurable
		Config: cc,
	})
	if err != nil {
		return err
	}

	err = s.policy.DeletePolicy(ctx, client.DeletePolicyOpts{
		ID:        opts.ID,
		StackName: cfn.NewStackNamer().CloudwatchDatasource(opts.ID.ClusterName),
	})
	if err != nil {
		return err
	}

	err = s.policy.DeletePolicy(ctx, client.DeletePolicyOpts{
		ID:        opts.ID,
		StackName: cfn.NewStackNamer().FargateCloudwatch(opts.ID.ClusterName),
	})
	if err != nil {
		return err
	}

	err = s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      notifiersProvisioningConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	err = s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      cloudwatchDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	err = s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      "aws-logging",
		Namespace: constant.DefaultFargateObservabilityNamespace,
	})
	if err != nil {
		return err
	}

	err = s.manifest.DeleteNamespace(ctx, api.DeleteNamespaceOpts{
		ID:        opts.ID,
		Namespace: constant.DefaultFargateObservabilityNamespace,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveKubePromStack()
	if err != nil {
		return err
	}

	return nil
}

// nolint: funlen, gocyclo
func (s *monitoringService) CreateKubePromStack(ctx context.Context, opts client.CreateKubePromStackOpts) (*client.KubePromStack, error) {
	cft, err := cfn.New(components.NewCloudwatchDatasourcePolicyComposer(opts.ID.ClusterName)).Build()
	if err != nil {
		return nil, err
	}

	policy, err := s.policy.CreatePolicy(ctx, client.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              cfn.NewStackNamer().CloudwatchDatasource(opts.ID.ClusterName),
		PolicyOutputName:       "CloudwatchDatasourcePolicy", // We need to cleanup the way we name outputs
		CloudFormationTemplate: cft,
	})
	if err != nil {
		return nil, err
	}

	cc, err := clusterconfig.NewCloudwatchDatasourceServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		constant.DefaultMonitoringNamespace,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
		[]string{policy.PolicyARN},
	)
	if err != nil {
		return nil, err
	}

	_, err = s.service.CreateServiceAccount(ctx, client.CreateServiceAccountOpts{
		ID:        opts.ID,
		Name:      "cloudwatch-datasource", // Like, why? We need to make these configurable
		PolicyArn: policy.PolicyARN,
		Config:    cc,
	})
	if err != nil {
		return nil, err
	}

	cert, err := s.cert.CreateCertificate(ctx, client.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         dns.Fqdn(grafanaDomain(opts.Domain)),
		Domain:       grafanaDomain(opts.Domain),
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, err
	}

	poolClient, err := s.ident.CreateIdentityPoolClient(ctx, client.CreateIdentityPoolClientOpts{
		ID:          opts.ID,
		UserPoolID:  opts.UserPoolID,
		Purpose:     grafanaPurpose,
		CallbackURL: fmt.Sprintf("https://%s/login/generic_oauth", cert.Domain),
	})
	if err != nil {
		return nil, err
	}

	clientSecret, err := s.param.CreateSecret(ctx, client.CreateSecretOpts{
		ID:     opts.ID,
		Name:   clientSecretName,
		Secret: poolClient.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	cookieSecret, err := s.param.CreateSecret(ctx, client.CreateSecretOpts{
		ID:     opts.ID,
		Name:   secretKeyName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	adminUser, err := s.param.CreateSecret(ctx, client.CreateSecretOpts{
		ID:     opts.ID,
		Name:   adminUserName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	adminPass, err := s.param.CreateSecret(ctx, client.CreateSecretOpts{
		ID:     opts.ID,
		Name:   adminPassName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	_, err = s.manifest.CreateNamespace(ctx, api.CreateNamespaceOpts{
		ID:        opts.ID,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return nil, err
	}

	manifest, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID:        opts.ID,
		Name:      secretsCfgName,
		Namespace: constant.DefaultMonitoringNamespace,
		Manifest: api.Manifest{
			Name:      secretsCfgName,
			Namespace: constant.DefaultMonitoringNamespace,
			Backend:   api.BackendTypeParameterStore,
			Data: []api.Data{
				{
					Key:  clientSecret.Path,
					Name: clientSecret.Name,
				},
				{
					Key:  cookieSecret.Path,
					Name: cookieSecret.Name,
				},
				{
					Key:  adminUser.Path,
					Name: adminUser.Name,
				},
				{
					Key:  adminPass.Path,
					Name: adminPass.Name,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = s.manifest.CreateConfigMap(ctx, client.CreateConfigMapOpts{
		ID:        opts.ID,
		Name:      notifiersProvisioningConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
		Data: map[string]string{
			"notifiers.yaml": string(""),
		},
		Labels: nil,
	})
	if err != nil {
		return nil, err
	}

	chart := kubepromstack.New(&kubepromstack.Values{ // nolint: gomnd
		GrafanaServiceAccountName:          constant.DefaultGrafanaCloudWatchDatasourceName,
		GrafanaCertificateARN:              cert.ARN,
		GrafanaHostname:                    cert.Domain,
		AuthHostname:                       opts.AuthDomain,
		ClientID:                           poolClient.ClientID,
		SecretsConfigName:                  secretsCfgName,
		SecretsGrafanaCookieSecretKey:      cookieSecret.Name,
		SecretsGrafanaOauthClientSecretKey: clientSecret.Name,
		SecretsGrafanaAdminUserKey:         adminUser.Name,
		SecretsGrafanaAdminPassKey:         adminPass.Name,
	}, kubepromChartTimeout)

	values, err := chart.ValuesYAML()
	if err != nil {
		return nil, err
	}

	c, err := s.helm.CreateHelmRelease(ctx, client.CreateHelmReleaseOpts{
		ID:             opts.ID,
		RepositoryName: chart.RepositoryName,
		RepositoryURL:  chart.RepositoryURL,
		ReleaseName:    chart.ReleaseName,
		Version:        chart.Version,
		Chart:          chart.Chart,
		Namespace:      chart.Namespace,
		Values:         values,
	})
	if err != nil {
		return nil, err
	}

	stack := &client.KubePromStack{
		ID:                     opts.ID,
		CertificateARN:         cert.ARN,
		Hostname:               cert.Domain,
		AuthHostname:           opts.AuthDomain,
		ClientID:               poolClient.ClientID,
		SecretsConfigName:      secretsCfgName,
		SecretsCookieSecretKey: cookieSecret.Name,
		SecretsClientSecretKey: clientSecret.Name,
		SecretsAdminUserKey:    adminUser.Name,
		SecretsAdminPassKey:    adminPass.Name,
		Chart:                  c,
		Certificate:            cert,
		IdentityPoolClient:     poolClient,
		ExternalSecret:         manifest,
	}

	data, err := yaml.Marshal(datasource.NewCloudWatch(opts.ID.Region))
	if err != nil {
		return nil, err
	}

	// should move this into the default datasources instead
	// probably
	_, err = s.manifest.CreateConfigMap(ctx, client.CreateConfigMapOpts{
		ID:        opts.ID,
		Name:      cloudwatchDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
		Data: map[string]string{
			"cloudwatch-datasource.yaml": string(data),
		},
		Labels: map[string]string{
			"grafana_datasource": "1",
		},
	})
	if err != nil {
		return nil, err
	}

	// The datasources are only loaded during grafana startup, so we need
	// to cycle grafana to have it pick up the changes
	for _, replicas := range []int32{0, 1} {
		err = s.manifest.ScaleDeployment(ctx, api.ScaleDeploymentOpts{
			ID:        opts.ID,
			Name:      constant.DefaultKubePrometheusStackGrafanaName,
			Namespace: constant.DefaultMonitoringNamespace,
			Replicas:  replicas,
		})
		if err != nil {
			return nil, err
		}
	}

	_, err = s.manifest.CreateNamespace(ctx, api.CreateNamespaceOpts{
		ID:        opts.ID,
		Namespace: constant.DefaultFargateObservabilityNamespace,
		Labels: map[string]string{
			"aws-observability": "enabled",
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = s.manifest.CreateConfigMap(ctx, client.CreateConfigMapOpts{
		ID:        opts.ID,
		Name:      "aws-logging",
		Namespace: constant.DefaultFargateObservabilityNamespace,
		Data: map[string]string{
			"output.conf": fmt.Sprintf(`
[OUTPUT]
    Name cloudwatch_logs
    Match *
    region %s
    log_group_name okctl-fluent-cloudwatch
    log_stream_prefix from-fluent-bit
    auto_create_group true
`, opts.ID.Region),
		},
		Labels: nil,
	})
	if err != nil {
		return nil, err
	}

	fcp, err := cfn.New(components.NewFargateCloudwatchPolicyComposer(opts.ID.ClusterName)).Build()
	if err != nil {
		return nil, err
	}

	fargatePolicy, err := s.policy.CreatePolicy(ctx, client.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              cfn.NewStackNamer().FargateCloudwatch(opts.ID.ClusterName),
		PolicyOutputName:       "FargateCloudwatchPolicy", // We need to cleanup the way we name outputs
		CloudFormationTemplate: fcp,
	})
	if err != nil {
		return nil, err
	}

	roleARN, err := eksapi.New(opts.ID.ClusterName, s.provider).
		FargateProfilePodExecutionRoleARN("fp-default") // another string that needs to be set explicitly
	if err != nil {
		return nil, err
	}

	err = iamapi.New(s.provider).AttachRolePolicy(fargatePolicy.PolicyARN, roleARN)
	if err != nil {
		return nil, err
	}

	stack.FargateProfilePodExecutionRoleARN = roleARN
	stack.FargateCloudWatchPolicyARN = fargatePolicy.PolicyARN

	err = s.state.SaveKubePromStack(stack)
	if err != nil {
		return nil, err
	}

	return stack, nil
}

// NewMonitoringServiceOpts contains necessary data for handling monitoring
type NewMonitoringServiceOpts struct {
	State                 client.MonitoringState
	Helm                  client.HelmService
	CertificateService    client.CertificateService
	IdentityService       client.IdentityManagerService
	ManifestService       client.ManifestService
	ParameterService      client.ParameterService
	ServiceAccountService client.ServiceAccountService
	PolicyService         client.ManagedPolicyService
	ObjectStorageService  api.ObjectStorageService
	KeyValueStoreService  api.KeyValueStoreService
	Provider              v1alpha1.CloudProvider
}

// NewMonitoringService returns an initialised service
func NewMonitoringService(opts NewMonitoringServiceOpts) client.MonitoringService {
	return &monitoringService{
		state:                opts.State,
		helm:                 opts.Helm,
		cert:                 opts.CertificateService,
		ident:                opts.IdentityService,
		param:                opts.ParameterService,
		manifest:             opts.ManifestService,
		service:              opts.ServiceAccountService,
		policy:               opts.PolicyService,
		objectStorageService: opts.ObjectStorageService,
		keyValueStoreService: opts.KeyValueStoreService,
		provider:             opts.Provider,
	}
}

const (
	defaultLokiServiceAccountName   = "loki"
	defaultLokiDynamoDBPolicyDomain = "loki"
	defaultLokiBucketName           = "loki"
)

type createBucketOpts struct {
	ctx                  context.Context
	id                   api.ID
	objectStorageService api.ObjectStorageService
	policyService        client.ManagedPolicyService
	bucketName           string
}

func createBucket(opts createBucketOpts) (string, error) {
	bucketARN, err := opts.objectStorageService.CreateBucket(api.CreateBucketOpts{
		ClusterID:  opts.id,
		BucketName: opts.bucketName,
		Private:    true,
		Encrypted:  true,
	})
	if err != nil {
		return "", fmt.Errorf("creating bucket: %w", err)
	}

	s3PolicyCFNStack, err := cfn.New(components.NewLokiS3PolicyComposer(opts.id.ClusterName, bucketARN)).Build()
	if err != nil {
		return "", fmt.Errorf("building S3 policy template: %w", err)
	}

	s3Policy, err := opts.policyService.CreatePolicy(opts.ctx, client.CreatePolicyOpts{
		ID:                     opts.id,
		StackName:              cfn.NewStackNamer().LokiS3Policy(opts.id.ClusterName, defaultLokiBucketName),
		PolicyOutputName:       "LokiS3ServiceAccountPolicy",
		CloudFormationTemplate: s3PolicyCFNStack,
	})
	if err != nil {
		return "", fmt.Errorf("creating S3 policy: %w", err)
	}

	return s3Policy.PolicyARN, nil
}

type createDynamoDBPolicyOpts struct {
	ctx           context.Context
	id            api.ID
	policyService client.ManagedPolicyService
	tablePrefix   string
}

func createDynamoDBPolicy(opts createDynamoDBPolicyOpts) (string, error) {
	// Should cover all tables starting with a certain prefix
	tablePrefix := opts.tablePrefix + "*"

	dynamoDBCFNStack, err := cfn.New(components.NewLokiDynamoDBPolicyComposer(opts.id, tablePrefix)).Build()
	if err != nil {
		return "", fmt.Errorf("building DynamoDB policy template: %w", err)
	}

	dynamoDBPolicy, err := opts.policyService.CreatePolicy(opts.ctx, client.CreatePolicyOpts{
		ID:                     opts.id,
		StackName:              cfn.NewStackNamer().LokiDynamoDBPolicy(opts.id.ClusterName, defaultLokiDynamoDBPolicyDomain),
		PolicyOutputName:       "LokiDynamoDBServiceAccountPolicy",
		CloudFormationTemplate: dynamoDBCFNStack,
	})
	if err != nil {
		return "", fmt.Errorf("creating DynamoDB policy: %w", err)
	}

	return dynamoDBPolicy.PolicyARN, nil
}

type deleteDynamoDBOpts struct {
	ctx                  context.Context
	policyService        client.ManagedPolicyService
	keyValueStoreService api.KeyValueStoreService
	provider             v1alpha1.CloudProvider
	clusterID            api.ID
}

func deleteDynamoDBIntegration(opts deleteDynamoDBOpts) error {
	err := opts.policyService.DeletePolicy(opts.ctx, client.DeletePolicyOpts{
		ID:        opts.clusterID,
		StackName: cfn.NewStackNamer().LokiDynamoDBPolicy(opts.clusterID.ClusterName, defaultLokiDynamoDBPolicyDomain),
	})
	if err != nil {
		return fmt.Errorf("deleting DynamoDB policy: %w", err)
	}

	tables, err := opts.keyValueStoreService.ListStores()
	if err != nil {
		return fmt.Errorf("listing tables: %w", err)
	}

	deletionQueue := make([]string, 0)
	tablePrefix := tablePrefixGenerator(opts.clusterID.ClusterName)

	for _, table := range tables {
		if strings.HasPrefix(table, tablePrefix) {
			deletionQueue = append(deletionQueue, table)
		}
	}

	dynamoDBAPI := opts.provider.DynamoDB()

	for _, table := range deletionQueue {
		_, err = dynamoDBAPI.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
		if err != nil {
			return fmt.Errorf("deleting table: %w", err)
		}
	}

	return nil
}

type deleteS3IntegrationOpts struct {
	ctx                  context.Context
	policyService        client.ManagedPolicyService
	objectStorageService api.ObjectStorageService
	clusterID            api.ID
}

func deleteS3Integration(opts deleteS3IntegrationOpts) error {
	absoluteBucketName := bucketNameGenerator(opts.clusterID.ClusterName)

	err := opts.policyService.DeletePolicy(opts.ctx, client.DeletePolicyOpts{
		ID:        opts.clusterID,
		StackName: cfn.NewStackNamer().LokiS3Policy(opts.clusterID.ClusterName, defaultLokiBucketName),
	})
	if err != nil {
		return fmt.Errorf("deleting S3 policy: %w", err)
	}

	err = opts.objectStorageService.EmptyBucket(api.EmptyBucketOpts{
		BucketName: absoluteBucketName,
	})
	if err != nil && !errors.Is(err, api.ErrObjectStorageBucketNotExist) {
		return fmt.Errorf("emptying bucket: %w", err)
	}

	err = opts.objectStorageService.DeleteBucket(api.DeleteBucketOpts{
		ClusterID:  opts.clusterID,
		BucketName: absoluteBucketName,
	})
	if err != nil && !errors.Is(err, api.ErrObjectStorageBucketNotExist) {
		return fmt.Errorf("deleting bucket: %w", err)
	}

	return nil
}

func tablePrefixGenerator(clusterName string) string {
	return fmt.Sprintf("okctl-%s-loki-index_", clusterName)
}

func bucketNameGenerator(clusterName string) string {
	return fmt.Sprintf("okctl-%s-loki", clusterName)
}
