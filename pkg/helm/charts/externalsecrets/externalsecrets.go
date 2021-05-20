// Package externalsecrets provides a helm chart for installing external-secrets:
// - https://external-secrets.github.io/kubernetes-external-secrets
package externalsecrets

import (
	"bytes"
	"text/template"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

const (
	// ReleaseName is the name of the release
	ReleaseName = "external-secrets"
	// Namespace is the default namespace
	Namespace = "kube-system"
)

// New returns an initialised external secrets chart
// - https://github.com/external-secrets/kubernetes-external-secrets/blob/master/charts/kubernetes-external-secrets/README.md
func New(values *Values, timeout time.Duration) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "external-secrets",
		RepositoryURL:  "https://external-secrets.github.io/kubernetes-external-secrets/",
		ReleaseName:    ReleaseName,
		Version:        "6.4.0",
		Chart:          "kubernetes-external-secrets",
		Namespace:      Namespace,
		Timeout:        timeout,
		Values:         values,
	}
}

// NewDefaultValues returns the default secrets values
// for the helm chart
func NewDefaultValues(region string) *Values {
	return &Values{
		Region: region,
	}
}

// Values maps up the values.yaml file
// nolint: maligned
type Values struct {
	Region string
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

const valuesTemplate = `crds:
    create: false
customResourceManagerDisabled: true
env:
    AWS_REGION: {{.Region}}
    AWS_DEFAULT_REGION: {{.Region}}
    POLLER_INTERVAL_MILLISECONDS: 10000
    LOG_LEVEL: debug
    LOG_MESSAGE_KEY: msg
    METRICS_PORT: 3001
    USE_HUMAN_READABLE_LOG_LEVELS: true
image:
    repository: ghcr.io/external-secrets/kubernetes-external-secrets
    tag: 6.4.0
    pullPolicy: IfNotPresent
rbac:
    create: true
replicaCount: 1
resources:
    limits:
        cpu: 100m
        memory: 128Mi
    requests:
        cpu: 100m
        memory: 128Mi
securityContext:
    runAsNonRoot: true
    fsGroup: 65534
serviceAccount:
    create: false
    name: external-secrets
serviceMonitor:
    enabled: false
    interval: 30s
    namespace: ""
`
