// Package promtail provides a Helm chart for installing:
// - https://github.com/grafana/helm-charts/tree/main/charts/promtail
package promtail

import (
	"bytes"
	"text/template"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/helm"
)

// New returns an initialised Helm chart for installing cluster-promtail
func New(values *Values) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "grafana",
		RepositoryURL:  "https://grafana.github.io/helm-charts",
		ReleaseName:    "promtail",
		Version:        "3.1.0",
		Chart:          "promtail",
		Namespace:      "monitoring",
		Timeout:        config.DefaultChartApplyTimeout,
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
	// Need to change the default delimiters as these collide with
	// Helm's own delimiters.
	tmpl, err := template.New("values").Delims("{{{", "}}}").Parse(valuesTemplate)
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
const valuesTemplate = `# -- Overrides the chart's name
nameOverride: null

# -- Overrides the chart's computed fullname
fullnameOverride: null

initContainer:
  # -- Specifies whether the init container for setting inotify max user instances is to be enabled
  enabled: false
  image:
    # -- The Docker registry for the init container
    registry: docker.io
    # -- Docker image repository for the init container
    repository: busybox
    # -- Docker tag for the init container
    tag: 1.33
    # -- Docker image pull policy for the init container image
    pullPolicy: IfNotPresent
  # -- The inotify max user instances to configure
  fsInotifyMaxUserInstances: 128

image:
  # -- The Docker registry
  registry: docker.io
  # -- Docker image repository
  repository: grafana/promtail
  # -- Overrides the image tag whose default is the chart's appVersion
  tag: null
  # -- Docker image pull policy
  pullPolicy: IfNotPresent

# -- Image pull secrets for Docker images
imagePullSecrets: []

# -- Annotations for the SaemonSet
annotations: {}

# -- The update strategy for the DaemonSet
updateStrategy: {}

# -- Pod labels
podLabels: {}

# -- Pod annotations
podAnnotations: {}
#  prometheus.io/scrape: "true"
#  prometheus.io/port: "http-metrics"

# -- The name of the PriorityClass
priorityClassName: null

# -- Liveness probe
livenessProbe: {}

# -- Readiness probe
# @default -- See "values.yaml"
readinessProbe:
  failureThreshold: 5
  httpGet:
    path: /ready
    port: http-metrics
  initialDelaySeconds: 10
  periodSeconds: 10
  successThreshold: 1
  timeoutSeconds: 1

# -- Resource requests and limits
resources: {}
#  limits:
#    cpu: 200m
#    memory: 128Mi
#  requests:
#    cpu: 100m
#    memory: 128Mi

# -- The security context for pods
podSecurityContext:
  runAsUser: 0
  runAsGroup: 0

# -- The security context for containers
containerSecurityContext:
  readOnlyRootFilesystem: true
  capabilities:
    drop:
      - ALL
  allowPrivilegeEscalation: false

rbac:
  # -- Specifies whether RBAC resources are to be created
  create: true
  # -- Specifies whether a PodSecurityPolicy is to be created
  pspEnabled: false

serviceAccount:
  # -- Specifies whether a ServiceAccount should be created
  create: true
  # -- The name of the ServiceAccount to use.
  # If not set and "create" is true, a name is generated using the fullname template
  name: null
  # -- Image pull secrets for the service account
  imagePullSecrets: []
  # -- Annotations for the service account
  annotations: {}

# -- Node selector for pods
nodeSelector: {}

# -- Affinity configuration for pods
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
            - key: "eks.amazonaws.com/compute-type"
              operator: NotIn
              values:
                - fargate

# -- Tolerations for pods. By default, pods will be scheduled on master nodes.
tolerations:
  - key: node-role.kubernetes.io/master
    operator: Exists
    effect: NoSchedule

# -- Default volumes that are mounted into pods. In most cases, these should not be changed.
# Use "extraVolumes"/"extraVolumeMounts" for additional custom volumes.
# @default -- See "values.yaml"
defaultVolumes:
  - name: containers
    hostPath:
      path: /var/lib/docker/containers
  - name: pods
    hostPath:
      path: /var/log/pods

# -- Default volume mounts. Corresponds to "volumes".
# @default -- See "values.yaml"
defaultVolumeMounts:
  - name: containers
    mountPath: /var/lib/docker/containers
    readOnly: true
  - name: pods
    mountPath: /var/log/pods
    readOnly: true

# Extra volumes to be added in addition to those specified under "defaultVolumes".
extraVolumes: []

# Extra volume mounts together. Corresponds to "extraVolumes".
extraVolumeMounts: []

# Extra args for the Promtail container.
extraArgs: []
# -- Example:
# -- extraArgs:
# --   - -client.external-labels=hostname=$(HOSTNAME)

# -- Extra environment variables
extraEnv: []

# -- Extra environment variables from secrets or configmaps
extraEnvFrom: []

# ServiceMonitor configuration
serviceMonitor:
  # -- If enabled, ServiceMonitor resources for Prometheus Operator are created
  enabled: false
  # -- Alternative namespace for ServiceMonitor resources
  namespace: null
  # -- Namespace selector for ServiceMonitor resources
  namespaceSelector: {}
  # -- ServiceMonitor annotations
  annotations: {}
  # -- Additional ServiceMonitor labels
  labels: {}
  # -- ServiceMonitor scrape interval
  interval: null
  # -- ServiceMonitor scrape timeout in Go duration format (e.g. 15s)
  scrapeTimeout: null

# -- Configure additional ports and services. For each configured port, a corresponding service is created.
# See values.yaml for details
extraPorts: {}
#  syslog:
#    name: tcp-syslog
#    containerPort: 1514
#    protocol: TCP
#    service:
#      type: ClusterIP
#      clusterIP: null
#      port: 1514
#      externalIPs: []
#      nodePort: null
#      annotations: {}
#      labels: {}
#      loadBalancerIP: null
#      loadBalancerSourceRanges: []
#      externalTrafficPolicy: null

# -- PodSecurityPolicy configuration.
# @default -- See "values.yaml"
podSecurityPolicy:
  privileged: true
  allowPrivilegeEscalation: true
  volumes:
    - 'secret'
    - 'hostPath'
    - 'downwardAPI'
  hostNetwork: false
  hostIPC: false
  hostPID: false
  runAsUser:
    rule: 'RunAsAny'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
  readOnlyRootFilesystem: true
  requiredDropCapabilities:
    - ALL

# -- Section for crafting Promtails config file. The only directly relevant value is "config.file"
# which is a templated string that references the other values and snippets below this key.
# @default -- See "values.yaml"
config:
  # -- The port of the Promtail server
  # Must be reference in "config.file" to configure "server.http_listen_port"
  # See default config in "values.yaml"
  serverPort: 3101
  # -- The Loki address to post logs to.
  # Must be reference in "config.file" to configure "client.url".
  # See default config in "values.yaml"
  lokiAddress: http://loki:3100/loki/api/v1/push
  # -- A section of reusable snippets that can be reference in "config.file".
  # Custom snippets may be added in order to reduce redundancy.
  # This is especially helpful when multiple "kubernetes_sd_configs" are use which usually have large parts in common.
  # @default -- See "values.yaml"
  snippets:
    pipelineStages:
      - cri: {}
    common:
      - action: replace
        source_labels:
          - __meta_kubernetes_pod_node_name
        target_label: node_name
      - action: replace
        source_labels:
          - __meta_kubernetes_namespace
        target_label: namespace
      - action: replace
        replacement: $1
        separator: /
        source_labels:
          - namespace
          - app
        target_label: job
      - action: replace
        source_labels:
          - __meta_kubernetes_pod_name
        target_label: pod
      - action: replace
        source_labels:
          - __meta_kubernetes_pod_container_name
        target_label: container
      - action: replace
        replacement: /var/log/pods/*$1/*.log
        separator: /
        source_labels:
          - __meta_kubernetes_pod_uid
          - __meta_kubernetes_pod_container_name
        target_label: __path__
      - action: replace
        replacement: /var/log/pods/*$1/*.log
        separator: /
        source_labels:
          - __meta_kubernetes_pod_annotation_kubernetes_io_config_hash
          - __meta_kubernetes_pod_container_name
        target_label: __path__

    # If set to true, adds an additional label for the scrape job.
    # This helps debug the Promtail config.
    addScrapeJobLabel: false

    # -- You can put here any keys that will be directly added to the config file's 'client' block.
    # @default -- empty
    extraClientConfigs: ""

    # -- You can put here any additional scrape configs you want to add to the config file.
    # @default -- empty
    extraScrapeConfigs: ""

    scrapeConfigs: |
      # See also https://github.com/grafana/loki/blob/master/production/ksonnet/promtail/scrape_config.libsonnet for reference

      # Pods with a label 'app.kubernetes.io/name'
      - job_name: kubernetes-pods-app-kubernetes-io-name
        pipeline_stages:
          {{- toYaml .Values.config.snippets.pipelineStages | nindent 4 }}
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          - action: replace
            source_labels:
              - __meta_kubernetes_pod_label_app_kubernetes_io_name
            target_label: app
          - action: drop
            regex: ''
            source_labels:
              - app
          - action: replace
            source_labels:
              - __meta_kubernetes_pod_label_app_kubernetes_io_component
            target_label: component
          {{- if .Values.config.snippets.addScrapeJobLabel }}
          - action: replace
            replacement: kubernetes-pods-app-kubernetes-io-name
            target_label: scrape_job
          {{- end }}
          {{- toYaml .Values.config.snippets.common | nindent 4 }}

      # Pods with a label 'app'
      - job_name: kubernetes-pods-app
        pipeline_stages:
          {{- toYaml .Values.config.snippets.pipelineStages | nindent 4 }}
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          # Drop pods with label 'app.kubernetes.io/name'. They are already considered above
          - action: drop
            regex: .+
            source_labels:
              - __meta_kubernetes_pod_label_app_kubernetes_io_name
          - action: replace
            source_labels:
              - __meta_kubernetes_pod_label_app
            target_label: app
          - action: drop
            regex: ''
            source_labels:
              - app
          - action: replace
            source_labels:
              - __meta_kubernetes_pod_label_component
            target_label: component
          {{- if .Values.config.snippets.addScrapeJobLabel }}
          - action: replace
            replacement: kubernetes-pods-app
            target_label: scrape_job
          {{- end }}
          {{- toYaml .Values.config.snippets.common | nindent 4 }}

      # Pods with direct controllers, such as StatefulSet
      - job_name: kubernetes-pods-direct-controllers
        pipeline_stages:
          {{- toYaml .Values.config.snippets.pipelineStages | nindent 4 }}
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          # Drop pods with label 'app.kubernetes.io/name' or 'app'. They are already considered above
          - action: drop
            regex: .+
            separator: ''
            source_labels:
              - __meta_kubernetes_pod_label_app_kubernetes_io_name
              - __meta_kubernetes_pod_label_app
          - action: drop
            regex: '[0-9a-z-.]+-[0-9a-f]{8,10}'
            source_labels:
              - __meta_kubernetes_pod_controller_name
          - action: replace
            source_labels:
              - __meta_kubernetes_pod_controller_name
            target_label: app
          {{- if .Values.config.snippets.addScrapeJobLabel }}
          - action: replace
            replacement: kubernetes-pods-direct-controllers
            target_label: scrape_job
          {{- end }}
          {{- toYaml .Values.config.snippets.common | nindent 4 }}

      # Pods with indirect controllers, such as Deployment
      - job_name: kubernetes-pods-indirect-controller
        pipeline_stages:
          {{- toYaml .Values.config.snippets.pipelineStages | nindent 4 }}
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          # Drop pods with label 'app.kubernetes.io/name' or 'app'. They are already considered above
          - action: drop
            regex: .+
            separator: ''
            source_labels:
              - __meta_kubernetes_pod_label_app_kubernetes_io_name
              - __meta_kubernetes_pod_label_app
          - action: keep
            regex: '[0-9a-z-.]+-[0-9a-f]{8,10}'
            source_labels:
              - __meta_kubernetes_pod_controller_name
          - action: replace
            regex: '([0-9a-z-.]+)-[0-9a-f]{8,10}'
            source_labels:
              - __meta_kubernetes_pod_controller_name
            target_label: app
          {{- if .Values.config.snippets.addScrapeJobLabel }}
          - action: replace
            replacement: kubernetes-pods-indirect-controller
            target_label: scrape_job
          {{- end }}
          {{- toYaml .Values.config.snippets.common | nindent 4 }}
      # All remaining pods not yet covered
      - job_name: kubernetes-other
        pipeline_stages:
          {{- toYaml .Values.config.snippets.pipelineStages | nindent 4 }}
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          # Drop what has already been covered
          - action: drop
            regex: .+
            separator: ''
            source_labels:
              - __meta_kubernetes_pod_label_app_kubernetes_io_name
              - __meta_kubernetes_pod_label_app
          - action: drop
            regex: .+
            source_labels:
              - __meta_kubernetes_pod_controller_name
          - action: replace
            source_labels:
              - __meta_kubernetes_pod_name
            target_label: app
          - action: replace
            source_labels:
              - __meta_kubernetes_pod_label_component
            target_label: component
          {{- if .Values.config.snippets.addScrapeJobLabel }}
          - action: replace
            replacement: kubernetes-other
            target_label: scrape_job
          {{- end }}
          {{- toYaml .Values.config.snippets.common | nindent 4 }}

  # -- Config file contents for Promtail.
  # Must be configured as string.
  # It is templated so it can be assembled from reusable snippets in order to avoid redundancy.
  # @default -- See "values.yaml"
  file: |
    server:
      log_level: info
      http_listen_port: {{ .Values.config.serverPort }}

    client:
      url: {{ .Values.config.lokiAddress }}
      {{- tpl .Values.config.snippets.extraClientConfigs . | nindent 2 }}

    positions:
      filename: /run/promtail/positions.yaml

    scrape_configs:
      {{- tpl .Values.config.snippets.scrapeConfigs . | nindent 2 }}
      {{- tpl .Values.config.snippets.extraScrapeConfigs . | nindent 2 }}
`
