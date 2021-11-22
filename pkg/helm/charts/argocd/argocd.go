// Package argocd provides a Helm chart for installing:
// - https://argoproj.github.io/argo-cd/
// - https://github.com/argoproj/argo-helm
package argocd

import (
	"bytes"
	_ "embed"
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

//go:embed values-template.yaml
//nolint:gochecknoglobals
var valuesTemplate string

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
