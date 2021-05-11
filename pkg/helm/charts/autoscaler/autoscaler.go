// Package autoscaler provides a Helm chart for installing:
// - https://github.com/kubernetes/autoscaler/tree/master/charts/cluster-autoscaler
package autoscaler

import (
	"bytes"
	"text/template"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

const (
	// ReleaseName is the name of the release
	ReleaseName = "cluster-autoscaler"
	// Namespace is the default namespace
	Namespace = "kube-system"
)

// New returns an initialised Helm chart for installing cluster-autoscaler
func New(values *Values, timeout time.Duration) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "autoscaler",
		RepositoryURL:  "https://kubernetes.github.io/autoscaler",
		ReleaseName:    ReleaseName,
		Version:        "9.4.0",
		Chart:          "cluster-autoscaler",
		Namespace:      Namespace,
		Timeout:        timeout,
		Values:         values,
	}
}

// NewDefaultValues returns the mapped values.yml containing
// the default values
func NewDefaultValues(region, clusterName, serviceAccount string) *Values {
	return &Values{
		Region:         region,
		ClusterName:    clusterName,
		ServiceAccount: serviceAccount,
	}
}

// Values contains the required inputs for generating the values.yml
type Values struct {
	Region         string
	ClusterName    string
	ServiceAccount string
}

