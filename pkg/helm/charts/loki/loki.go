// Package loki provides a Helm chart for installing:
// - https://github.com/grafana/helm-charts/tree/main/charts/loki
// - https://grafana.com/oss/loki/
package loki

import (
	"bytes"
	"text/template"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/helm"
)

// New returns an initialised Helm chart for installing cluster-Loki
func New(values *Values) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "grafana",
		RepositoryURL:  "https://grafana.github.io/helm-charts",
		ReleaseName:    "loki",
		Version:        "2.3.0",
		Chart:          "loki",
		Namespace:      "monitoring",
		Timeout:        constant.DefaultChartApplyTimeout,
		Values:         values,
	}
}

// NewDefaultValues returns the mapped values.yml containing
// the default values
func NewDefaultValues() *Values {
	return &Values{}
}

// Values contains the required inputs for generating the values.yml
// One of those cases where there really isn't much to change, but
// I will leave these structures here nonetheless.
type Values struct{}

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
const valuesTemplate = `image:
  repository: grafana/loki
  tag: 2.1.0
  pullPolicy: IfNotPresent

  ## Optionally specify an array of imagePullSecrets.
  ## Secrets must be manually created in the namespace.
  ## ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
  ##
  # pullSecrets:
  #   - myRegistryKeySecretName

ingress:
  enabled: false
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: chart-example.local
      paths: []
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

## Affinity for pod assignment
## ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
affinity: {}
# podAntiAffinity:
#   requiredDuringSchedulingIgnoredDuringExecution:
#   - labelSelector:
#       matchExpressions:
#       - key: app
#         operator: In
#         values:
#         - loki
#     topologyKey: "kubernetes.io/hostname"

## StatefulSet annotations
annotations: {}

# enable tracing for debug, need install jaeger and specify right jaeger_agent_host
tracing:
  jaegerAgentHost:

config:
  auth_enabled: false
  ingester:
    chunk_idle_period: 3m
    chunk_block_size: 262144
    chunk_retain_period: 1m
    max_transfer_retries: 0
    lifecycler:
      ring:
        kvstore:
          store: inmemory
        replication_factor: 1

      ## Different ring configs can be used. E.g. Consul
      # ring:
      #   store: consul
      #   replication_factor: 1
      #   consul:
      #     host: "consul:8500"
      #     prefix: ""
      #     http_client_timeout: "20s"
      #     consistent_reads: true
  limits_config:
    enforce_metric_name: false
    reject_old_samples: true
    reject_old_samples_max_age: 168h
  schema_config:
    configs:
    - from: "2020-10-24"
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h
  server:
    http_listen_port: 3100
  storage_config:
    boltdb_shipper:
      active_index_directory: /data/loki/boltdb-shipper-active
      cache_location: /data/loki/boltdb-shipper-cache
      cache_ttl: 24h         # Can be increased for faster performance over longer query periods, uses more disk space
      shared_store: filesystem
    filesystem:
      directory: /data/loki/chunks
  chunk_store_config:
    max_look_back_period: 0s
  table_manager:
    retention_deletes_enabled: false
    retention_period: 0s
  compactor:
    working_directory: /data/loki/boltdb-shipper-compactor
    shared_store: filesystem
# Needed for Alerting: https://grafana.com/docs/loki/latest/alerting/
# This is just a simple example, for more details: https://grafana.com/docs/loki/latest/configuration/#ruler_config
#  ruler:
#    storage:
#      type: local
#      local:
#        directory: /rules
#    rule_path: /tmp/scratch
#    alertmanager_url: http://alertmanager.svc.namespace:9093
#    ring:
#      kvstore:
#        store: inmemory
#    enable_api: true

## Additional Loki container arguments, e.g. log level (debug, info, warn, error)
extraArgs: {}
  # log.level: debug

livenessProbe:
  httpGet:
    path: /ready
    port: http-metrics
  initialDelaySeconds: 45

## ref: https://kubernetes.io/docs/concepts/services-networking/network-policies/
networkPolicy:
  enabled: false

## The app name of loki clients
client: {}
  # name:

## ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
nodeSelector: {}

## ref: https://kubernetes.io/docs/concepts/storage/persistent-volumes/
## If you set enabled as "True", you need :
## - create a pv which above 10Gi and has same namespace with loki
## - keep storageClassName same with below setting
persistence:
  enabled: false
  accessModes:
  - ReadWriteOnce
  size: 10Gi
  annotations: {}
  # selector:
  #   matchLabels:
  #     app.kubernetes.io/name: loki
  # subPath: ""
  # existingClaim:

## Pod Labels
podLabels: {}

## Pod Annotations
podAnnotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "http-metrics"

podManagementPolicy: OrderedReady

## Assign a PriorityClassName to pods if set
# priorityClassName:

rbac:
  create: true
  pspEnabled: true

readinessProbe:
  httpGet:
    path: /ready
    port: http-metrics
  initialDelaySeconds: 45

replicas: 1

resources: {}
# limits:
#   cpu: 200m
#   memory: 256Mi
# requests:
#   cpu: 100m
#   memory: 128Mi

securityContext:
  fsGroup: 10001
  runAsGroup: 10001
  runAsNonRoot: true
  runAsUser: 10001

service:
  type: ClusterIP
  nodePort:
  port: 3100
  annotations: {}
  labels: {}
  targetPort: http-metrics

serviceAccount:
  create: true
  name:
  annotations: {}

terminationGracePeriodSeconds: 4800

## Tolerations for pod assignment
## ref: https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
tolerations: []

# The values to set in the PodDisruptionBudget spec
# If not set then a PodDisruptionBudget will not be created
podDisruptionBudget: {}
# minAvailable: 1
# maxUnavailable: 1

updateStrategy:
  type: RollingUpdate

serviceMonitor:
  enabled: false
  interval: ""
  additionalLabels: {}
  annotations: {}
  # scrapeTimeout: 10s

initContainers: []
## Init containers to be added to the loki pod.
# - name: my-init-container
#   image: busybox:latest
#   command: ['sh', '-c', 'echo hello']

extraContainers: []
## Additional containers to be added to the loki pod.
# - name: reverse-proxy
#   image: angelbarrera92/basic-auth-reverse-proxy:dev
#   args:
#     - "serve"
#     - "--upstream=http://localhost:3100"
#     - "--auth-config=/etc/reverse-proxy-conf/authn.yaml"
#   ports:
#     - name: http
#       containerPort: 11811
#       protocol: TCP
#   volumeMounts:
#     - name: reverse-proxy-auth-config
#       mountPath: /etc/reverse-proxy-conf


extraVolumes: []
## Additional volumes to the loki pod.
# - name: reverse-proxy-auth-config
#   secret:
#     secretName: reverse-proxy-auth-config

## Extra volume mounts that will be added to the loki container
extraVolumeMounts: []

extraPorts: []
## Additional ports to the loki services. Useful to expose extra container ports.
# - port: 11811
#   protocol: TCP
#   name: http
#   targetPort: http

# Extra env variables to pass to the loki container
env: []

# Specify Loki Alerting rules based on this documentation: https://grafana.com/docs/loki/latest/alerting/
# When specified, you also need to add a ruler config section above. An example is shown in the alerting docs.
alerting_groups: []
#  - name: example
#    rules:
#    - alert: HighThroughputLogStreams
#      expr: sum by(container) (rate({job=~"loki-dev/.*"}[1m])) > 1000
#      for: 2m
`
