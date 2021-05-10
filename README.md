[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md) [![codecov](https://codecov.io/gh/oslokommune/okctl/branch/master/graph/badge.svg)](https://codecov.io/gh/oslokommune/okctl) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=oslokommune_okctl&metric=security_rating)](https://sonarcloud.io/dashboard?id=oslokommune_okctl) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=oslokommune_okctl&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=oslokommune_okctl) [![Go Report Card](https://goreportcard.com/badge/github.com/oslokommune/okctl)](https://goreportcard.com/report/github.com/oslokommune/okctl) [![Nightly build](https://github.com/oslokommune/okctl/actions/workflows/nightly.yaml/badge.svg)](https://github.com/oslokommune/okctl/actions/workflows/nightly.yaml)

# `okctl` - Opinionated and effortless infrastructure and application management

<span style="display:block;text-align:center">![okctl](logo/okctl-gopher.png)</span>

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

## Getting started

The following is a guide for how to create an environment that contains the elements described in
[Compare and contrast](#compare-and-contrast).

### 1. Create a new GitHub repository

Go to [github.com/oslokommune](https://github.com/oslokommune) and create a new private git repository. No credentials are stored in this
repository, we want it private as a precaution until we are certain it is safe to have it set as public.

:information_source: This repository will be used by okctl to store infrastructure-as-code, which means files containing
various configuration for your up-and-coming cluster.

Now, run:

```bash
# Clone the repository you just made
$ git clone git@github.com:oslokommune/<the new repository>.git
$ cd <the new repository>
```

### 2. Create a cluster

A "cluster" is a [Kubernetes](https://kubernetes.io) cluster with many addons and integrations, creating a production grade environment as described in [Functionality](#functionality).

You will soon be running `okctl apply cluster`, which will ask you for the following information:

* Username and password: This is your Oslo Kommune AD organization username (e.g., oooXXXXX) and its password.
* Multi factor token (MFA): The same one you use to login to
[AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws). If you haven't set up MFA yet, you
can do that [here](https://login.oslo.kommune.no/auth/realms/AD/account/totp).
* AWS account ID: This identifies which account you want to use. You can see which accounts you have access to just
after logging in to [AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws):

![okctl](userdocs/src/img/aws-account-id.png)

```bash
# Scaffold a cluster. Format:
okctl scaffold cluster -f cluster.yaml
# <edit cluster.yaml>
okctl apply cluster -f cluster.yaml
```

Follow the instructions.

When done, verify that you have a working cluster by running

```bash
$ okctl venv -c cluster.yaml
$ kubectl get service

```

The last command should show something like

```bash
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.100.0.1   <none>        443/TCP   1h
```

## Common commands

```bash
# Get help for any command
$ okctl --help

# Show credentials for cluster
$ okctl show credentials -c cluster.yaml

# Run a sub shell with environment variables from the above command and a custom command prompt (PS1)
$ okctl venv -c cluster.yaml

# Delete the cluster
$ okctl delete cluster -c cluster.yaml
```

## Functionality

The core cluster is up and running, and we are currently working on building a seamless experience on top of this cluster with integrations for common functionality:

### Core cluster

- [x] [ExternalSecrets](https://external-secrets.github.io/kubernetes-external-secrets/) for storing secrets securely
- [x] [AWS Load Balancer Controller](https://github.com/kubernetes-sigs/aws-load-balancer-controller) creates load balancers for incoming traffic
- [x] [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) ensures humane DNS hostnames 
- [x] [Autoscaler](https://github.com/kubernetes/autoscaler/) for adjusting the size of pods and nodes
- [x] [Argo CD](https://github.com/argoproj/argo-cd) gives us continuous delivery
- [x] [Kubernetes Prometheus](https://github.com/prometheus-operator/kube-prometheus) for prometheus and grafana
- [x] [Amazon Elastic Block Store (EBS) CSI driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/) for block storage
- [x] [Loki](https://github.com/grafana/loki) for processing logs
- [x] [Tempo](https://grafana.com/oss/tempo/) for processing traces
- [x] [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) for collecting logs

### Application lifecycle

- [x] [Postgres](https://aws.amazon.com/rds/postgresql/) for creating and integrating a postgres database with your application
- [x] [Amazon Elastic Container Registry](https://aws.amazon.com/ecr/) for creating and assisting with the publication of container images for use in the cluster
- [ ] Reference application that demonstrates how to use the cluster and its integrations

## Compare and contrast

The intended purpose of `okctl` is to be an opinionated solver of infrastructure problems, this includes setting up CI/CD pipelines et al.

The following table is not present to determine what tool is better or worse, but rather how these tools compare to `okctl` and the problems we are focused on solving.

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
