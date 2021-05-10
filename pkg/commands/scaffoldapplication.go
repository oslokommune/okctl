package commands

import (
	"fmt"
	"io"
	"text/template"
)

// ScaffoldApplicationOpts contains attributes used in an application declaration
type ScaffoldApplicationOpts struct {
	PrimaryHostedZone string
}

// ScaffoldApplicationDeclaration scaffolds an application declaration based on ScaffoldApplicationOpts
func ScaffoldApplicationDeclaration(out io.Writer, opts ScaffoldApplicationOpts) error {
	t, err := template.New("application.yaml").Parse(applicationTemplate)
	if err != nil {
		return fmt.Errorf("parsing template string: %w", err)
	}

	err = t.Execute(out, opts)
	if err != nil {
		return fmt.Errorf("interpolating template: %w", err)
	}

	return nil
}

const applicationTemplate = `apiVersion: okctl.io/v1alpha1
kind: Application

metadata:
  # A name that identifies your app
  name: my-app
  # The Kubernetes namespace where your app will live
  namespace: my-namespace

# The Docker image containing the application. image.uri and image.name is mutually exclusive. Either specify the URI or
# define a name of an ECR repository for which okctl will create for you.
image: 
  # uri defines where the image can be pulled from
  uri: ghcr.io/oslokommune/okctl-hello:v0.0.3
  # name defines the name of an ECR repository to create
  #name: my-app

# The subdomain of the URL your app should be available on
# Example in a cluster with {{ .PrimaryHostedZone }} as root cluster URL (as defined by primary DNS zone in the
# cluster declaration):
#
# subDomain: okctl
# result: okctl.{{ .PrimaryHostedZone }}
# Comment this out to avoid setting up an ingress, in other words - avoid exposing it on the internet
#
subDomain: my-app

# The port your app listens on
# Comment this out to avoid setting up a service (required if url is specified)
#
port: 3000

# Enable prometheus scraping of metrics
#prometheus:
#  path: /metrics

# How many replicas of your application should we scaffold
#replicas: 3 # 1 by default

# A Docker repository secret for pulling your image
#imagePullSecret: my-pull-secret-name

# The environment your app requires
#environment:
#  MY_VARIABLE: my-value

# Volumes to mount
#volumes:
#  - /path/to/mount/volume: # Requests 1Gi by default
#  - /path/to/mount/volume: 24Gi

# Annotations for your ingress
#ingress:
#  annotations:
#    nginx.ingress.kubernetes.io/cors-allow-origin: http://localhost:8080
#    cert-manager.io/cluster-issuer: letsencrypt-production
`
