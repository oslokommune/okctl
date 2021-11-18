// Package argocd provides a Helm chart for installing:
// - https://argoproj.github.io/argo-cd/
// - https://github.com/argoproj/argo-helm
package argocd

import (
	"bytes"
	"text/template"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

const (
	// ReleaseName is the name of the release
	ReleaseName = "argocd"
	// Namespace is the default namespace
	Namespace = "argocd"
)

// New returns an initialised Helm chart
func New(values *Values, timeout time.Duration) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "argo",
		RepositoryURL:  "https://argoproj.github.io/argo-helm",
		ReleaseName:    ReleaseName,
		Version:        "3.26.9",
		Chart:          "argo-cd",
		Namespace:      Namespace,
		Timeout:        timeout,
		Values:         values,
	}
}

// ValuesOpts contains input for creating the default values
type ValuesOpts struct {
	URL                  string
	HostName             string
	CertificateARN       string
	Region               string
	ClientID             string
	Organisation         string
	AuthDomain           string
	UserPoolID           string
	RepoURL              string
	RepoName             string
	PrivateKeySecretName string
	PrivateKeySecretKey  string
}

// NewDefaultValues returns the default values for the chart
//nolint:gomnd,funlen,lll
func NewDefaultValues(opts ValuesOpts) *Values {
	return &Values{
		URL:                  opts.URL,
		HostName:             opts.HostName,
		CertificateARN:       opts.CertificateARN,
		Region:               opts.Region,
		ClientID:             opts.ClientID,
		Organisation:         opts.Organisation,
		AuthDomain:           opts.AuthDomain,
		UserPoolID:           opts.UserPoolID,
		RepoURL:              opts.RepoURL,
		RepoName:             opts.RepoName,
		PrivateKeySecretName: opts.PrivateKeySecretName,
		PrivateKeySecretKey:  opts.PrivateKeySecretKey,
	}
}

