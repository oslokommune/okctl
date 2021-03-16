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
	Environment     string
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
		Environment:     strings.ToLower(opts.Environment),
		Organization:    opts.Organization,
		RepositoryName:  opts.RepositoryName,
		OutputDirectory: opts.OutputDirectory,
	}
}

const clusterTemplate = `apiVersion: okctl.io/v1alpha1
kind: Cluster

## For help finding values, see https://okctl.io/getting-started/create-cluster
metadata:
  ## Account ID is your AWS account ID
  accountID: '{{ .AWSAccountID }}'

  ## Environment is the name you use to identify the type of cluster it is. Common names are production, test, staging
  environment: {{ .Environment }}

  ## Name can be anything, but should define the scope of the cluster. Meaning if the cluster is scoped to one product,
  ## you might want to name it the name of the product. If the cluster contains all services and products owned by a
  ## team, the team name might be more fitting.
  name: {{ .Name }}

  ## Region defines the AWS region to prefer when creating resources
  # region: eu-west-1

## The cluster root domain defines the domain of which to create services beneath. For example; okctl will setup ArgoCD
## which has a frontend. The frontend will be available at https://argocd.<clusterRootDomain>. For Cognito it will be 
## https://auth.<clusterRootDomain>
clusterRootDomain: {{ .Name }}-{{ .Environment }}.oslo.systems

## For okctl to be able to setup ArgoCD correctly for you, it needs to know what repository on Github that will contain
## your infrastructure.
github:
  ## The name of the repository
  repository: {{ .RepositoryName }}

  ## The organization that owns the repository
  # organisation: {{ .Organization }}

  ## The folder to place infrastructure declarations
  # outputPath: {{ .OutputDirectory }}

integrations:
  ## ArgoCD is a service that watches a repository for Kubernetes charts and ensures the defined resources are running
  ## as declared in the cluster
  argoCD: true

  ## Autoscaler automatically adjusts the size of pods and nodes in your cluster depending on load
  autoscaler: true

  ## AWS Load Balancer Controller handles routing from the internet to your application running inside your okctl
  ## Kubernetes cluster. If you want your applications and services accessible from the internet, this needs to be
  ## enabled.
  awsLoadBalancerController: true

  ## Block storage provides persistent storage for your cluster (Persistent Volumes)
  blockstorage: true

  ## Cognito is an authentication provider that okctl uses to control access to different resources, like ArgoCD and
  ## Grafana
  cognito: true

  ## External DNS handles defining the necessary DNS records required to route traffic to your defined service or 
  ## application
  externalDNS: true

  ## External Secrets fetches secrets from external sources and exposes them as native Kubernetes secrets inside the
  ## cluster 
  externalSecrets: true

  ## KubePromStack enables Prometheus and Grafana for metrics
  kubePromStack: true

  ## Loki collects logs and exposes them as a data source in Grafana
  loki: true

  ## Promtail scrapes logs from pods and feeds them to Loki
  promtail: true

  ## Tempo collects traces and exposes them as a data source in Grafana. Supports formats like jaeger, zipkin, open 
  ## telemetry
  tempo: true

## okctl creates a Virtual Private Cloud for you which it organizes all the intended resources that require networking.
## A VPC is mandatory, but can be configured by the following attributes.
#vpc:
  ## CIDR defines the VPC IP range. Leave this be if you don't know what it is/does
  # cidr: 192.168.0.0/20

  ## HighAvailability means we create redundancy in the network setup. If set to true we will create a NAT gateway per
  ## public subnet, instead of routing all traffic through one.
  # highAvailability: true
`
