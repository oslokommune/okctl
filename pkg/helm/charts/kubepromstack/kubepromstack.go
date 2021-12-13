// Package kubepromstack provides a helm chart for installing kube-prometheus-stack:
// - https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
package kubepromstack

import (
	"bytes"
	_ "embed"
	"text/template"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

const (
	// ReleaseName is the name of the release
	ReleaseName = "kube-prometheus-stack"
	// Namespace is the default namespace
	Namespace = "monitoring"
)

// New returns an initialised helm chart:
// - https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
func New(values *Values, timeout time.Duration) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "prometheus-community",
		RepositoryURL:  "https://prometheus-community.github.io/helm-charts",
		ReleaseName:    ReleaseName,
		Version:        "13.9.1",
		Chart:          "kube-prometheus-stack",
		Namespace:      Namespace,
		Timeout:        timeout,
		Values:         values,
	}
}

// Values maps up the values.yaml file
type Values struct {
	GrafanaServiceAccountName          string
	GrafanaCertificateARN              string
	GrafanaHostname                    string
	AuthHostname                       string
	ClientID                           string
	SecretsConfigName                  string
	SecretsGrafanaCookieSecretKey      string
	SecretsGrafanaOauthClientSecretKey string
	SecretsGrafanaAdminUserKey         string
	SecretsGrafanaAdminPassKey         string
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

//go:embed values.yaml
var valuesTemplate string //nolint:gochecknoglobals