// Values contains the parameters we map up
type Values struct {
	URL                  string
	HostName             string
	CertificateARN       string
	Region               string
	ClientID             string
	Organisation         string
	AuthDomain           string
	UserPoolID           string
	RepoURL              string
	RepoName             string
	PrivateKeySecretName string
	PrivateKeySecretKey  string
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

const valuesTemplate = `
## ArgoCD configuration
## Ref: https://github.com/argoproj/argo-cd
##

# -- Provide a name in place of argocd
nameOverride: argocd
# -- String to fully override "argo-cd.fullname"
fullnameOverride: ""
# -- Override the Kubernetes version, which is used to evaluate certain manifests
kubeVersionOverride: ""

global:
  image:
    # -- If defined, a repository applied to all ArgoCD deployments
    repository: quay.io/argoproj/argocd
    # -- Overrides the global ArgoCD image tag whose default is the chart appVersion
    tag: ""
    # -- If defined, a imagePullPolicy applied to all ArgoCD deployments
    imagePullPolicy: IfNotPresent
  # -- Annotations for the all deployed pods
  podAnnotations: {}
  # -- Labels for the all deployed pods
  podLabels: {}
  # -- Toggle and define securityContext. See [values.yaml]
  securityContext: {}
  #  runAsUser: 999
  #  runAsGroup: 999
  #  fsGroup: 999

  # -- If defined, uses a Secret to pull an image from a private Docker registry or repository
  imagePullSecrets: []
  # -- Mapping between IP and hostnames that will be injected as entries in the pod's hosts files
  hostAliases: []
  # - ip: 10.20.30.40
  #   hostnames:
  #   - git.myhostname

  networkPolicy:
    # -- Create NetworkPolicy objects for all components
    create: false
    # -- Default deny all ingress traffic
    defaultDenyIngress: false

# Override APIVersions
# If you want to template helm charts but cannot access k8s API server
# you can set api versions here
apiVersionOverrides:
  # -- String to override apiVersion of certmanager resources rendered by this helm chart
  certmanager: "" # cert-manager.io/v1
  # -- String to override apiVersion of ingresses rendered by this helm chart
  ingress: "" # networking.k8s.io/v1beta1

# -- Create clusterroles that extend existing clusterroles to interact with argo-cd crds
## Ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#aggregated-clusterroles
createAggregateRoles: false

## Controller
controller:
  # -- Application controller name string
  name: application-controller

  image:
    # -- Repository to use for the application controller
    # @default -- "" (defaults to global.image.repository)
    repository: ""
    # -- Tag to use for the application controller
    # @default -- "" (defaults to global.image.tag)
    tag: ""
    # -- Image pull policy for the application controller
    # @default -- "" (defaults to global.image.imagePullPolicy)
    imagePullPolicy: ""

  # -- The number of application controller pods to run.
  # If changing the number of replicas you must pass the number as ARGOCD_CONTROLLER_REPLICAS as an environment variable
  replicas: 1

  # -- Deploy the application controller as a StatefulSet instead of a Deployment, this is required for HA capability.
  # This is a feature flag that will become the default in chart version 3.x
  enableStatefulSet: false

  ## Application controller commandline flags
  args:
    # -- define the application controller --status-processors
    statusProcessors: "20"
    # -- define the application controller --operation-processors
    operationProcessors: "10"
    # -- define the application controller --app-resync
    appResyncPeriod: "180"
    # -- define the application controller --self-heal-timeout-seconds
    selfHealTimeout: "5"
    # -- define the application controller --repo-server-timeout-seconds
    repoServerTimeoutSeconds: "60"

  # -- Application controller log format. Either text or json
  logFormat: text
  # -- Application controller log level
  logLevel: info

  # -- Additional command line arguments to pass to application controller
  extraArgs: []

  # -- Environment variables to pass to application controller
  env:
    []
    # - name: "ARGOCD_CONTROLLER_REPLICAS"
    #   value: ""

  # -- envFrom to pass to application controller
  # @default -- [] (See [values.yaml])
  envFrom: []
  # - configMapRef:
  #     name: config-map-name
  # - secretRef:
  #     name: secret-name

  # -- Annotations to be added to application controller pods
  podAnnotations: {}

  # -- Labels to be added to application controller pods
  podLabels: {}

  # -- Application controller container-level security context
  containerSecurityContext:
    {}
    # capabilities:
    #   drop:
    #     - all
    # readOnlyRootFilesystem: true
    # runAsNonRoot: true

  # -- Application controller listening port
  containerPort: 8082

  ## Readiness and liveness probes for default backend
  ## Ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
  ##
  readinessProbe:
    # -- Minimum consecutive failures for the [probe] to be considered failed after having succeeded
    failureThreshold: 3
    # -- Number of seconds after the container has started before [probe] is initiated
    initialDelaySeconds: 10
    # -- How often (in seconds) to perform the [probe]
    periodSeconds: 10
    # -- Minimum consecutive successes for the [probe] to be considered successful after having failed
    successThreshold: 1
    # -- Number of seconds after which the [probe] times out
    timeoutSeconds: 1
  livenessProbe:
    # -- Minimum consecutive failures for the [probe] to be considered failed after having succeeded
    failureThreshold: 3
    # -- Number of seconds after the container has started before [probe] is initiated
    initialDelaySeconds: 10
    # -- How often (in seconds) to perform the [probe]
    periodSeconds: 10
    # -- Minimum consecutive successes for the [probe] to be considered successful after having failed
    successThreshold: 1
    # -- Number of seconds after which the [probe] times out
    timeoutSeconds: 1

  # -- Additional volumeMounts to the application controller main container
  volumeMounts: []

  # -- Additional volumes to the application controller pod
  volumes: []

  ## Controller service configuration
  service:
    # -- Application controller service annotations
    annotations: {}
    # -- Application controller service labels
    labels: {}
    # -- Application controller service port
    port: 8082
    # -- Application controller service port name
    portName: https-controller

  # -- [Node selector]
  nodeSelector: {}
  # -- [Tolerations] for use with node taints
  tolerations: []
  # -- Assign custom [affinity] rules to the deployment
  affinity: {}

  # -- Assign custom [TopologySpreadConstraints] rules to the application controller
  ## Ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
  ## If labelSelector is left out, it will default to the labelSelector configuration of the deployment
  topologySpreadConstraints: []
  # - maxSkew: 1
  #   topologyKey: topology.kubernetes.io/zone
  #   whenUnsatisfiable: DoNotSchedule

  # -- Priority class for the application controller pods
  priorityClassName: ""

  # -- Resource limits and requests for the application controller pods
  resources: {}
  #  limits:
  #    cpu: 500m
  #    memory: 512Mi
  #  requests:
  #    cpu: 250m
  #    memory: 256Mi

  serviceAccount:
    # -- Create a service account for the application controller
    create: true
    # -- Service account name
    name: argocd-application-controller
    # -- Annotations applied to created service account
    annotations: {}
    # -- Automount API credentials for the Service Account
    automountServiceAccountToken: true

  ## Application controller metrics configuration
  metrics:
    # -- Deploy metrics service
    enabled: false
    service:
      # -- Metrics service annotations
      annotations: {}
      # -- Metrics service labels
      labels: {}
      # -- Metrics service port
      servicePort: 8082
    serviceMonitor:
      # -- Enable a prometheus ServiceMonitor
      enabled: false
      # -- Prometheus ServiceMonitor interval
      interval: 30s
      # -- Prometheus [RelabelConfigs] to apply to samples before scraping
      relabelings: []
      # -- Prometheus [MetricRelabelConfigs] to apply to samples before ingestion
      metricRelabelings: []
      # -- Prometheus ServiceMonitor selector
      selector: {}
        # prometheus: kube-prometheus

      # -- Prometheus ServiceMonitor namespace
      namespace: "" # "monitoring"
      # -- Prometheus ServiceMonitor labels
      additionalLabels: {}
    rules:
      # -- Deploy a PrometheusRule for the application controller
      enabled: false
      # -- PrometheusRule.Spec for the application controller
      spec: []
      # - alert: ArgoAppMissing
      #   expr: |
      #     absent(argocd_app_info)
      #   for: 15m
      #   labels:
      #     severity: critical
      #   annotations:
      #     summary: "[ArgoCD] No reported applications"
      #     description: >
      #       ArgoCD has not reported any applications data for the past 15 minutes which
      #       means that it must be down or not functioning properly.  This needs to be
      #       resolved for this cloud to continue to maintain state.
      # - alert: ArgoAppNotSynced
      #   expr: |
      #     argocd_app_info{sync_status!="Synced"} == 1
      #   for: 12h
      #   labels:
      #     severity: warning
      #   annotations:
      #     summary: "[{{{{$labels.name}}}}] Application not synchronized"
      #     description: >
      #       The application [{{{{$labels.name}}}} has not been synchronized for over
      #       12 hours which means that the state of this cloud has drifted away from the
      #       state inside Git.
    #   selector:
    #     prometheus: kube-prometheus
    #   namespace: monitoring
    #   additionalLabels: {}

  ## Enable if you would like to grant rights to ArgoCD to deploy to the local Kubernetes cluster.
  clusterAdminAccess:
    # -- Enable RBAC for local cluster deployments
    enabled: true

  ## Enable this and set the rules: to whatever custom rules you want for the Cluster Role resource.
  ## Defaults to off
  clusterRoleRules:
    # -- Enable custom rules for the application controller's ClusterRole resource
    enabled: false
    # -- List of custom rules for the application controller's ClusterRole resource
    rules: []

  # -- Additional containers to be added to the application controller pod
  extraContainers: []

## Dex
dex:
  # -- Enable dex
  enabled: true
  # -- Dex name
  name: dex-server

  metrics:
    # -- Deploy metrics service
    enabled: false
    service:
      # -- Metrics service annotations
      annotations: {}
      # -- Metrics service labels
      labels: {}
    serviceMonitor:
      # -- Enable a prometheus ServiceMonitor
      enabled: false
      # -- Prometheus ServiceMonitor interval
      interval: 30s
      # -- Prometheus [RelabelConfigs] to apply to samples before scraping
      relabelings: []
      # -- Prometheus [MetricRelabelConfigs] to apply to samples before ingestion
      metricRelabelings: []
      # -- Prometheus ServiceMonitor selector
      selector: {}
        # prometheus: kube-prometheus

      # -- Prometheus ServiceMonitor namespace
      namespace: "" # "monitoring"
      # -- Prometheus ServiceMonitor labels
      additionalLabels: {}

  image:
    # -- Dex image repository
    repository: ghcr.io/dexidp/dex
    # -- Dex image tag
    tag: v2.30.0
    # -- Dex imagePullPolicy
    imagePullPolicy: IfNotPresent
  initImage:
    # -- Argo CD init image repository
    # @default -- "" (defaults to global.image.repository)
    repository: ""
    # -- Argo CD init image tag
    # @default -- "" (defaults to global.image.tag)
    tag: ""
    # -- Argo CD init image imagePullPolicy
    # @default -- "" (defaults to global.image.imagePullPolicy)
    imagePullPolicy: ""

  # -- Environment variables to pass to the Dex server
  env: []

  # -- envFrom to pass to the Dex server
  # @default -- [] (See [values.yaml])
  envFrom: []
  # - configMapRef:
  #     name: config-map-name
  # - secretRef:
  #     name: secret-name

  # -- Annotations to be added to the Dex server pods
  podAnnotations: {}

  # -- Labels to be added to the Dex server pods
  podLabels: {}

  ## Probes for Dex server
  ## Supported from Dex >= 2.28.0
  livenessProbe:
    # -- Enable Kubernetes liveness probe for Dex >= 2.28.0
    enabled: false
    # -- Minimum consecutive failures for the [probe] to be considered failed after having succeeded
    failureThreshold: 3
    # -- Number of seconds after the container has started before [probe] is initiated
    initialDelaySeconds: 10
    # -- How often (in seconds) to perform the [probe]
    periodSeconds: 10
    # -- Minimum consecutive successes for the [probe] to be considered successful after having failed
    successThreshold: 1
    # -- Number of seconds after which the [probe] times out
    timeoutSeconds: 1
  readinessProbe:
    # -- Enable Kubernetes readiness probe for Dex >= 2.28.0
    enabled: false
    # -- Minimum consecutive failures for the [probe] to be considered failed after having succeeded
    failureThreshold: 3
    # -- Number of seconds after the container has started before [probe] is initiated
    initialDelaySeconds: 10
    # -- How often (in seconds) to perform the [probe]
    periodSeconds: 10
    # -- Minimum consecutive successes for the [probe] to be considered successful after having failed
    successThreshold: 1
    # -- Number of seconds after which the [probe] times out
    timeoutSeconds: 1

  serviceAccount:
    # -- Create dex service account
    create: true
    # -- Dex service account name
    name: argocd-dex-server
    # -- Annotations applied to created service account
    annotations: {}
    # -- Automount API credentials for the Service Account
    automountServiceAccountToken: true

  # -- Additional volumeMounts to the dex main container
  volumeMounts:
    - name: static-files
      mountPath: /shared

  # -- Additional volumes to the dex pod
  volumes:
    - name: static-files
      emptyDir: {}

  # -- Extra volumes to the dex pod
  extraVolumes: []

  # -- Extra volumeMounts to the dex pod
  extraVolumeMounts: []

  # -- Container port for HTTP access
  containerPortHttp: 5556
  # -- Service port for HTTP access
  servicePortHttp: 5556
  # -- Service port name for HTTP access
  servicePortHttpName: http
  # -- Container port for gRPC access
  containerPortGrpc: 5557
  # -- Service port for gRPC access
  servicePortGrpc: 5557
  # -- Service port name for gRPC access
  servicePortGrpcName: grpc
  # -- Container port for metrics access
  containerPortMetrics: 5558
  # -- Service port for metrics access
  servicePortMetrics: 5558

  # -- [Node selector]
  nodeSelector: {}
  # -- [Tolerations] for use with node taints
  tolerations: []
  # -- Assign custom [affinity] rules to the deployment
  affinity: {}

  # -- Assign custom [TopologySpreadConstraints] rules to dex
  ## Ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
  ## If labelSelector is left out, it will default to the labelSelector configuration of the deployment
  topologySpreadConstraints: []
  # - maxSkew: 1
  #   topologyKey: topology.kubernetes.io/zone
  #   whenUnsatisfiable: DoNotSchedule

  # -- Priority class for dex
  priorityClassName: ""

  # -- Dex container-level security context
  containerSecurityContext:
    {}
    # capabilities:
    #   drop:
    #     - all
    # readOnlyRootFilesystem: true

# -- Resource limits and requests for dex
  resources: {}
  #  limits:
  #    cpu: 50m
  #    memory: 64Mi
  #  requests:
  #    cpu: 10m
  #    memory: 32Mi

  # -- Additional containers to be added to the dex pod
  extraContainers: []

## Redis
redis:
  # -- Enable redis
  enabled: true
  # -- Redis name
  name: redis

  image:
    # -- Redis repository
    repository: redis
    # -- Redis tag
    tag: 6.2.4-alpine
    # -- Redis imagePullPolicy
    imagePullPolicy: IfNotPresent

  # -- Additional command line arguments to pass to redis-server
  extraArgs: []
  # - --bind
  # - "0.0.0.0"

  # -- Redis container port
  containerPort: 6379
  # -- Redis service port
  servicePort: 6379

  # -- Environment variables to pass to the Redis server
  env: []

  # -- envFrom to pass to the Redis server
  # @default -- [] (See [values.yaml])
  envFrom: []
  # - configMapRef:
  #     name: config-map-name
  # - secretRef:
  #     name: secret-name

  # -- Annotations to be added to the Redis server pods
  podAnnotations: {}

  # -- Labels to be added to the Redis server pods
  podLabels: {}

  # -- [Node selector]
  nodeSelector: {}
  # -- [Tolerations] for use with node taints
  tolerations: []
  # -- Assign custom [affinity] rules to the deployment
  affinity: {}

  # -- Assign custom [TopologySpreadConstraints] rules to redis
  ## Ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
  ## If labelSelector is left out, it will default to the labelSelector configuration of the deployment
  topologySpreadConstraints: []
  # - maxSkew: 1
  #   topologyKey: topology.kubernetes.io/zone
  #   whenUnsatisfiable: DoNotSchedule

  # -- Priority class for redis
  priorityClassName: ""

  # -- Redis container-level security context
  containerSecurityContext:
    {}
    # capabilities:
    #   drop:
    #     - all
    # readOnlyRootFilesystem: true

  # -- Redis pod-level security context
  securityContext:
    runAsNonRoot: true
    runAsUser: 999

  serviceAccount:
    # -- Create a service account for the redis pod
    create: false
    # -- Service account name for redis pod
    name: ""
    # -- Annotations applied to created service account
    annotations: {}
    # -- Automount API credentials for the Service Account
    automountServiceAccountToken: false

  # -- Resource limits and requests for redis
  resources: {}
  #  limits:
  #    cpu: 200m
  #    memory: 128Mi
  #  requests:
  #    cpu: 100m
  #    memory: 64Mi

  # -- Additional volumeMounts to the redis container
  volumeMounts: []
  # -- Additional volumes to the redis pod
  volumes: []

  # -- Additional containers to be added to the redis pod
  extraContainers: []

  service:
    # -- Redis service annotations
    annotations: {}
    # -- Additional redis service labels
    labels: {}

  metrics:
    # -- Deploy metrics service and redis-exporter sidecar
    enabled: false
    image:
      # -- redis-exporter image repository
      repository: quay.io/bitnami/redis-exporter
      # -- redis-exporter image tag
      tag: 1.26.0-debian-10-r2
      # -- redis-exporter image PullPolicy
      imagePullPolicy: IfNotPresent
    # -- Port to use for redis-exporter sidecar
    containerPort: 9121
    # -- Resource limits and requests for redis-exporter sidecar
    resources: {}
      # limits:
      #   cpu: 50m
      #   memory: 64Mi
      # requests:
      #   cpu: 10m
      #   memory: 32Mi
    service:
      # -- Metrics service type
      type: ClusterIP
      # -- Metrics service clusterIP. None makes a "headless service" (no virtual IP)
      clusterIP: None
      # -- Metrics service annotations
      annotations: {}
      # -- Metrics service labels
      labels: {}
      # -- Metrics service port
      servicePort: 9121
      # -- Metrics service port name
      portName: http-metrics
    serviceMonitor:
      # -- Enable a prometheus ServiceMonitor
      enabled: false
      # -- Interval at which metrics should be scraped
      interval: 30s
      # -- Prometheus [RelabelConfigs] to apply to samples before scraping
      relabelings: []
      # -- Prometheus [MetricRelabelConfigs] to apply to samples before ingestion
      metricRelabelings: []
      # -- Prometheus ServiceMonitor selector
      selector: {}
        # prometheus: kube-prometheus

      # -- Prometheus ServiceMonitor namespace
      namespace: "" # "monitoring"
      # -- Prometheus ServiceMonitor labels
      additionalLabels: {}

# This key configures Redis-HA subchart and when enabled (redis-ha.enabled=true)
# the custom redis deployment is omitted
# Check the redis-ha chart for more properties
redis-ha:
  # -- Enables the Redis HA subchart and disables the custom Redis single node deployment
  enabled: false
  exporter:
    # -- If true, the prometheus exporter sidecar is enabled
    enabled: true
  persistentVolume:
    # -- Configures persistency on Redis nodes
    enabled: false
  redis:
    # -- Redis convention for naming the cluster group: must match ^[\\w-\\.]+$ and can be templated
    masterGroupName: argocd
    # -- Any valid redis config options in this section will be applied to each server (see redis-ha chart)
    # @default -- See [values.yaml]
    config:
      # -- Will save the DB if both the given number of seconds and the given number of write operations against the DB occurred. ""  is disabled
      save: '""'
  haproxy:
    # -- Enabled HAProxy LoadBalancing/Proxy
    enabled: true
    metrics:
      # -- HAProxy enable prometheus metric scraping
      enabled: true
  image:
    # -- Redis tag
    tag: 6.2.4-alpine

## Server
server:
  # -- Argo CD server name
  name: server

  # -- The number of server pods to run
  replicas: 1

  autoscaling:
    # -- Enable Horizontal Pod Autoscaler ([HPA]) for the Argo CD server
    enabled: false
    # -- Minimum number of replicas for the Argo CD server [HPA]
    minReplicas: 1
    # -- Maximum number of replicas for the Argo CD server [HPA]
    maxReplicas: 5
    # -- Average CPU utilization percentage for the Argo CD server [HPA]
    targetCPUUtilizationPercentage: 50
    # -- Average memory utilization percentage for the Argo CD server [HPA]
    targetMemoryUtilizationPercentage: 50

  image:
    # -- Repository to use for the Argo CD server
    # @default -- "" (defaults to global.image.repository)
    repository: "" # defaults to global.image.repository
    # -- Tag to use for the Argo CD server
    # @default -- "" (defaults to global.image.tag)
    tag: "" # defaults to global.image.tag
    # -- Image pull policy for the Argo CD server
    # @default -- "" (defaults to global.image.imagePullPolicy)
    imagePullPolicy: "" # IfNotPresent

  # -- Additional command line arguments to pass to Argo CD server
  extraArgs: []
  #  - --insecure

  # This flag is used to either remove or pass the CLI flag --staticassets /shared/app to the Argo CD server app
  staticAssets:
    # -- Disable deprecated flag --staticassets
    enabled: true

  # -- Environment variables to pass to Argo CD server
  env: []

  # -- envFrom to pass to Argo CD server
  # @default -- [] (See [values.yaml])
  envFrom: []
  # - configMapRef:
  #     name: config-map-name
  # - secretRef:
  #     name: secret-name

  # -- Specify postStart and preStop lifecycle hooks for your argo-cd-server container
  lifecycle: {}

  # -- Argo CD server log format: Either text or json
  logFormat: text
  # -- Argo CD server log level
  logLevel: info

  # -- Annotations to be added to server pods
  podAnnotations: {}

  # -- Labels to be added to server pods
  podLabels: {}

  # -- Configures the server port
  containerPort: 8080

  ## Readiness and liveness probes for default backend
  ## Ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
  ##
  readinessProbe:
    # -- Minimum consecutive failures for the [probe] to be considered failed after having succeeded
    failureThreshold: 3
    # -- Number of seconds after the container has started before [probe] is initiated
    initialDelaySeconds: 10
    # -- How often (in seconds) to perform the [probe]
    periodSeconds: 10
    # -- Minimum consecutive successes for the [probe] to be considered successful after having failed
    successThreshold: 1
    # -- Number of seconds after which the [probe] times out
    timeoutSeconds: 1
  livenessProbe:
    # -- Minimum consecutive failures for the [probe] to be considered failed after having succeeded
    failureThreshold: 3
    # -- Number of seconds after the container has started before [probe] is initiated
    initialDelaySeconds: 10
    # -- How often (in seconds) to perform the [probe]
    periodSeconds: 10
    # -- Minimum consecutive successes for the [probe] to be considered successful after having failed
    successThreshold: 1
    # -- Number of seconds after which the [probe] times out
    timeoutSeconds: 1

  # -- Additional volumeMounts to the server main container
  volumeMounts: []

  # -- Additional volumes to the server pod
  volumes: []

  # -- [Node selector]
  nodeSelector: {}
  # -- [Tolerations] for use with node taints
  tolerations: []
  # -- Assign custom [affinity] rules to the deployment
  affinity: {}

  # -- Assign custom [TopologySpreadConstraints] rules to the Argo CD server
  ## Ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
  ## If labelSelector is left out, it will default to the labelSelector configuration of the deployment
  topologySpreadConstraints: []
  # - maxSkew: 1
  #   topologyKey: topology.kubernetes.io/zone
  #   whenUnsatisfiable: DoNotSchedule

  # -- Priority class for the Argo CD server
  priorityClassName: ""

  # -- Servers container-level security context
  containerSecurityContext:
    {}
    # capabilities:
    #   drop:
    #     - all
    # readOnlyRootFilesystem: true

  # -- Resource limits and requests for the Argo CD server
  resources: {}
  #  limits:
  #    cpu: 100m
  #    memory: 128Mi
  #  requests:
  #    cpu: 50m
  #    memory: 64Mi

  ## Certificate configuration
  certificate:
    # -- Enables a certificate manager certificate
    enabled: false
    # -- Certificate manager domain
    domain: argocd.example.com
    issuer:
      # -- Certificate manager issuer
      kind: # ClusterIssuer
      # -- Certificate manager name
      name: # letsencrypt
    # -- Certificate manager additional hosts
    additionalHosts: []
    # -- Certificate manager secret name
    secretName: argocd-server-tls

  ## Server service configuration
  service:
    # -- Server service annotations
    annotations: {}
    # -- Server service labels
    labels: {}
    # -- Server service type
    type: ClusterIP
    # -- Server service http port for NodePort service type (only if server.service.type is set to "NodePort")
    nodePortHttp: 30080
    # -- Server service https port for NodePort service type (only if server.service.type is set to "NodePort")
    nodePortHttps: 30443
    # -- Server service http port
    servicePortHttp: 80
    # -- Server service https port
    servicePortHttps: 443
    # -- Server service http port name, can be used to route traffic via istio
    servicePortHttpName: http
    # -- Server service https port name, can be used to route traffic via istio
    servicePortHttpsName: https
    # -- Use named target port for argocd
    ## Named target ports are not supported by GCE health checks, so when deploying argocd on GKE
    ## and exposing it via GCE ingress, the health checks fail and the load balancer returns a 502.
    namedTargetPort: true
    # -- LoadBalancer will get created with the IP specified in this field
    loadBalancerIP: ""
    # -- Source IP ranges to allow access to service from
    loadBalancerSourceRanges: []
    # -- Server service external IPs
    externalIPs: []
    # -- Denotes if this Service desires to route external traffic to node-local or cluster-wide endpoints
    externalTrafficPolicy: ""
    # -- Used to maintain session affinity. Supports ClientIP and None
    sessionAffinity: ""

  ## Server metrics service configuration
  metrics:
    # -- Deploy metrics service
    enabled: false
    service:
      # -- Metrics service annotations
      annotations: {}
      # -- Metrics service labels
      labels: {}
      # -- Metrics service port
      servicePort: 8083
    serviceMonitor:
      # -- Enable a prometheus ServiceMonitor
      enabled: false
      # -- Prometheus ServiceMonitor interval
      interval: 30s
      # -- Prometheus [RelabelConfigs] to apply to samples before scraping
      relabelings: []
      # -- Prometheus [MetricRelabelConfigs] to apply to samples before ingestion
      metricRelabelings: []
      # -- Prometheus ServiceMonitor selector
      selector: {}
        # prometheus: kube-prometheus

      # -- Prometheus ServiceMonitor namespace
      namespace: ""  # monitoring
      # -- Prometheus ServiceMonitor labels
      additionalLabels: {}

  serviceAccount:
    # -- Create server service account
    create: true
    # -- Server service account name
    name: argocd-server
    # -- Annotations applied to created service account
    annotations: {}
    # -- Automount API credentials for the Service Account
    automountServiceAccountToken: true

  ingress:
    # -- Enable an ingress resource for the Argo CD server
    enabled: false
    # -- Additional ingress annotations
    annotations: {}
    # -- Additional ingress labels
    labels: {}
    # -- Defines which ingress controller will implement the resource
    ingressClassName: ""

    # -- List of ingress hosts
    ## Argo Ingress.
    ## Hostnames must be provided if Ingress is enabled.
    ## Secrets must be manually created in the namespace
    hosts:
      []
      # - argocd.example.com

    # -- List of ingress paths
    paths:
      - /
    # -- Ingress path type. One of Exact, Prefix or ImplementationSpecific
    pathType: Prefix
    # -- Additional ingress paths
    extraPaths:
      []
      # - path: /*
      #   backend:
      #     serviceName: ssl-redirect
      #     servicePort: use-annotation
      ## for Kubernetes >=1.19 (when "networking.k8s.io/v1" is used)
      # - path: /*
      #   pathType: Prefix
      #   backend:
      #     service:
      #       name: ssl-redirect
      #       port:
      #         name: use-annotation

    # -- Ingress TLS configuration
    tls:
      []
      # - secretName: argocd-tls-certificate
      #   hosts:
      #     - argocd.example.com

    # -- Uses server.service.servicePortHttps instead server.service.servicePortHttp
    https: false

  # dedicated ingress for gRPC as documented at
  # Ref: https://argoproj.github.io/argo-cd/operator-manual/ingress/
  ingressGrpc:
    # -- Enable an ingress resource for the Argo CD server for dedicated [gRPC-ingress]
    enabled: false
    # -- Setup up gRPC ingress to work with an AWS ALB
    isAWSALB: false
    # -- Additional ingress annotations for dedicated [gRPC-ingress]
    annotations: {}
    # -- Additional ingress labels for dedicated [gRPC-ingress]
    labels: {}
    # -- Defines which ingress controller will implement the resource [gRPC-ingress]
    ingressClassName: ""

    awsALB:
      # -- Service type for the AWS ALB gRPC service
      ## Service Type if isAWSALB is set to true
      ## Can be of type NodePort or ClusterIP depending on which mode you are
      ## are running. Instance mode needs type NodePort, IP mode needs type
      ## ClusterIP
      ## Ref: https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.2/how-it-works/#ingress-traffic
      serviceType: NodePort
      # -- Backend protocol version for the AWS ALB gRPC service
      ## This tells AWS to send traffic from the ALB using HTTP2. Can use gRPC as well if you want to leverage gRPC specific features
      backendProtocolVersion: HTTP2

    # -- List of ingress hosts for dedicated [gRPC-ingress]
    ## Argo Ingress.
    ## Hostnames must be provided if Ingress is enabled.
    ## Secrets must be manually created in the namespace
    ##
    hosts:
      []
      # - argocd.example.com

    # -- List of ingress paths for dedicated [gRPC-ingress]
    paths:
      - /
    # -- Ingress path type for dedicated [gRPC-ingress]. One of Exact, Prefix or ImplementationSpecific
    pathType: Prefix
    # -- Additional ingress paths for dedicated [gRPC-ingress]
    extraPaths:
      []
      # - path: /*
      #   backend:
      #     serviceName: ssl-redirect
      #     servicePort: use-annotation
      ## for Kubernetes >=1.19 (when "networking.k8s.io/v1" is used)
      # - path: /*
      #   pathType: Prefix
      #   backend:
      #     service:
      #       name: ssl-redirect
      #       port:
      #         name: use-annotation

    # -- Ingress TLS configuration for dedicated [gRPC-ingress]
    tls:
      []
      # - secretName: argocd-tls-certificate
      #   hosts:
      #     - argocd.example.com

    # -- Uses server.service.servicePortHttps instead server.service.servicePortHttp
    https: false

  # Create a OpenShift Route with SSL passthrough for UI and CLI
  # Consider setting 'hostname' e.g. https://argocd.apps-crc.testing/ using your Default Ingress Controller Domain
  # Find your domain with: kubectl describe --namespace=openshift-ingress-operator ingresscontroller/default | grep Domain:
  # If 'hostname' is an empty string "" OpenShift will create a hostname for you.
  route:
    # -- Enable an OpenShift Route for the Argo CD server
    enabled: false
    # -- Openshift Route annotations
    annotations: {}
    # -- Hostname of OpenShift Route
    hostname: ""
    # -- Termination type of Openshift Route
    termination_type: passthrough
    # -- Termination policy of Openshift Route
    termination_policy: None

  # -- Manage ArgoCD configmap (Declarative Setup)
  ## Ref: https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/argocd-cm.yaml
  configEnabled: true
  # -- [General Argo CD configuration]
  # @default -- See [values.yaml]
  config:
    # Argo CD's externally facing base URL (optional). Required when configuring SSO
    url: https://argocd.example.com
    # Argo CD instance label key
    application.instanceLabelKey: argocd.argoproj.io/instance

    # DEPRECATED: Please instead use configs.credentialTemplates and configs.repositories
    # repositories: |
    #   - url: git@github.com:group/repo.git
    #     sshPrivateKeySecret:
    #       name: secret-name
    #       key: sshPrivateKey
    #   - type: helm
    #     url: https://charts.helm.sh/stable
    #     name: stable
    #   - type: helm
    #     url: https://argoproj.github.io/argo-helm
    #     name: argo

    # oidc.config: |
    #   name: AzureAD
    #   issuer: https://login.microsoftonline.com/TENANT_ID/v2.0
    #   clientID: CLIENT_ID
    #   clientSecret: $oidc.azuread.clientSecret
    #   requestedIDTokenClaims:
    #     groups:
    #       essential: true
    #   requestedScopes:
    #     - openid
    #     - profile
    #     - email

  # -- Annotations to be added to ArgoCD ConfigMap
  configAnnotations: {}

  # -- ArgoCD rbac config ([ArgoCD RBAC policy])
  ## Ref: https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/rbac.md
  rbacConfig:
    {}
    # policy.csv is an file containing user-defined RBAC policies and role definitions (optional).
    # Policy rules are in the form:
    #   p, subject, resource, action, object, effect
    # Role definitions and bindings are in the form:
    #   g, subject, inherited-subject
    # See https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/rbac.md for additional information.
    # policy.csv: |
    #   # Grant all members of the group 'my-org:team-alpha; the ability to sync apps in 'my-project'
    #   p, my-org:team-alpha, applications, sync, my-project/*, allow
    #   # Grant all members of 'my-org:team-beta' admins
    #   g, my-org:team-beta, role:admin
    # policy.default is the name of the default role which Argo CD will falls back to, when
    # authorizing API requests (optional). If omitted or empty, users may be still be able to login,
    # but will see no apps, projects, etc...
    # policy.default: role:readonly
    # scopes controls which OIDC scopes to examine during rbac enforcement (in addition to sub scope).
    # If omitted, defaults to: '[groups]'. The scope value can be a string, or a list of strings.
    # scopes: '[cognito:groups, email]'

  # -- Annotations to be added to ArgoCD rbac ConfigMap
  rbacConfigAnnotations: {}

  # -- Whether or not to create the configmap. If false, it is expected the configmap will be created
  # by something else. ArgoCD will not work if there is no configMap created with the name above.
  rbacConfigCreate: true

  # -- Deploy ArgoCD Applications within this helm release
  # @default -- [] (See [values.yaml])
  ## Ref: https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/
  additionalApplications: []
  # - name: guestbook
  #   namespace: argocd
  #   additionalLabels: {}
  #   additionalAnnotations: {}
  #   finalizers:
  #   - resources-finalizer.argocd.argoproj.io
  #   project: guestbook
  #   source:
  #     repoURL: https://github.com/argoproj/argocd-example-apps.git
  #     targetRevision: HEAD
  #     path: guestbook
  #     directory:
  #       recurse: true
  #   destination:
  #     server: https://kubernetes.default.svc
  #     namespace: guestbook
  #   syncPolicy:
  #     automated:
  #       prune: false
  #       selfHeal: false
  #   ignoreDifferences:
  #   - group: apps
  #     kind: Deployment
  #     jsonPointers:
  #     - /spec/replicas
  #   info:
  #   - name: url
  #     value: https://argoproj.github.io/

  # -- Deploy ArgoCD Projects within this helm release
  # @default -- [] (See [values.yaml])
  ## Ref: https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/
  additionalProjects: []
  # - name: guestbook
  #   namespace: argocd
  #   additionalLabels: {}
  #   additionalAnnotations: {}
  #   finalizers:
  #   - resources-finalizer.argocd.argoproj.io
  #   description: Example Project
  #   sourceRepos:
  #   - '*'
  #   destinations:
  #   - namespace: guestbook
  #     server: https://kubernetes.default.svc
  #   clusterResourceWhitelist: []
  #   namespaceResourceBlacklist:
  #   - group: ''
  #     kind: ResourceQuota
  #   - group: ''
  #     kind: LimitRange
  #   - group: ''
  #     kind: NetworkPolicy
  #     orphanedResources: {}
  #     roles: []
  #   namespaceResourceWhitelist:
  #   - group: 'apps'
  #     kind: Deployment
  #   - group: 'apps'
  #     kind: StatefulSet
  #   orphanedResources: {}
  #   roles: []
  #   syncWindows:
  #   - kind: allow
  #     schedule: '10 1 * * *'
  #     duration: 1h
  #     applications:
  #     - '*-prod'
  #     manualSync: true
  #   signatureKeys:
  #   - keyID: ABCDEF1234567890

  ## Enable Admin ClusterRole resources.
  ## Enable if you would like to grant rights to ArgoCD to deploy to the local Kubernetes cluster.
  clusterAdminAccess:
    # -- Enable RBAC for local cluster deployments
    enabled: true

  GKEbackendConfig:
    # -- Enable BackendConfig custom resource for Google Kubernetes Engine
    enabled: false
    # -- [BackendConfigSpec]
    spec: {}
  #  spec:
  #    iap:
  #      enabled: true
  #      oauthclientCredentials:
  #        secretName: argocd-secret

  # -- Additional containers to be added to the server pod
  ## See https://github.com/lemonldap-ng-controller/lemonldap-ng-controller as example.
  extraContainers: []
  # - name: my-sidecar
  #   image: nginx:latest
  # - name: lemonldap-ng-controller
  #   image: lemonldapng/lemonldap-ng-controller:0.2.0
  #   args:
  #     - /lemonldap-ng-controller
  #     - --alsologtostderr
  #     - --configmap=$(POD_NAMESPACE)/lemonldap-ng-configuration
  #   env:
  #     - name: POD_NAME
  #       valueFrom:
  #         fieldRef:
  #           fieldPath: metadata.name
  #     - name: POD_NAMESPACE
  #       valueFrom:
  #         fieldRef:
  #           fieldPath: metadata.namespace
  #   volumeMounts:
  #   - name: copy-portal-skins
  #     mountPath: /srv/var/lib/lemonldap-ng/portal/skins

## Repo Server
repoServer:
  # -- Repo server name
  name: repo-server

  # -- The number of repo server pods to run
  replicas: 1

  autoscaling:
    # -- Enable Horizontal Pod Autoscaler ([HPA]) for the repo server
    enabled: false
    # -- Minimum number of replicas for the repo server [HPA]
    minReplicas: 1
    # -- Maximum number of replicas for the repo server [HPA]
    maxReplicas: 5
    # -- Average CPU utilization percentage for the repo server [HPA]
    targetCPUUtilizationPercentage: 50
    # -- Average memory utilization percentage for the repo server [HPA]
    targetMemoryUtilizationPercentage: 50

  image:
    # -- Repository to use for the repo server
    # @default -- "" (defaults to global.image.repository)
    repository: "" # defaults to global.image.repository
    # -- Tag to use for the repo server
    # @default -- "" (defaults to global.image.tag)
    tag: "" # defaults to global.image.tag
    # -- Image pull policy for the repo server
    # @default -- "" (defaults to global.image.imagePullPolicy)
    imagePullPolicy: "" # IfNotPresent

  # -- Additional command line arguments to pass to repo server
  extraArgs: []

  # -- Environment variables to pass to repo server
  env: []

  # -- envFrom to pass to repo server
  # @default -- [] (See [values.yaml])
  envFrom: []
  # - configMapRef:
  #     name: config-map-name
  # - secretRef:
  #     name: secret-name

  # -- Repo server log format: Either text or json
  logFormat: text
  # -- Repo server log level
  logLevel: info

  # -- Annotations to be added to repo server pods
  podAnnotations: {}

  # -- Labels to be added to repo server pods
  podLabels: {}

  # -- Configures the repo server port
  containerPort: 8081

  ## Readiness and liveness probes for default backend
  ## Ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
  ##
  readinessProbe:
    # -- Minimum consecutive failures for the [probe] to be considered failed after having succeeded
    failureThreshold: 3
    # -- Number of seconds after the container has started before [probe] is initiated
    initialDelaySeconds: 10
    # -- How often (in seconds) to perform the [probe]
    periodSeconds: 10
    # -- Minimum consecutive successes for the [probe] to be considered successful after having failed
    successThreshold: 1
    # -- Number of seconds after which the [probe] times out
    timeoutSeconds: 1
  livenessProbe:
    # -- Minimum consecutive failures for the [probe] to be considered failed after having succeeded
    failureThreshold: 3
    # -- Number of seconds after the container has started before [probe] is initiated
    initialDelaySeconds: 10
    # -- How often (in seconds) to perform the [probe]
    periodSeconds: 10
    # -- Minimum consecutive successes for the [probe] to be considered successful after having failed
    successThreshold: 1
    # -- Number of seconds after which the [probe] times out
    timeoutSeconds: 1

  # -- Additional volumeMounts to the repo server main container
  volumeMounts: []

  # -- Additional volumes to the repo server pod
  volumes: []
  ## Use init containers to configure custom tooling
  ## https://argoproj.github.io/argo-cd/operator-manual/custom_tools/
  ## When using the volumes & volumeMounts section bellow, please comment out those above.
  #  - name: custom-tools
  #    emptyDir: {}

  # -- [Node selector]
  nodeSelector: {}
  # -- [Tolerations] for use with node taints
  tolerations: []
  # -- Assign custom [affinity] rules to the deployment
  affinity: {}

  # -- Assign custom [TopologySpreadConstraints] rules to the repo server
  ## Ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
  ## If labelSelector is left out, it will default to the labelSelector configuration of the deployment
  topologySpreadConstraints: []
  # - maxSkew: 1
  #   topologyKey: topology.kubernetes.io/zone
  #   whenUnsatisfiable: DoNotSchedule

  # -- Priority class for the repo server
  priorityClassName: ""

  # -- Repo server container-level security context
  containerSecurityContext:
    {}
    # capabilities:
    #   drop:
    #     - all
    # readOnlyRootFilesystem: true

  # -- Resource limits and requests for the repo server pods
  resources: {}
  #  limits:
  #    cpu: 50m
  #    memory: 128Mi
  #  requests:
  #    cpu: 10m
  #    memory: 64Mi

  ## Repo server service configuration
  service:
    # -- Repo server service annotations
    annotations: {}
    # -- Repo server service labels
    labels: {}
    # -- Repo server service port
    port: 8081
    # -- Repo server service port name
    portName: https-repo-server

  ## Repo server metrics service configuration
  metrics:
    # -- Deploy metrics service
    enabled: false
    service:
      # -- Metrics service annotations
      annotations: {}
      # -- Metrics service labels
      labels: {}
      # -- Metrics service port
      servicePort: 8084
    serviceMonitor:
      # -- Enable a prometheus ServiceMonitor
      enabled: false
      # -- Prometheus ServiceMonitor interval
      interval: 30s
      # -- Prometheus [RelabelConfigs] to apply to samples before scraping
      relabelings: []
      # -- Prometheus [MetricRelabelConfigs] to apply to samples before ingestion
      metricRelabelings: []
      # -- Prometheus ServiceMonitor selector
      selector: {}
        # prometheus: kube-prometheus

      # -- Prometheus ServiceMonitor namespace
      namespace: "" # "monitoring"
      # -- Prometheus ServiceMonitor labels
      additionalLabels: {}

  ## Enable Admin ClusterRole resources.
  ## Enable if you would like to grant cluster rights to ArgoCD repo server.
  clusterAdminAccess:
    # -- Enable RBAC for local cluster deployments
    enabled: false
  ## Enable Custom Rules for the Repo server's Cluster Role resource
  ## Enable this and set the rules: to whatever custom rules you want for the Cluster Role resource.
  ## Defaults to off
  clusterRoleRules:
    # -- Enable custom rules for the Repo server's Cluster Role resource
    enabled: false
    # -- List of custom rules for the Repo server's Cluster Role resource
    rules: []

  ## Repo server service account
  ## If create is set to true, make sure to uncomment the name and update the rbac section below
  serviceAccount:
    # -- Create repo server service account
    create: false
    # -- Repo server service account name
    name: "" # "argocd-repo-server"
    # -- Annotations applied to created service account
    annotations: {}
    # -- Automount API credentials for the Service Account
    automountServiceAccountToken: true

  # -- Additional containers to be added to the repo server pod
  extraContainers: []

  # -- Repo server rbac rules
  rbac: []
  #   - apiGroups:
  #     - argoproj.io
  #     resources:
  #     - applications
  #     verbs:
  #     - get
  #     - list
  #     - watch

  # -- Init containers to add to the repo server pods
  initContainers: []
  #  - name: download-tools
  #    image: alpine:3.8
  #    command: [sh, -c]
  #    args:
  #      - wget -qO- https://get.helm.sh/helm-v2.16.1-linux-amd64.tar.gz | tar -xvzf - &&
  #        mv linux-amd64/helm /custom-tools/
  #    volumeMounts:
  #      - mountPath: /custom-tools
  #        name: custom-tools
  #  volumeMounts:
  #  - mountPath: /usr/local/bin/helm
  #    name: custom-tools
  #    subPath: helm

## Argo Configs
configs:
  # -- Provide one or multiple [external cluster credentials]
  # @default -- [] (See [values.yaml])
  ## Ref:
  ## - https://argoproj.github.io/argo-cd/operator-manual/declarative-setup/#clusters
  ## - https://argoproj.github.io/argo-cd/operator-manual/security/#external-cluster-credentials
  clusterCredentials: []
    # - name: mycluster
    #   server: https://mycluster.com
    #   labels: {}
    #   annotations: {}
    #   config:
    #     bearerToken: "<authentication token>"
    #     tlsClientConfig:
    #       insecure: false
    #       caData: "<base64 encoded certificate>"
    # - name: mycluster2
    #   server: https://mycluster2.com
    #   labels: {}
    #   annotations: {}
    #   namespaces: namespace1,namespace2
    #   config:
    #     bearerToken: "<authentication token>"
    #     tlsClientConfig:
    #       insecure: false
    #       caData: "<base64 encoded certificate>"

  # -- GnuPG key ring annotations
  gpgKeysAnnotations: {}
  # -- [GnuPG](https://argoproj.github.io/argo-cd/user-guide/gpg-verification/) keys to add to the key ring
  # @default -- {} (See [values.yaml])
  gpgKeys: {}
    # 4AEE18F83AFDEB23: |
    #     -----BEGIN PGP PUBLIC KEY BLOCK-----
    #
    #     mQENBFmUaEEBCACzXTDt6ZnyaVtueZASBzgnAmK13q9Urgch+sKYeIhdymjuMQta
    #     x15OklctmrZtqre5kwPUosG3/B2/ikuPYElcHgGPL4uL5Em6S5C/oozfkYzhwRrT
    #     SQzvYjsE4I34To4UdE9KA97wrQjGoz2Bx72WDLyWwctD3DKQtYeHXswXXtXwKfjQ
    #     7Fy4+Bf5IPh76dA8NJ6UtjjLIDlKqdxLW4atHe6xWFaJ+XdLUtsAroZcXBeWDCPa
    #     buXCDscJcLJRKZVc62gOZXXtPfoHqvUPp3nuLA4YjH9bphbrMWMf810Wxz9JTd3v
    #     yWgGqNY0zbBqeZoGv+TuExlRHT8ASGFS9SVDABEBAAG0NUdpdEh1YiAod2ViLWZs
    #     b3cgY29tbWl0IHNpZ25pbmcpIDxub3JlcGx5QGdpdGh1Yi5jb20+iQEiBBMBCAAW
    #     BQJZlGhBCRBK7hj4Ov3rIwIbAwIZAQAAmQEH/iATWFmi2oxlBh3wAsySNCNV4IPf
    #     DDMeh6j80WT7cgoX7V7xqJOxrfrqPEthQ3hgHIm7b5MPQlUr2q+UPL22t/I+ESF6
    #     9b0QWLFSMJbMSk+BXkvSjH9q8jAO0986/pShPV5DU2sMxnx4LfLfHNhTzjXKokws
    #     +8ptJ8uhMNIDXfXuzkZHIxoXk3rNcjDN5c5X+sK8UBRH092BIJWCOfaQt7v7wig5
    #     4Ra28pM9GbHKXVNxmdLpCFyzvyMuCmINYYADsC848QQFFwnd4EQnupo6QvhEVx1O
    #     j7wDwvuH5dCrLuLwtwXaQh0onG4583p0LGms2Mf5F+Ick6o/4peOlBoZz48=
    #     =Bvzs
    #     -----END PGP PUBLIC KEY BLOCK-----

  # -- Known Hosts configmap annotations
  knownHostsAnnotations: {}
  knownHosts:
    data:
      # -- Known Hosts
      # @default -- See [values.yaml]
      ssh_known_hosts: |
        bitbucket.org ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
        github.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=
        github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl
        github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
        gitlab.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=
        gitlab.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAfuCHKVTjquxvt6CM6tdG4SLp1Btn/nOeHHE5UOzRdf
        gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9
        ssh.dev.azure.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
        vs-ssh.visualstudio.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
  # -- TLS certificate configmap annotations
  tlsCertsAnnotations: {}
  # -- TLS certificate
  # @default -- See [values.yaml]
  tlsCerts:
    {}
    # data:
    #   argocd.example.com: |
    #     -----BEGIN CERTIFICATE-----
    #     MIIF1zCCA7+gAwIBAgIUQdTcSHY2Sxd3Tq/v1eIEZPCNbOowDQYJKoZIhvcNAQEL
    #     BQAwezELMAkGA1UEBhMCREUxFTATBgNVBAgMDExvd2VyIFNheG9ueTEQMA4GA1UE
    #     BwwHSGFub3ZlcjEVMBMGA1UECgwMVGVzdGluZyBDb3JwMRIwEAYDVQQLDAlUZXN0
    #     c3VpdGUxGDAWBgNVBAMMD2Jhci5leGFtcGxlLmNvbTAeFw0xOTA3MDgxMzU2MTda
    #     Fw0yMDA3MDcxMzU2MTdaMHsxCzAJBgNVBAYTAkRFMRUwEwYDVQQIDAxMb3dlciBT
    #     YXhvbnkxEDAOBgNVBAcMB0hhbm92ZXIxFTATBgNVBAoMDFRlc3RpbmcgQ29ycDES
    #     MBAGA1UECwwJVGVzdHN1aXRlMRgwFgYDVQQDDA9iYXIuZXhhbXBsZS5jb20wggIi
    #     MA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCv4mHMdVUcafmaSHVpUM0zZWp5
    #     NFXfboxA4inuOkE8kZlbGSe7wiG9WqLirdr39Ts+WSAFA6oANvbzlu3JrEQ2CHPc
    #     CNQm6diPREFwcDPFCe/eMawbwkQAPVSHPts0UoRxnpZox5pn69ghncBR+jtvx+/u
    #     P6HdwW0qqTvfJnfAF1hBJ4oIk2AXiip5kkIznsAh9W6WRy6nTVCeetmIepDOGe0G
    #     ZJIRn/OfSz7NzKylfDCat2z3EAutyeT/5oXZoWOmGg/8T7pn/pR588GoYYKRQnp+
    #     YilqCPFX+az09EqqK/iHXnkdZ/Z2fCuU+9M/Zhrnlwlygl3RuVBI6xhm/ZsXtL2E
    #     Gxa61lNy6pyx5+hSxHEFEJshXLtioRd702VdLKxEOuYSXKeJDs1x9o6cJ75S6hko
    #     Ml1L4zCU+xEsMcvb1iQ2n7PZdacqhkFRUVVVmJ56th8aYyX7KNX6M9CD+kMpNm6J
    #     kKC1li/Iy+RI138bAvaFplajMF551kt44dSvIoJIbTr1LigudzWPqk31QaZXV/4u
    #     kD1n4p/XMc9HYU/was/CmQBFqmIZedTLTtK7clkuFN6wbwzdo1wmUNgnySQuMacO
    #     gxhHxxzRWxd24uLyk9Px+9U3BfVPaRLiOPaPoC58lyVOykjSgfpgbus7JS69fCq7
    #     bEH4Jatp/10zkco+UQIDAQABo1MwUTAdBgNVHQ4EFgQUjXH6PHi92y4C4hQpey86
    #     r6+x1ewwHwYDVR0jBBgwFoAUjXH6PHi92y4C4hQpey86r6+x1ewwDwYDVR0TAQH/
    #     BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEAFE4SdKsX9UsLy+Z0xuHSxhTd0jfn
    #     Iih5mtzb8CDNO5oTw4z0aMeAvpsUvjJ/XjgxnkiRACXh7K9hsG2r+ageRWGevyvx
    #     CaRXFbherV1kTnZw4Y9/pgZTYVWs9jlqFOppz5sStkfjsDQ5lmPJGDii/StENAz2
    #     XmtiPOgfG9Upb0GAJBCuKnrU9bIcT4L20gd2F4Y14ccyjlf8UiUi192IX6yM9OjT
    #     +TuXwZgqnTOq6piVgr+FTSa24qSvaXb5z/mJDLlk23npecTouLg83TNSn3R6fYQr
    #     d/Y9eXuUJ8U7/qTh2Ulz071AO9KzPOmleYPTx4Xty4xAtWi1QE5NHW9/Ajlv5OtO
    #     OnMNWIs7ssDJBsB7VFC8hcwf79jz7kC0xmQqDfw51Xhhk04kla+v+HZcFW2AO9so
    #     6ZdVHHQnIbJa7yQJKZ+hK49IOoBR6JgdB5kymoplLLiuqZSYTcwSBZ72FYTm3iAr
    #     jzvt1hxpxVDmXvRnkhRrIRhK4QgJL0jRmirBjDY+PYYd7bdRIjN7WNZLFsgplnS8
    #     9w6CwG32pRlm0c8kkiQ7FXA6BYCqOsDI8f1VGQv331OpR2Ck+FTv+L7DAmg6l37W
    #     +LB9LGh4OAp68ImTjqf6ioGKG0RBSznwME+r4nXtT1S/qLR6ASWUS4ViWRhbRlNK
    #     XWyb96wrUlv+E8I=
    #     -----END CERTIFICATE-----

  # -- *DEPRECATED:* Instead, use configs.credentialTemplates and/or configs.repositories
  repositoryCredentials: {}

  # -- Repository credentials to be used as Templates for other repos
  ## Creates a secret for each key/value specified below to create repository credentials
  credentialTemplates: {}
    # github-enterprise-creds-1:
    #   url: https://github.com/argoproj
    #   githubAppID: 1
    #   githubAppInstallationID: 2
    #   githubAppEnterpriseBaseUrl: https://ghe.example.com/api/v3
    #   githubAppPrivateKey: |
    #     -----BEGIN OPENSSH PRIVATE KEY-----
    #     ...
    #     -----END OPENSSH PRIVATE KEY-----
    # https-creds:
    #   url: https://github.com/argoproj
    #   password: my-password
    #   username: my-username
    # ssh-creds:
    #  url: git@github.com:argoproj-labs
    #  sshPrivateKey: |
    #    -----BEGIN OPENSSH PRIVATE KEY-----
    #    ...
    #    -----END OPENSSH PRIVATE KEY-----

  # -- Repositories list to be used by applications
  ## Creates a secret for each key/value specified below to create repositories
  ## Note: the last example in the list would use a repository credential template, configured under "configs.repositoryCredentials".
  repositories: {}
    # istio-helm-repo:
    #   url: https://storage.googleapis.com/istio-prerelease/daily-build/master-latest-daily/charts
    #   name: istio.io
    #   type: helm
    # private-helm-repo:
    #   url: https://my-private-chart-repo.internal
    #   name: private-repo
    #   type: helm
    #   password: my-password
    #   username: my-username
    # private-repo:
    #   url: https://github.com/argoproj/private-repo

  secret:
    # -- Create the argocd-secret
    createSecret: true
    # -- Annotations to be added to argocd-secret
    annotations: {}

    # -- Shared secret for authenticating GitHub webhook events
    githubSecret: ""
    # -- Shared secret for authenticating GitLab webhook events
    gitlabSecret: ""
    # -- Shared secret for authenticating BitbucketServer webhook events
    bitbucketServerSecret: ""
    # -- UUID for authenticating Bitbucket webhook events
    bitbucketUUID: ""
    # -- Shared secret for authenticating Gogs webhook events
    gogsSecret: ""

    # -- add additional secrets to be added to argocd-secret
    ## Custom secrets. Useful for injecting SSO secrets into environment variables.
    ## Ref: https://argoproj.github.io/argo-cd/operator-manual/sso/
    ## Note that all values must be non-empty.
    extra:
      {}
      # LDAP_PASSWORD: "mypassword"

    # -- Argo TLS Data
    argocdServerTlsConfig:
      {}
      # key:
      # crt: |
      #   -----BEGIN CERTIFICATE-----
      #   <cert data>
      #   -----END CERTIFICATE-----
      #   -----BEGIN CERTIFICATE-----
      #   <ca cert data>
      #   -----END CERTIFICATE-----

    # -- Bcrypt hashed admin password
    ## Argo expects the password in the secret to be bcrypt hashed. You can create this hash with
    ## htpasswd -nbBC 10 "" $ARGO_PWD | tr -d ':\n' | sed 's/$2y/$2a/'
    argocdServerAdminPassword: ""
    # -- Admin password modification time. Eg. "2006-01-02T15:04:05Z"
    # @default -- "" (defaults to current time)
    argocdServerAdminPasswordMtime: ""

  # -- Define custom [CSS styles] for your argo instance.
  # This setting will automatically mount the provided CSS and reference it in the argo configuration.
  # @default -- "" (See [values.yaml])
  ## Ref: https://argo-cd.readthedocs.io/en/stable/operator-manual/custom-styles/
  styles: ""
  # styles: |
  #  .nav-bar {
  #    background: linear-gradient(to bottom, #999, #777, #333, #222, #111);
  #  }

openshift:
  # -- enables using arbitrary uid for argo repo server
  enabled: false
`
