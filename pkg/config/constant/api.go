// Package constant contains constants used throughout okctl
package constant

import "time"

// nolint: golint
const (
	// DefaultDir is the default location directory for the okctl application config
	DefaultDir = ".okctl"
	// DefaultConfig is the default filename of the okctl application config
	DefaultConfig = "conf.yml"
	// DefaultConfigName is the default name of the okctl application config
	DefaultConfigName = "conf"
	// DefaultConfigType is the default type of the okctl application config
	DefaultConfigType = "yml"
	// DefaultLogDir is the default directory name for logs
	DefaultLogDir = "logs"
	// DefaultLogName is the default name of the file to log to
	DefaultLogName = "console.log"
	// DefaultLogDays determines how many days we keep the logs
	DefaultLogDays = 28
	// DefaultLogBackups determines how many backups we will keep
	DefaultLogBackups = 3
	// DefaultLogSizeInMb determines how much storage we will consume
	DefaultLogSizeInMb = 50
	// DefaultCredentialsDirName sets the name of the directory for creds
	DefaultCredentialsDirName = "credentials"

	DefaultRepositoryStateFile = ".okctl.yml"

	DefaultEKSKubernetesVersion = "1.19"

	DefaultChartApplyTimeout  = 5 * time.Minute
	DefaultChartRemoveTimeout = 5 * time.Minute

	DefaultGithubHost                      = "git@github.com"
	DefaultGithubOrganization              = "oslokommune"
	DefaultFargateObservabilityNamespace   = "aws-observability"
	DefaultArgoCDNamespace                 = "argocd"
	DefaultMonitoringNamespace             = "monitoring"
	DefaultKubePrometheusStackGrafanaName  = "kube-prometheus-stack-grafana"
	DefaultGrafanaCloudWatchDatasourceName = "cloudwatch-datasource"

	DefaultClusterConfig         = "cluster.yml"
	DefaultClusterKubeConfig     = "kubeconfig"
	DefaultClusterAwsConfig      = "aws-config"
	DefaultClusterAwsCredentials = "aws-credentials"
	DefaultClusterBaseDir        = "cluster"

	DefaultHelmBaseDir          = "helm"
	DefaultHelmRegistryConfig   = "registry.json"
	DefaultHelmRepositoryConfig = "repositories.yaml"
	DefaultHelmRepositoryCache  = "repository"
	DefaultHelmPluginsDirectory = "plugins"

	DefaultVpcOutputs                = "vpc-outputs.json"
	DefaultVpcCloudFormationTemplate = "vpc-cf.yml"
	DefaultVpcBaseDir                = "vpc"

	DefaultMonitoringBaseDir                        = "monitoring"
	DefaultPromtailBaseDir                          = "promtail"
	DefaultLokiBaseDir                              = "loki"
	DefaultTempoBaseDir                             = "tempo"
	DefaultAWSLoadBalancerControllerBaseDir         = "aws-load-balancer-controller"
	DefaultAliasBaseDir                             = "alias"
	DefaultAliasCloudFormationTemplate              = "alias-cf.yaml"
	DefaultArgoCDBaseDir                            = "argocd"
	DefaultArgoOutputsFile                          = "argocd-outputs.json"
	DefaultCertificateBaseDir                       = "certificates"
	DefaultCertificateCloudFormationTemplate        = "certificate-cf.yml"
	DefaultDomainBaseDir                            = "domains"
	DefaultDomainCloudFormationTemplate             = "domains-cf.yml"
	DefaultExternalDNSBaseDir                       = "external-dns"
	DefaultExternalSecretsBaseDir                   = "external-secrets"
	DefaultAutoscalerBaseDir                        = "autoscaler"
	DefaultBlockstorageBaseDir                      = "blockstorage"
	DefaultKubePromStackBaseDir                     = "kubepromstack"
	DefaultKubePromStackOutputsFile                 = "kubepromstack-outputs.json"
	DefaultHelmChartFile                            = "helm-chart.json"
	DefaultHelmOutputsFile                          = "helm-outputs.json"
	DefaultHelmReleaseFile                          = "helm-release.json"
	DefaultIdentityPoolBaseDir                      = "identitypool"
	DefaultIdentityPoolCloudFormationTemplate       = "identitypool-cf.yaml"
	DefaultIdentityPoolOutputsFile                  = "identitypool-outputs.json"
	DefaultIdentityPoolClientsBaseDir               = "clients"
	DefaultIdentityPoolClientCloudFormationTemplate = "ipc-cf.yaml"
	DefaultIdentityPoolClientOutputsFile            = "ipc-outputs.json"
	DefaultIdentityPoolUsersBaseDir                 = "users"
	DefaultIdentityPoolUserOutputsFile              = "ipu-outputs.json"
	DefaultIdentityPoolUserCloudFormationTemplate   = "ipu-cf.yaml"
	DefaultKubeOutputsFile                          = "kube-outputs.json"
	DefaultParameterBaseDir                         = "parameters"
	DefaultParameterOutputsFile                     = "parameter-outputs.json"
	DefaultPolicyCloudFormationTemplateFile         = "policy-cf.yml"
	DefaultPolicyOutputFile                         = "policy-outputs.json"
	DefaultServiceAccountConfigFile                 = "service-account-config.yml"
	DefaultServiceAccountOutputsFile                = "service-account-outputs.json"

	DefaultStormDBName                       = "state.db"
	DefaultStormNodeDomains                  = "domains"
	DefaultStormNodeCertificates             = "certificates"
	DefaultStormNodeServiceAccounts          = "serviceaccounts"
	DefaultStormNodeKubernetesManifest       = "manifests"
	DefaultStormNodeBlockStorage             = "blockstorage"
	DefaultStormNodeMonitoring               = "monitoring"
	DefaultStormNodeComponent                = "components"
	DefaultStormNodeArgoCD                   = "argocd"
	DefaultStormNodeHelm                     = "helm"
	DefaultStormNodeExternalSecrets          = "externalsecrets"
	DefaultStormNodeAWSLoadBalanerController = "awsloadbalancercontroller"

	// EnvPrefix of environment variables that will be processed by okctl
	EnvPrefix = "OKCTL"
	// EnvHome is the default env var parsed for determining the application home
	EnvHome = "OKCTL_HOME"

	// DefaultApplicationsOutputDir is where the application declarations reside
	DefaultApplicationsOutputDir = "applications"
	// DefaultApplicationBaseDir is where the directory where application base files reside
	DefaultApplicationBaseDir = "base"
	// DefaultApplicationOverlayDir is where the directory where application overlay files reside
	DefaultApplicationOverlayDir = "overlays"
	// DefaultIngressPatchFilename defines the filename of the ingress patch
	DefaultIngressPatchFilename = "ingress-patch.json"

	// DefaultKeyringServiceName is the name of the keyring or encrypted file used to store client secrets
	DefaultKeyringServiceName = "okctlService"

	// DefaultRequiredEpis number of elastic ips required for cluster creation
	DefaultRequiredEpis = 3
	// DefaultRequiredVpcs number of vpc(s) required for cluster creation
	DefaultRequiredVpcs = 1
	// DefaultRequiredIgws number of internet gateways required for cluster creation
	DefaultRequiredIgws = 1

	DefaultNameserverRecordTTL = 300

	// DefaultMaxReconciliationRequeues defines the maximum allowed times a reconciliation can be requeued
	DefaultMaxReconciliationRequeues = 3

	// DefaultClusterCIDR defines the default CIDR to use when creating cluster VPCs
	DefaultClusterCIDR     = "192.168.0.0/20"
	DefaultOutputDirectory = "infrastructure"

	DefaultComponentBaseDir           = "component"
	DefaultPostgresBaseDir            = "postgres"
	DefaultPostgresOutputFile         = "postgres-output.json"
	DefaultPostgresCloudFormationFile = "postgres-cf.yaml"
)
