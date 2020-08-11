[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md) [![codecov](https://codecov.io/gh/oslokommune/okctl/branch/master/graph/badge.svg)](https://codecov.io/gh/oslokommune/okctl) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=oslokommune_okctl&metric=security_rating)](https://sonarcloud.io/dashboard?id=oslokommune_okctl) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=oslokommune_okctl&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=oslokommune_okctl) [![Go Report Card](https://goreportcard.com/badge/github.com/oslokommune/okctl)](https://goreportcard.com/report/github.com/oslokommune/okctl)

# `okctl` - Opinionated and effortless infrastructure and application management

We will be building up the functionality of this CLI within the coming weeks, and the aim is to release frequently. At this stage we are mostly setting the CI/CD and release processes.

![okctl](logo/okctl-fit.png)

## Installation

To download the latest release, run the command matching your operating system:

```bash
# Linux
curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/okctl /usr/local/bin

# macOS
brew tap oslokommune/tap
brew install oslokommune/tap/okctl
```

## Usage

```bash
$ okctl --help

# Create a cluster
$ okctl create cluster pro 123456789012

# Show credentials for cluster
$ okctl show credentials pro

# Delete the cluster
$ okctl delete cluster pro
```

## Roadmap

We have a basic cluster up and running, and we are currently working on getting the basic components of the cluster deployed:

- [ ] [SealedSecrets](https://github.com/bitnami-labs/sealed-secrets) for storing secrets securely in git
- [ ] [AWS ALB Ingress Controller](https://github.com/kubernetes-sigs/aws-alb-ingress-controller) creates load balancers for incoming traffic
- [ ] [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) ensures humane DNS hostnames 
- [ ] [Autoscaler](https://github.com/kubernetes/autoscaler/) for adjusting the size of pods and nodes
- [ ] [Argo CD](https://github.com/argoproj/argo-cd) gives us continuous delivery
- [ ] [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator) for monitoring
- [ ] [Amazon Elastic Block Store (EBS) CSI driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/) for block storage
- [ ] [Loki](https://github.com/grafana/loki) for consuming logs

## Compare and contrast

`okctl` intended purpose is to be an opinionated solver of infrastructure problems, this includes setting up CI/CD pipelines et al.

The purpose of this table is not to determine what tool is better or worse, but rather how these tools compare to `okctl` with the problem they are trying to solve.

| | okctl | [eksctl](https://eksctl.io) | [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) | [serverless.tf](https://serverless.tf) | 
|---|---|---|---|---|
|Defines a CI/CD scheme|✅|❌|❌|❌|
|Creates a Kubernetes cluster|✅|✅|❌|❌|
|Facilitates with application creation|✅|❌|❌|✅|
|Integrates with Github (actions, packages, oauth)|✅|❌|❌|❌|
|Sets up monitoring|✅|❌|❌|❌|

## Inspiration

We have begged, borrowed and stolen various ideas from the following CLIs:

- https://github.com/deislabs/porter
- https://github.com/weaveworks/eksctl
- https://github.com/kubernetes/kops
- https://kubernetes.io/docs/tasks/tools/install-kubectl/
- https://github.com/go-kit/kit
- https://github.com/Versent/saml2aws/