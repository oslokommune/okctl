package core

import (
	"context"
	"fmt"
	"time"

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

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/google/uuid"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type monitoringService struct {
	spinner spinner.Spinner
	api     client.MonitoringAPI
	store   client.MonitoringStore
	state   client.MonitoringState
	report  client.MonitoringReport

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

	lokiDatasourceConfigMapName       = "loki-datasource"
	tempoDatasourceConfigMapName      = "tempo-datasource"
	cloudwatchDatasourceConfigMapName = "cloudwatch-datasource"
)

func grafanaDomain(baseDomain string) string {
	return fmt.Sprintf("%s.%s", grafanaSubDomain, baseDomain)
}

func (s *monitoringService) DeleteTempo(ctx context.Context, opts client.DeleteTempoOpts) error {
	err := s.spinner.Start("Tempo")
	if err != nil {
		return err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	err = s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      tempoDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	chart := tempo.New(nil)

	err = s.api.DeleteTempo(api.DeleteHelmReleaseOpts{
		ID:          opts.ID,
		ReleaseName: chart.ReleaseName,
		Namespace:   chart.Namespace,
	})
	if err != nil {
		return err
	}

	report, err := s.store.RemoveTempo(opts.ID)
	if err != nil {
		return err
	}

	return s.report.ReportRemoveTempo(report)
}

// nolint: funlen
func (s *monitoringService) CreateTempo(ctx context.Context, opts client.CreateTempoOpts) (*client.Tempo, error) {
	err := s.spinner.Start("Tempo")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	chart := tempo.New(tempo.NewDefaultValues())

	values, err := chart.ValuesYAML()
	if err != nil {
		return nil, err
	}

	c, err := s.api.CreateTempo(api.CreateHelmReleaseOpts{
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

	data, err := yaml.Marshal(datasource.NewTempo())
	if err != nil {
		return nil, err
	}

	_, err = s.manifest.CreateConfigMap(ctx, client.CreateConfigMapOpts{
		ID:        opts.ID,
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
			ID:        opts.ID,
			Name:      constant.DefaultKubePrometheusStackGrafanaName,
			Namespace: constant.DefaultMonitoringNamespace,
			Replicas:  replicas,
		})
		if err != nil {
			return nil, err
		}
	}

	l := &client.Tempo{
		ID:    opts.ID,
		Chart: c,
	}

	report, err := s.store.SaveTempo(l)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportSaveTempo(l, report)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *monitoringService) DeletePromtail(_ context.Context, opts client.DeletePromtailOpts) error {
	err := s.spinner.Start("promtail")
	if err != nil {
		return err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	chart := promtail.New(nil)

	err = s.api.DeletePromtail(api.DeleteHelmReleaseOpts{
		ID:          opts.ID,
		ReleaseName: chart.ReleaseName,
		Namespace:   chart.Namespace,
	})
	if err != nil {
		return err
	}

	report, err := s.store.RemovePromtail(opts.ID)
	if err != nil {
		return err
	}

	return s.report.ReportRemovePromtail(report)
}

func (s *monitoringService) CreatePromtail(_ context.Context, opts client.CreatePromtailOpts) (*client.Promtail, error) {
	err := s.spinner.Start("promtail")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	chart, err := s.api.CreatePromtail(opts)
	if err != nil {
		return nil, err
	}

	l := &client.Promtail{
		ID:    opts.ID,
		Chart: chart,
	}

	report, err := s.store.SavePromtail(l)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportSavePromtail(l, report)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *monitoringService) DeleteLoki(ctx context.Context, opts client.DeleteLokiOpts) error {
	err := s.spinner.Start("loki")
	if err != nil {
		return err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	chart := loki.New(nil)

	err = s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      lokiDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	err = s.api.DeleteLoki(api.DeleteHelmReleaseOpts{
		ID:          opts.ID,
		ReleaseName: chart.ReleaseName,
		Namespace:   chart.Namespace,
	})
	if err != nil {
		return err
	}

	report, err := s.store.RemoveLoki(opts.ID)
	if err != nil {
		return err
	}

	return s.report.ReportRemoveLoki(report)
}

// nolint: funlen
func (s *monitoringService) CreateLoki(ctx context.Context, opts client.CreateLokiOpts) (*client.Loki, error) {
	err := s.spinner.Start("loki")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	chart, err := s.api.CreateLoki(opts)
	if err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(datasource.NewLoki())
	if err != nil {
		return nil, err
	}

	_, err = s.manifest.CreateConfigMap(ctx, client.CreateConfigMapOpts{
		ID:        opts.ID,
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
			ID:        opts.ID,
			Name:      constant.DefaultKubePrometheusStackGrafanaName,
			Namespace: constant.DefaultMonitoringNamespace,
			Replicas:  replicas,
		})
		if err != nil {
			return nil, err
		}
	}

	l := &client.Loki{
		ID:    opts.ID,
		Chart: chart,
	}

	report, err := s.store.SaveLoki(l)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportSaveLoki(l, report)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// nolint: funlen gocyclo
func (s *monitoringService) DeleteKubePromStack(ctx context.Context, opts client.DeleteKubePromStackOpts) error {
	err := s.spinner.Start("kubepromstack")
	if err != nil {
		return err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	err = s.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      cloudwatchDatasourceConfigMapName,
		Namespace: constant.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	// Do we like this? Probably not.
	chart := kubepromstack.New(0*time.Second, nil)

	err = s.api.DeleteKubePromStack(api.DeleteHelmReleaseOpts{
		ID:          opts.ID,
		ReleaseName: chart.ReleaseName,
		Namespace:   chart.Namespace,
	})
	if err != nil {
		return err
	}

	err = s.manifest.DeleteExternalSecret(ctx, client.DeleteExternalSecretOpts{
		ID: opts.ID,
		Secrets: map[string]string{
			secretsCfgName: constant.DefaultMonitoringNamespace,
		},
	})
	if err != nil {
		return err
	}

	err = s.ident.DeleteIdentityPoolClient(ctx, api.DeleteIdentityPoolClientOpts{
		ID:      opts.ID,
		Purpose: grafanaPurpose,
	})
	if err != nil {
		return err
	}

	for _, secretName := range []string{clientSecretName, secretKeyName, adminUserName, adminPassName} {
		if err = s.param.DeleteSecret(ctx, api.DeleteSecretOpts{Name: secretName}); err != nil {
			return err
		}
	}

	err = s.cert.DeleteCertificate(ctx, api.DeleteCertificateOpts{
		ID:     opts.ID,
		Domain: grafanaDomain(opts.Domain),
	})
	if err != nil {
		return err
	}

	cc, err := clusterconfig.NewCloudwatchDatasourceServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		"N/A",
		constant.DefaultMonitoringNamespace,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
	)
	if err != nil {
		return err
	}

	err = s.service.DeleteServiceAccount(ctx, api.DeleteServiceAccountOpts{
		ID:     opts.ID,
		Name:   "cloudwatch-datasource", // Make this configurable
		Config: cc,
	})
	if err != nil {
		return err
	}

	err = s.policy.DeletePolicy(ctx, api.DeletePolicyOpts{
		ID:        opts.ID,
		StackName: cfn.NewStackNamer().CloudwatchDatasource(opts.ID.Repository, opts.ID.Environment),
	})
	if err != nil {
		return err
	}

	err = s.policy.DeletePolicy(ctx, api.DeletePolicyOpts{
		ID:        opts.ID,
		StackName: cfn.NewStackNamer().FargateCloudwatch(opts.ID.Repository, opts.ID.Environment),
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

	r1, err := s.store.RemoveKubePromStack(opts.ID)
	if err != nil {
		return err
	}

	r2, err := s.state.RemoveKubePromStack(opts.ID)
	if err != nil {
		return err
	}

	return s.report.ReportRemoveKubePromStack([]*store.Report{r1, r2})
}

// nolint: funlen, gocyclo
func (s *monitoringService) CreateKubePromStack(ctx context.Context, opts client.CreateKubePromStackOpts) (*client.KubePromStack, error) {
	err := s.spinner.Start("kubepromstack")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	cft, err := cfn.New(components.NewCloudwatchDatasourcePolicyComposer(opts.ID.Repository, opts.ID.Environment)).Build()
	if err != nil {
		return nil, err
	}

	policy, err := s.policy.CreatePolicy(ctx, api.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              cfn.NewStackNamer().CloudwatchDatasource(opts.ID.Repository, opts.ID.Environment),
		PolicyOutputName:       "CloudwatchDatasourcePolicy", // We need to cleanup the way we name outputs
		CloudFormationTemplate: cft,
	})
	if err != nil {
		return nil, err
	}

	cc, err := clusterconfig.NewCloudwatchDatasourceServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		policy.PolicyARN,
		constant.DefaultMonitoringNamespace,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
	)
	if err != nil {
		return nil, err
	}

	_, err = s.service.CreateServiceAccount(ctx, api.CreateServiceAccountOpts{
		ID:        opts.ID,
		Name:      "cloudwatch-datasource", // Like, why? We need to make these configurable
		PolicyArn: policy.PolicyARN,
		Config:    cc,
	})
	if err != nil {
		return nil, err
	}

	cert, err := s.cert.CreateCertificate(ctx, api.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         dns.Fqdn(grafanaDomain(opts.Domain)),
		Domain:       grafanaDomain(opts.Domain),
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, err
	}

	poolClient, err := s.ident.CreateIdentityPoolClient(ctx, api.CreateIdentityPoolClientOpts{
		ID:          opts.ID,
		UserPoolID:  opts.UserPoolID,
		Purpose:     grafanaPurpose,
		CallbackURL: fmt.Sprintf("https://%s/login/generic_oauth", cert.Domain),
	})
	if err != nil {
		return nil, err
	}

	clientSecret, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   clientSecretName,
		Secret: poolClient.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	cookieSecret, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   secretKeyName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	adminUser, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   adminUserName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	adminPass, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   adminPassName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	manifest, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID: opts.ID,
		Manifests: []api.Manifest{
			{
				Name:      secretsCfgName,
				Namespace: constant.DefaultMonitoringNamespace,
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
		},
	})
	if err != nil {
		return nil, err
	}

	chart, err := s.api.CreateKubePromStack(api.CreateKubePrometheusStackOpts{
		ID:                                  opts.ID,
		GrafanaCloudWatchServiceAccountName: constant.DefaultGrafanaCloudWatchDatasourceName,
		CertificateARN:                      cert.CertificateARN,
		Hostname:                            cert.Domain,
		AuthHostname:                        opts.AuthDomain,
		ClientID:                            poolClient.ClientID,
		SecretsConfigName:                   secretsCfgName,
		SecretsCookieSecretKey:              cookieSecret.Name,
		SecretsClientSecretKey:              clientSecret.Name,
		SecretsAdminUserKey:                 adminUser.Name,
		SecretsAdminPassKey:                 adminPass.Name,
	})
	if err != nil {
		return nil, err
	}

	stack := &client.KubePromStack{
		ID:                     opts.ID,
		CertificateARN:         cert.CertificateARN,
		Hostname:               cert.Domain,
		AuthHostname:           opts.AuthDomain,
		ClientID:               poolClient.ClientID,
		SecretsConfigName:      secretsCfgName,
		SecretsCookieSecretKey: cookieSecret.Name,
		SecretsClientSecretKey: clientSecret.Name,
		SecretsAdminUserKey:    adminUser.Name,
		SecretsAdminPassKey:    adminPass.Name,
		Chart:                  chart,
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

	fcp, err := cfn.New(components.NewFargateCloudwatchPolicyComposer(opts.ID.Repository, opts.ID.Environment)).Build()
	if err != nil {
		return nil, err
	}

	fargatePolicy, err := s.policy.CreatePolicy(ctx, api.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              cfn.NewStackNamer().FargateCloudwatch(opts.ID.Repository, opts.ID.Environment),
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

	r1, err := s.store.SaveKubePromStack(stack)
	if err != nil {
		return nil, err
	}

	r2, err := s.state.SaveKubePromStack(stack)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportSaveKubePromStack(stack, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return stack, nil
}

// NewMonitoringService returns an initialised service
func NewMonitoringService(
	spinner spinner.Spinner,
	api client.MonitoringAPI,
	store client.MonitoringStore,
	state client.MonitoringState,
	report client.MonitoringReport,
	cert client.CertificateService,
	ident client.IdentityManagerService,
	manifest client.ManifestService,
	param client.ParameterService,
	service client.ServiceAccountService,
	policy client.ManagedPolicyService,
	provider v1alpha1.CloudProvider,
) client.MonitoringService {
	return &monitoringService{
		spinner:  spinner,
		api:      api,
		store:    store,
		state:    state,
		report:   report,
		cert:     cert,
		ident:    ident,
		param:    param,
		manifest: manifest,
		service:  service,
		policy:   policy,
		provider: provider,
	}
}
