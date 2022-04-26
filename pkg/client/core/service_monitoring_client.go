package core

import (
	"context"
	"fmt"
	"time"

	"github.com/mishudark/errors"

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
	state    client.MonitoringState
	helm     client.HelmService
	cert     client.CertificateService
	ident    client.IdentityManagerService
	param    client.ParameterService
	manifest client.ManifestService
	service  client.ServiceAccountService
	policy   client.ManagedPolicyService
	provider v1alpha1.CloudProvider
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
			if errors.IsKind(err, errors.NotExist) {
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
	if err != nil {
		return err
	}

	return nil
}

// nolint: funlen
func (s *monitoringService) CreateLoki(ctx context.Context, id api.ID) (*client.Helm, error) {
	chart := loki.New(loki.NewDefaultValues(), constant.DefaultChartApplyTimeout)

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
			if errors.IsKind(err, errors.NotExist) {
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

// NewMonitoringService returns an initialised service
func NewMonitoringService(
	state client.MonitoringState,
	helm client.HelmService,
	cert client.CertificateService,
	ident client.IdentityManagerService,
	manifest client.ManifestService,
	param client.ParameterService,
	service client.ServiceAccountService,
	policy client.ManagedPolicyService,
	provider v1alpha1.CloudProvider,
) client.MonitoringService {
	return &monitoringService{
		state:    state,
		helm:     helm,
		cert:     cert,
		ident:    ident,
		param:    param,
		manifest: manifest,
		service:  service,
		policy:   policy,
		provider: provider,
	}
}
