package commands

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

// ScaffoldClusterOpts contains attributes used in a cluster declaration
type ScaffoldClusterOpts struct {
	Name string

	AWSAccountID    string
	Organization    string
	RepositoryName  string
	OutputDirectory string
}

// ScaffoldClusterDeclaration scaffolds a cluster declaration based on ScaffoldClusterOpts
func ScaffoldClusterDeclaration(out io.Writer, opts ScaffoldClusterOpts) error {
	sanitizedOpts := sanitizeOpts(opts)

	t, err := template.New("cluster.yaml").Parse(clusterTemplate)
	if err != nil {
		return fmt.Errorf("parsing template string: %w", err)
	}

	err = t.Execute(out, sanitizedOpts)
	if err != nil {
		return fmt.Errorf("interpolating template: %w", err)
	}

	return nil
}

func sanitizeOpts(opts ScaffoldClusterOpts) ScaffoldClusterOpts {
	return ScaffoldClusterOpts{
		Name:            strings.ToLower(opts.Name),
		AWSAccountID:    opts.AWSAccountID,
		Organization:    opts.Organization,
		OutputDirectory: opts.OutputDirectory,
		RepositoryName:  opts.RepositoryName,
	}
}

const clusterTemplate = `apiVersion: okctl.io/v1alpha1
kind: Cluster

## For help finding values, see https://okctl.io/cluster/reference

metadata:
  ## Add your AWS account ID
  accountID: '{{ .AWSAccountID }}'

  ## Create a cluster name that defines the scope of the cluster, typically the name of the product or your team
  name: {{ .Name }}

  ## Choose the AWS region to prefer when creating resources
  # region: eu-west-1

## Create a cluster root domain name
clusterRootDomain: {{ .Name }}.oslo.systems

## To be able to setup ArgoCD correctly, okctl needs to know which repository on Github that will contain your infrastructure
github:
  ## Specify your IAC repository name
  repository: {{ .RepositoryName }}

  ## Specify the organization that owns the repository
  # organisation: {{ .Organization }}

  ## Specify the okctl root directory name
  # outputPath: {{ .OutputDirectory }}

## Add admin users to your okctl environment
#users:
#- email: user.email@emailprovider.org

#integrations:
#  # Enable ArgoCD, Git based infrastructure as code (IAC)
#  argoCD: true
#
#  # Enable Autoscaler, automatic scaling of pods and nodes
#  autoscaler: true
#
#  # Enable AWS Load Balancer Controller, automatic provisioning of load balancers
#  awsLoadBalancerController: true
#
#  # Enable Block Storage, persistent storage for pods
#  blockstorage: true
#
#  # Enable Cognito, authentication provider for ArgoCD and Grafana
#  cognito: true
#
#  # Enable ExternalDNS, automatic DNS configuration
#  externalDNS: true
#
#  # Enable External Secrets, external referencing of secrets
#  externalSecrets: true
#
#  # Enable Prometheus and Grafana, for metrics
#  kubePromStack: true
#
#  # Enable Promtail, feeding pod logs to Loki
#  promtail: true
#
#  # Enable Loki, exposing pod logs in Grafana
#  loki: true
#
#  # Enable Tempo, exposing traces in Grafana
#  tempo: true

## Specify the database(s) for okctl to provision
## namespace is for placing database information
## user defines the admin user to provision for operations
#databases:
#  postgres:
#  - name: dbname
#    namespace: relevantnamespace
#    user: postgres

## For help finding values, see https://okctl.io/cluster/reference
`
