# `okctl` - Opinionated and effortless infrastructure and application management

`okctl` is a CLI for creating a production environment on [Amazon Web Services](https://aws.amazon.com/) (AWS) using [Elastic Kubernetes Service](https://aws.amazon.com/eks/) (EKS) as its foundation, we aim to be:

- **Specialised**: by tightly integrating with the services provided by Oslo kommune
- **Highly opinionated**: you will get something that works well that may not be exactly what you want
- **Modern**: base our decisions on best-practise methodologies for a modern development lifecycle
- **Sensible**: provide a set of sensible default integrations

<span style="display:block;text-align:center">![okctl](img/logo.png)</span>

## Installation

To download the latest release, run the command matching your operating system:

### Linux

```bash
curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/okctl /usr/local/bin
```

### macOS

```bash
brew tap oslokommune/tap
brew install oslokommune/tap/okctl
```

## Getting started

The following is a guide for how to create an environment that contains the elements described in
[Compare and contrast](#compare-and-contrast).

### 1. Create a new github repository

Go to github.com/oslokommune and create a new **private** git repository.

:information_source: This repository will be used by okctl to store infrastructure-as-code, which means files containing
various configuration for your up-and-coming cluster.

Now, run:

```bash

# Clone the repository you just made
$ git clone git@github.com:oslokommune/<the new repository>.git
$ cd <the new repository>
```

### 2. Create a cluster

A "cluster" is a Kubernetes cluster in an environment as described in [Compare and contrast](#compare-and-contrast).

You will soon be running `okctl create cluster`, which will ask you for the following information:
* Username and password: This is your Oslo Kommune AD organization username (oooXXXXX) and its password.
* Multi factor token (MFA): The same one you use to login to
[AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws). If you haven't set up MFA yet, you
can do that [here](https://login.oslo.kommune.no/auth/realms/AD/account/totp).
* AWS account ID: This identifies which account you want to use. You can see which accounts you have access to just
after logging in to [AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws):

<span style="display:block;text-align:center">![okctl](img/aws-account-id.png)</span>

```bash

# Create a cluster. Format:
# okctl create cluster <environment name> <AWS account ID>
#
# <environment name>    can be for instance "prod" or "test".
# <AWS account ID>      is the account ID described in the above
#
# Example:
$ okctl create cluster prod 123456789012
```

Follow the instructions.

When done, verify that you have a working cluster by running

```bash
$ okctl venv
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
$ okctl create cluster --help

# Show credentials for cluster
$ okctl show credentials prod

# Run a sub shell with environment variables from the above command and a custom command prompt (PS1)
$ okctl venv

# Delete the cluster
$ okctl delete cluster prod
```

## Compare and contrast

With `okctl` we are attempting to solve the production environment setup problem. What we include within the definition of a production environment, we can see below.

| Functionality | okctl | [eksctl](https://eksctl.io) | [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) | [serverless.tf](https://serverless.tf) | 
|---|---|---|---|---|
|Provide authentication and authorisation|✓||?||
|Define a continuous integration pipeline|?||||
|Assist with creating a deployment pipeline|✓|✓|||
|Create a Kubernetes cluster (run-time environment)|✓|✓|||
|Facilitate with application setup and integration|✓||✓|✓|
|Provide integrations with version control system|✓||||
|Setup monitoring with metrics, logs and traces|✓||||

## Roadmap

We have a basic cluster up and running, and we are currently working on getting the core components of the cluster deployed:

- [x] [ExternalSecrets](https://github.com/godaddy/kubernetes-external-secrets/) for storing secrets securely
- [x] [AWS ALB Ingress Controller](https://github.com/kubernetes-sigs/aws-alb-ingress-controller) creates load balancers for incoming traffic
- [x] [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) ensures humane DNS hostnames 
- [ ] [Autoscaler](https://github.com/kubernetes/autoscaler/) for adjusting the size of pods and nodes
- [x] [Argo CD](https://github.com/argoproj/argo-cd) gives us continuous delivery
- [ ] [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator) for monitoring
- [ ] [Amazon Elastic Block Store (EBS) CSI driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/) for block storage
- [ ] [Loki](https://github.com/grafana/loki) for consuming logs
- [ ] [Jaeger](https://www.jaegertracing.io/) for traces

## Inspiration

We stand on the shoulders of giants, we have begged, borrowed and stolen from the following great products and teams:

- [Porter](https://github.com/deislabs/porter)
- [eksctl](https://github.com/weaveworks/eksctl)
- [kops](https://github.com/kubernetes/kops)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [go-kit](https://github.com/go-kit/kit)
- [saml2aws](https://github.com/Versent/saml2aws/)