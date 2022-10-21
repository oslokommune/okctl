// Package constant contains constants used throughout okctl
package constant

import (
	"time"
)

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

	DefaultEKSKubernetesVersion = "1.21"

	DefaultChartApplyTimeout  = 10 * time.Minute
	DefaultChartRemoveTimeout = 5 * time.Minute
	DefaultChartFindTimeout   = 1 * time.Minute

	DefaultAwsRegion                       = "eu-west-1"
	DefaultGithubHost                      = "git@github.com"
	DefaultGithubOrganization              = "oslokommune"
	DefaultFargateObservabilityNamespace   = "aws-observability"
	DefaultArgoCDNamespace                 = "argocd"
	DefaultMonitoringNamespace             = "monitoring"
	DefaultKubePrometheusStackGrafanaName  = "kube-prometheus-stack-grafana"
	DefaultGrafanaCloudWatchDatasourceName = "cloudwatch-datasource"

	DefaultClusterKubeConfig     = "kubeconfig"
	DefaultClusterKubePermission = 0o600
	DefaultClusterAwsConfig      = "aws-config"
	DefaultClusterAwsCredentials = "aws-credentials"

	DefaultHelmBaseDir          = "helm"
	DefaultHelmRegistryConfig   = "registry.json"
	DefaultHelmRepositoryConfig = "repositories.yaml"
	DefaultHelmRepositoryCache  = "repository"
	DefaultHelmPluginsDirectory = "plugins"

	DefaultStormDBName                  = "state.db"
	DefaultStormNodeDomains             = "domains"
	DefaultStormNodeCertificates        = "certificates"
	DefaultStormNodeKubernetesManifest  = "manifests"
	DefaultStormNodeMonitoring          = "monitoring"
	DefaultStormNodeComponent           = "components"
	DefaultStormNodeArgoCD              = "argocd"
	DefaultStormNodeHelm                = "helm"
	DefaultStormNodeExternalDNS         = "externaldns"
	DefaultStormNodeParameter           = "parameter"
	DefaultStormNodeIdentityManager     = "identitymanager"
	DefaultStormNodeVpc                 = "vpc"
	DefaultStormNodeCluster             = "cluster"
	DefaultStormNodeGithub              = "github"
	DefaultStormNodeManagedPolicy       = "managedpolicy"
	DefaultStormNodeServiceAccount      = "serviceaccount"
	DefaultStormNodeContainerRepository = "containers"
	DefaultStormNodeUpgrade             = "upgrade"

	// EnvPrefix of environment variables that will be processed by okctl
	EnvPrefix = "OKCTL"
	// EnvHome is the default env var parsed for determining the application home
	EnvHome = "OKCTL_HOME"
	// EnvClusterDeclaration specifies the path to the cluster declaration context
	EnvClusterDeclaration = EnvPrefix + "_" + "CLUSTER_DECLARATION"
	// EnvAWSCredentialsType enables setting the --aws-credentials-type flag through environment vars
	EnvAWSCredentialsType = EnvPrefix + "_" + "AWS_CREDENTIALS_TYPE"
	// EnvGithubCredentialsType enables setting the --github-credentials-type flag through environment vars
	EnvGithubCredentialsType = EnvPrefix + "_" + "GITHUB_CREDENTIALS_TYPE"

	// DefaultKeyringServiceName is the name of the keyring or encrypted file used to store client secrets
	DefaultKeyringServiceName = "okctlService"

	// DefaultRequiredEpis number of elastic ips required for cluster creation
	DefaultRequiredEpis = 3
	// DefaultRequiredVpcs number of vpc(s) required for cluster creation
	DefaultRequiredVpcs = 1
	// DefaultRequiredIgws number of internet gateways required for cluster creation
	DefaultRequiredIgws = 1

	// DefaultRequiredFargateVCPUOnDemandResourceCount is the minimum amount of Fargate vCPU units that should be available
	DefaultRequiredFargateVCPUOnDemandResourceCount = 6

	// DefaultMaxReconciliationRequeues defines the maximum allowed times a reconciliation can be requeued
	DefaultMaxReconciliationRequeues = 3
	// DefaultReconciliationLoopDelayDuration defines the default delay between each reconciliation
	DefaultReconciliationLoopDelayDuration = 1 * time.Second

	// DefaultClusterCIDR defines the default CIDR to use when creating cluster VPCs
	DefaultClusterCIDR     = "192.168.0.0/20"
	DefaultOutputDirectory = "infrastructure"

	DefaultAutomaticPullRequestMergeLabel = "automerge"

	DefaultAwsAuthDocumentationURL = "https://www.okctl.io/authenticating-to-aws/"

	ValidationOKUsername = "^[a-z]{3}[0-9]{4,7}$"
)
