// Package blockstorage provides a Helm chart installing:
// - https://github.com/kubernetes-sigs/aws-ebs-csi-driver
package blockstorage

import (
	"bytes"
	"text/template"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

const (
	// ReleaseName is the name of the release
	ReleaseName = "aws-ebs-csi-driver"
	// Namespace is the default namespace
	Namespace = "kube-system"
)

// New returns an initialised Helm chart for installing aws-ebs-csi-driver
func New(values *Values, timeout time.Duration) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "aws-ebs-csi-driver",
		RepositoryURL:  "https://kubernetes-sigs.github.io/aws-ebs-csi-driver",
		ReleaseName:    ReleaseName,
		Version:        "0.9.6",
		Chart:          "aws-ebs-csi-driver",
		Namespace:      Namespace,
		Timeout:        timeout,
		Values:         values,
	}
}

// NewDefaultValues returns the mapped values.yml containing the default
// values
func NewDefaultValues(region, clusterName, serviceAccount string) *Values {
	return &Values{
		Region:         region,
		ClusterName:    clusterName,
		ServiceAccount: serviceAccount,
	}
}

// Values contains the configurable values
type Values struct {
	Region         string
	ClusterName    string
	ServiceAccount string
}

// RawYAML implements the helm marshaller interface
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

const valuesTemplate = `# Default values for aws-ebs-csi-driver.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 2

image:
  repository: k8s.gcr.io/provider-aws/aws-ebs-csi-driver
  tag: "v0.9.0"
  pullPolicy: IfNotPresent

sidecars:
  provisionerImage:
    repository: k8s.gcr.io/sig-storage/csi-provisioner
    tag: "v2.0.2"
  attacherImage:
    repository: k8s.gcr.io/sig-storage/csi-attacher
    tag: "v3.0.0"
  snapshotterImage:
    repository: k8s.gcr.io/sig-storage/csi-snapshotter
    tag: "v3.0.3"
  livenessProbeImage:
    repository: k8s.gcr.io/sig-storage/livenessprobe
    tag: "v2.1.0"
  resizerImage:
    repository: k8s.gcr.io/sig-storage/csi-resizer
    tag: "v1.0.0"
  nodeDriverRegistrarImage:
    repository: k8s.gcr.io/sig-storage/csi-node-driver-registrar
    tag: "v2.0.1"

snapshotController:
  repository: k8s.gcr.io/sig-storage/snapshot-controller
  tag: "v3.0.3"

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

podAnnotations: {}

# True if enable volume scheduling for dynamic volume provisioning
enableVolumeScheduling: true

# True if enable volume resizing
enableVolumeResizing: true

# True if enable volume snapshot
enableVolumeSnapshot: true

# The "maximum number of attachable volumes" per node
volumeAttachLimit: ""

resources:
  {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

priorityClassName: ""
nodeSelector: {}
tolerateAllTaints: true
tolerations: []
affinity: {}

# Extra volume tags to attach to each dynamically provisioned volume.
# ---
# extraVolumeTags:
#   key1: value1
#   key2: value2
extraVolumeTags: {}

# If set, add pv/pvc metadata to plugin create requests as parameters.
extraCreateMetadata: true

# ID of the Kubernetes cluster used for tagging provisioned EBS volumes (optional).
k8sTagClusterId: {{.ClusterName}}

# AWS region to use. If not specified then the region will be looked up via the AWS EC2 metadata
# service.
# ---
# region: us-east-1
region: {{.Region}}

node:
  priorityClassName: ""
  nodeSelector: {}
  podAnnotations: {}
  tolerateAllTaints: true
  tolerations: []
  resources: {}

serviceAccount:
  controller:
    create: false # A service account will be created for you if set to true. Set to false if you want to use your own.
    name: {{.ServiceAccount}} # Name of the service-account to be used/created.
    annotations: {}
  snapshot:
    create: false
    name: {{.ServiceAccount}}
    annotations: {}
  node:
    create: false
    name: {{.ServiceAccount}}
    annotations: {}

storageClasses: []
`
