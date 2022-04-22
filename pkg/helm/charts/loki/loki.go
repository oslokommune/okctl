// Package loki provides a Helm chart for installing:
// - https://github.com/grafana/helm-charts/tree/main/charts/loki
// - https://grafana.com/oss/loki/
package loki

import (
	"bytes"
	_ "embed"
	"text/template"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

const (
	// ReleaseName is the name of the release
	ReleaseName = "loki"
	// Namespace is the default namespace
	Namespace = "monitoring"
)

// New returns an initialised Helm chart for installing cluster-Loki
func New(values *Values, timeout time.Duration) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "grafana",
		RepositoryURL:  "https://grafana.github.io/helm-charts",
		ReleaseName:    ReleaseName,
		Version:        "2.3.0",
		Chart:          "loki",
		Namespace:      Namespace,
		Timeout:        timeout,
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

//nolint:gochecknoglobals
//go:embed values.yaml
var valuesTemplate string
