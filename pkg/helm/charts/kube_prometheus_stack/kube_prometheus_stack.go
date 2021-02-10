// Package kube_prometheus_stack provides a helm chart for installing kube-prometheus-stack:
// - https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
package kube_prometheus_stack

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/oslokommune/okctl/pkg/helm"
)

// TODO check that these values/structs matches with actual values.yaml

// Values maps up the values.yaml file
// nolint: maligned
type Values struct {
	CRDs                          CRDs            `yaml:"crds"`
	CustomResourceManagerDisabled bool            `yaml:"customResourceManagerDisabled"`
	Env                           Env             `yaml:"env"`
	Image                         Image           `yaml:"image"`
	RBAC                          RBAC            `yaml:"rbac"`
	ReplicaCount                  int             `yaml:"replicaCount"`
	Resources                     Resources       `yaml:"resources"`
	SecurityContext               SecurityContext `yaml:"securityContext"`
	ServiceAccount                ServiceAccount  `yaml:"serviceAccount"`
	ServiceMonitor                ServiceMonitor  `yaml:"serviceMonitor"`
}

// CRDs determines how CRDs are created
type CRDs struct {
	Create bool `yaml:"create"`
}

// Env sets the pod environment variables
type Env struct {
	AwsRegion                  string `yaml:"AWS_REGION"`
	AwsDefaultRegion           string `yaml:"AWS_DEFAULT_REGION"`
	PollerIntervalMilliseconds int    `yaml:"POLLER_INTERVAL_MILLISECONDS"`
	LogLevel                   string `yaml:"LOG_LEVEL"`
	LogMessageKey              string `yaml:"LOG_MESSAGE_KEY"`
	MetricsPort                int    `yaml:"METRICS_PORT"`
	UseHumanReadableLogLevels  bool   `yaml:"USE_HUMAN_READABLE_LOG_LEVELS"`
}

// RBAC determines whether RBAC resources
// are created or not
type RBAC struct {
	Create bool `yaml:"create"`
}

// ServiceAccount determines how service
// accounts are created
type ServiceAccount struct {
	Create bool   `yaml:"create"`
	Name   string `yaml:"name"`
}

// Image set what container image
// to use
type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

// SecurityContext sets the pod security context
type SecurityContext struct {
	RunAsNonRoot bool `yaml:"runAsNonRoot"`
	FsGroup      int  `yaml:"fsGroup"`
}

// Resources set the resource consumption of the pod
type Resources struct {
	Limits   ResourceEntry `yaml:"limits"`
	Requests ResourceEntry `yaml:"requests"`
}

// ResourceEntry sets the values for cpu and memory
type ResourceEntry struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// ServiceMonitor determines if the service monitor
// is enabled or not
type ServiceMonitor struct {
	Enabled   bool   `yaml:"enabled"`
	Interval  string `yaml:"interval"`
	Namespace string `yaml:"namespace"`
}

// ExternalSecrets returns an initialised external secrets chart
// - https://github.com/external-secrets/kubernetes-external-secrets/blob/master/charts/kubernetes-external-secrets/README.md
func KubePrometheusStack(values interface{}) *helm.Chart {
	return &helm.Chart{
		RepositoryName: "external-secrets",
		RepositoryURL:  "https://external-secrets.github.io/kubernetes-external-secrets/",
		ReleaseName:    "external-secrets",
		Version:        "5.2.0",
		Chart:          "kubernetes-external-secrets",
		Namespace:      "kube-system",
		Timeout:        3 * time.Minute, // nolint: gomnd
		Values:         values,
	}
}

// DefaultKubePrometheusStackValues returns the default prometheus values
// for the helm chart
func DefaultKubePrometheusStackValues() (string, error) {
	result, err := goTemplateToString(kubePrometheusStackTemplate, struct {}{})

	if err != nil {
	    return "", fmt.Errorf("parsing KubePrometheusStackValues template: %w", err)
	}
	
	return result, nil
}

func goTemplateToString(templateString string, data interface{}) (string, error) {
	tmpl, err := template.New("t").Parse(templateString)
	if err != nil {
		return "", err
	}

	tmplBuffer := new(bytes.Buffer)
	err = tmpl.Execute(tmplBuffer, data)

	if err != nil {
		return "", err
	}

	return tmplBuffer.String(), nil
}