// RawYAML implements the raw marshaller interface in the Helm package
func (v *Values) RawYAML() ([]byte, error) {
	tmpl, err := template.New("values").Parse(valuesTemplate)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer

	err = tmpl.Execute(&buff, *v)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// nolint: lll
const valuesTemplate = `## Ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
# affinity -- Affinity for pod assignment
affinity: {}

autoDiscovery:
  # cloudProviders "aws", "gce" and "magnum" are supported by auto-discovery at this time
  # AWS: Set tags as described in https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md#auto-discovery-setup

  # autoDiscovery.clusterName -- Enable autodiscovery for "cloudProvider=aws", for groups matching "autoDiscovery.tags".
  # Enable autodiscovery for "cloudProvider=gce", but no MIG tagging required.
  # Enable autodiscovery for "cloudProvider=magnum", for groups matching "autoDiscovery.roles".
  clusterName:  {{.ClusterName}}

  # autoDiscovery.tags -- ASG tags to match, run through "tpl".
  tags:
  - k8s.io/cluster-autoscaler/enabled
  - k8s.io/cluster-autoscaler/{{.ClusterName}}

  # autoDiscovery.roles -- Magnum node group roles to match.
  roles:
  - worker

# autoscalingGroups -- For AWS, Azure AKS or Magnum. At least one element is required if not using "autoDiscovery". For example:
# <pre>
# - name: asg1<br />
#   maxSize: 2<br />
#   minSize: 1
# </pre>
autoscalingGroups: []
# - name: asg1
#   maxSize: 2
#   minSize: 1
# - name: asg2
#   maxSize: 2
#   minSize: 1

# autoscalingGroupsnamePrefix -- For GCE. At least one element is required if not using "autoDiscovery". For example:
# <pre>
# - name: ig01<br />
#   maxSize: 10<br />
#   minSize: 0
# </pre>
autoscalingGroupsnamePrefix: []
# - name: ig01
#   maxSize: 10
#   minSize: 0
# - name: ig02
#   maxSize: 10
#   minSize: 0

# awsAccessKeyID -- AWS access key ID ([if AWS user keys used](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md#using-aws-credentials))
awsAccessKeyID: ""

# awsRegion -- AWS region (required if "cloudProvider=aws")
awsRegion: {{.Region}}

# awsSecretAccessKey -- AWS access secret key ([if AWS user keys used](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md#using-aws-credentials))
awsSecretAccessKey: ""

# azureClientID -- Service Principal ClientID with contributor permission to Cluster and Node ResourceGroup.
# Required if "cloudProvider=azure"
azureClientID: ""

# azureClientSecret -- Service Principal ClientSecret with contributor permission to Cluster and Node ResourceGroup.
# Required if "cloudProvider=azure"
azureClientSecret: ""

# azureResourceGroup -- Azure resource group that the cluster is located.
# Required if "cloudProvider=azure"
azureResourceGroup: ""

# azureSubscriptionID -- Azure subscription where the resources are located.
# Required if "cloudProvider=azure"
azureSubscriptionID: ""

# azureTenantID -- Azure tenant where the resources are located.
# Required if "cloudProvider=azure"
azureTenantID: ""

# azureVMType -- Azure VM type.
azureVMType: "AKS"

# azureClusterName -- Azure AKS cluster name.
# Required if "cloudProvider=azure"
azureClusterName: ""

# azureNodeResourceGroup -- Azure resource group where the cluster's nodes are located, typically set as "MC_<cluster-resource-group-name>_<cluster-name>_<location>".
# Required if "cloudProvider=azure"
azureNodeResourceGroup: ""

# azureUseManagedIdentityExtension -- Whether to use Azure's managed identity extension for credentials. If using MSI, ensure subscription ID and resource group are set.
azureUseManagedIdentityExtension: false

# magnumClusterName -- Cluster name or ID in Magnum.
# Required if "cloudProvider=magnum" and not setting "autoDiscovery.clusterName".
magnumClusterName: ""

# magnumCABundlePath -- Path to the host's CA bundle, from "ca-file" in the cloud-config file.
magnumCABundlePath: "/etc/kubernetes/ca-bundle.crt"

# cloudConfigPath -- Configuration file for cloud provider.
cloudConfigPath: /etc/gce.conf

# cloudProvider -- The cloud provider where the autoscaler runs.
# Currently only "gce", "aws", "azure" and "magnum" are supported.
# "aws" supported for AWS. "gce" for GCE. "azure" for Azure AKS.
# "magnum" for OpenStack Magnum.
cloudProvider: aws

# containerSecurityContext -- [Security context for container](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
containerSecurityContext: {}
  # capabilities:
  #   drop:
  #   - ALL

# dnsPolicy -- Defaults to "ClusterFirst". Valid values are:
# "ClusterFirstWithHostNet", "ClusterFirst", "Default" or "None".
# If autoscaler does not depend on cluster DNS, recommended to set this to "Default".
dnsPolicy: ClusterFirst

## Priorities Expander
# expanderPriorities -- The expanderPriorities is used if "extraArgs.expander" is set to "priority" and expanderPriorities is also set with the priorities.
# If "extraArgs.expander" is set to "priority", then expanderPriorities is used to define cluster-autoscaler-priority-expander priorities.
# See: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/expander/priority/readme.md
expanderPriorities: {}

# extraArgs -- Additional container arguments.
# Refer to https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-are-the-parameters-to-ca for the full list of cluster autoscaler
# parameters and their default values.
extraArgs:
  logtostderr: true
  stderrthreshold: info
  v: 4
  # write-status-configmap: true
  # leader-elect: true
  # skip-nodes-with-local-storage: true
  # expander: random
  # scale-down-enabled: true
  # balance-similar-node-groups: true
  # min-replica-count: 0
  # scale-down-utilization-threshold: 0.5
  # scale-down-non-empty-candidates-count: 30
  # max-node-provision-time: 15m0s
  # scan-interval: 10s
  # scale-down-delay-after-add: 10m
  # scale-down-delay-after-delete: 0s
  # scale-down-delay-after-failure: 3m
  # scale-down-unneeded-time: 10m
  # skip-nodes-with-system-pods: true

# extraEnv -- Additional container environment variables.
extraEnv: {}

# extraEnvConfigMaps -- Additional container environment variables from ConfigMaps.
extraEnvConfigMaps: {}

# extraEnvSecrets -- Additional container environment variables from Secrets.
extraEnvSecrets: {}

# envFromConfigMap -- ConfigMap name to use as envFrom.
envFromConfigMap: ""

# envFromSecret -- Secret name to use as envFrom.
envFromSecret: ""

# extraVolumeSecrets -- Additional volumes to mount from Secrets.
extraVolumeSecrets: {}
  # autoscaler-vol:
  #   mountPath: /data/autoscaler/
  # custom-vol:
  #   name: custom-secret
  #   mountPath: /data/custom/
  #   items:
  #     - key: subkey
  #       path: mypath

# extraVolumes -- Additional volumes.
extraVolumes: []
  # - name: ssl-certs
  #   hostPath:
  #     path: /etc/ssl/certs/ca-bundle.crt

# extraVolumeMounts -- Additional volumes to mount.
extraVolumeMounts: []
  # - name: ssl-certs
  #   mountPath: /etc/ssl/certs/ca-certificates.crt
  #   readonly: true

# fullnameOverride -- String to fully override "cluster-autoscaler.fullname" template.
fullnameOverride: ""

image:
  # image.repository -- Image repository
  repository: us.gcr.io/k8s-artifacts-prod/autoscaling/cluster-autoscaler
  # image.tag -- Image tag
  tag: v1.18.1
  # image.pullPolicy -- Image pull policy
  pullPolicy: IfNotPresent
  ## Optionally specify an array of imagePullSecrets.
  ## Secrets must be manually created in the namespace.
  ## ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
  ##
  # image.pullSecrets -- Image pull secrets
  pullSecrets: []
  # - myRegistrKeySecretName

# kubeTargetVersionOverride -- Allow overriding the ".Capabilities.KubeVersion.GitVersion" check. Useful for "helm template" commands.
kubeTargetVersionOverride: ""

# nameOverride -- String to partially override "cluster-autoscaler.fullname" template (will maintain the release name)
nameOverride: ""

# nodeSelector -- Node labels for pod assignment. Ref: https://kubernetes.io/docs/user-guide/node-selection/.
nodeSelector: {}

# podAnnotations -- Annotations to add to each pod.
podAnnotations: {}

# podDisruptionBudget -- Pod disruption budget.
podDisruptionBudget:
  maxUnavailable: 1
  # minAvailable: 2

# podLabels -- Labels to add to each pod.
podLabels: {}

# additionalLabels -- Labels to add to each object of the chart.
additionalLabels: {}

# priorityClassName -- priorityClassName
priorityClassName: ""

rbac:
  # rbac.create -- If "true", create and use RBAC resources.
  create: true
  # rbac.pspEnabled -- If "true", creates and uses RBAC resources required in the cluster with [Pod Security Policies](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) enabled.
  # Must be used with "rbac.create" set to "true".
  pspEnabled: false
  serviceAccount:
    # rbac.serviceAccount.annotations -- Additional Service Account annotations.
    annotations: {}
    # rbac.serviceAccount.create -- If "true" and "rbac.create" is also true, a Service Account will be created.
    create: false
    # rbac.serviceAccount.name -- The name of the ServiceAccount to use. If not set and create is "true", a name is generated using the fullname template.
    name: {{.ServiceAccount}}

# replicaCount -- Desired number of pods
replicaCount: 1

# resources -- Pod resource requests and limits.
resources: {}
  # limits:
  #   cpu: 100m
  #   memory: 300Mi
  # requests:
  #   cpu: 100m
  #   memory: 300Mi

# securityContext -- [Security context for pod](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
securityContext: {}
  # runAsNonRoot: true
  # runAsUser: 1001
  # runAsGroup: 1001

service:
  # service.annotations -- Annotations to add to service
  annotations: {}
  # service.labels -- Labels to add to service
  labels: {}
  # service.externalIPs -- List of IP addresses at which the service is available. Ref: https://kubernetes.io/docs/user-guide/services/#external-ips.
  externalIPs: []

  # service.loadBalancerIP -- IP address to assign to load balancer (if supported).
  loadBalancerIP: ""
  # service.loadBalancerSourceRanges -- List of IP CIDRs allowed access to load balancer (if supported).
  loadBalancerSourceRanges: []
  # service.servicePort -- Service port to expose.
  servicePort: 8085
  # service.portName -- Name for service port.
  portName: http
  # service.type -- Type of service to create.
  type: ClusterIP

## Are you using Prometheus Operator?
serviceMonitor:
  # serviceMonitor.enabled -- If true, creates a Prometheus Operator ServiceMonitor.
  enabled: false
  # serviceMonitor.interval -- Interval that Prometheus scrapes Cluster Autoscaler metrics.
  interval: 10s
  # serviceMonitor.namespace -- Namespace which Prometheus is running in.
  namespace: monitoring
  ## [Prometheus Selector Label](https://github.com/helm/charts/tree/master/stable/prometheus-operator#prometheus-operator-1)
  ## [Kube Prometheus Selector Label](https://github.com/helm/charts/tree/master/stable/prometheus-operator#exporters)
  # serviceMonitor.selector -- Default to kube-prometheus install (CoreOS recommended), but should be set according to Prometheus install.
  selector:
    release: prometheus-operator
  # serviceMonitor.path -- The path to scrape for metrics; autoscaler exposes "/metrics" (this is standard)
  path: /metrics

# tolerations -- List of node taints to tolerate (requires Kubernetes >= 1.6).
tolerations: []

# updateStrategy -- [Deployment update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy)
updateStrategy: {}
  # rollingUpdate:
  #   maxSurge: 1
  #   maxUnavailable: 0
  # type: RollingUpdate
`
