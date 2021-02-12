# `okctl` - Opinionated and effortless infrastructure and application management

`okctl` is a CLI for creating a production environment on [Amazon Web Services](https://aws.amazon.com/) (AWS) using [Elastic Kubernetes Service](https://aws.amazon.com/eks/) (EKS) as its foundation, we aim to be:

- **Specialised**: by tightly integrating with the services provided by Oslo kommune
- **Highly opinionated**: you will get something that works well that may not be exactly what you want
- **Modern**: base our decisions on best-practise methodologies for a modern development lifecycle
- **Sensible**: provide a set of sensible default integrations

<span style="display:block;text-align:center">![okctl](img/logo.png)</span>

## Installation

To download the latest release, make sure you have the required pre-requisites. Then run the command matching your
operating system.

### Prerequisites

* Your own GPG key
    * List existing keys with `gpg --list-keys`
    * If you don't have any keys, create a new one with `gpg --gen-key`

### Linux

```bash
curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/okctl /usr/local/bin
```

### macOS

#### Prerequisites

* `pass`
    * If you do not have it, run `brew install pass`
* `pass` must be initalized
    * Run `pass ls` to see if you have a password store. If it returns `Error: password store is empty. Try "pass init".`
      then run `pass init <GPG ID>` where GPG ID is the ID of your GPG key - use `gpg --list-keys` to see your keys (both
      e-mail and the ID are valid values, but if you have several keys for the same e-mail, use the ID).

#### Installation

```bash
brew tap oslokommune/tap
brew install oslokommune/tap/okctl
```

### Shell autocompletion

To get autocompletion in your shell working run `okctl completion -h` and follow the steps for your preferred shell.

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
git clone git@github.com:oslokommune/<the new repository>.git
cd <the new repository>
```

### 2. Create a cluster

A "cluster" is a Kubernetes cluster in an environment as described in [Compare and contrast](#compare-and-contrast).

#### Gather input data

You will soon be running `okctl create cluster`, which will ask you for the following information:

* Username and password: This is your Oslo Kommune AD organization username (oooXXXXX) and its password.
* Multi factor token (MFA): The same one you use to login to
[AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws). If you haven't set up MFA yet, you
can do that [here](https://login.oslo.kommune.no/auth/realms/AD/account/totp).
* AWS account ID: This identifies which account you want to use. You can see which accounts you have access to just
after logging in to [AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws):

<span style="display:block;text-align:center">![okctl](img/aws-account-id.png)</span>

#### Run create command

```bash
# Create a cluster. Format:
# okctl create cluster <environment name> <AWS account ID>
#
# <environment name>    can be for instance "prod" or "test".
# <AWS account ID>      is the account ID described in the above
#
# Example:
okctl create cluster prod 123456789012
```

Follow the instructions. If something fails, it's safe to just run the same command again (it's idempotent).

When done, verify that you have a working cluster by running

```bash
okctl venv prod
kubectl get service

```

The last command should show something like

```bash
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.100.0.1   <none>        443/TCP   1h
```

## Common commands

```bash
# Get help for any command
okctl --help
okctl create cluster --help

# Show credentials for cluster
okctl show credentials prod

# Run a sub shell with environment variables from the above command and a custom command prompt (PS1)
okctl venv prod

# Delete the cluster
okctl delete cluster prod
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

- [x] [ExternalSecrets](https://github.com/external-secrets/kubernetes-external-secrets/) for storing secrets securely
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

## Known issues, they are on our backlog :)

**Device authentication flow fails**

If you are unable to complete device authentication against github, you need to install `pass`
This is because github token is stored in an encrypted keyring on your device.

Linux
```bash
apt-get install pass
```

Mac
```bash
brew install pass
```

**okctl keeps trying to do the Github Device Authentication Flow while trying to do \<any action\>**

This is known to happen if `pass init <gpg-key-id>` has not been run after installing `pass`. 

**On `okctl delete cluster`, some resources are not deleted (automatic deletion is coming in a later version)**

Workaround: manually deleted the following resources:

* Deploy key in IAC repo
* Hosted zone and records (unless you specify "i-know-what-i-am-doing-delete-hosted-zone-and-records" flag to "true")
* It is recommended to delete the infrastructure/<env> directory and .okctl.yaml file upon successful delete of cluster, as the last manual step.

**ArgoCD doesn't show my apps**

This is due to an authorization bug.

Workaround:
```bash
kubectl edit configmap argocd-rbac-cm --namespace argocd
```

Add a new line after `g, admins, role:admin`, so it becomes something like this:
```
policy.csv: |
  g, admins, role:admin
  g, my.email@mail.com, role:admin
```

**okctl create cluster: Create identitypool fails / Re-create cluster within short timespan fails**

If you do the following:

* Create a cluster
* Delete it
* Create a new cluster with the same domain name (e.g. whatever.oslo.systems)

This might fail if you do these steps within 15 minutes. This is due to DNS resolvers caching NS server records.
More details: https://github.com/oslokommune/okctl/pull/231

Workaround: Wait for up to 15 minutes before creating cluster again.

15 minutes is the TTL (Time to live, i.e. cache expiry) of the NS record. You can see this value in
Route 53 -> Hosted zones -> Your domain -> NS record for your top domain -> Edit -> See TTL field.

**okctl create cluster: Failed to create external secrets helm chart**

You get the following error (shortened):

```
..  creating: external-secrets (elapsed: 1 second 76 microseconds)WARN[0007] failed to process request, because: failed to create external secrets helm chart: failed to update repository: failed to fetch https://kubernetes-charts-incubator.storage.googleapis.com/index.yaml : 403 Forbidden  endpoint=create service=helm/externalSecrets
✓   creating
Error:
....
request failed with Internal Server Error, because: failed to create external secrets helm chart: failed to update repository: failed to fetch https://kubernetes-charts-incubator.storage.googleapis.com/index.yaml : 403 Forbidden
```

This happens because Helm
[changed URLs](https://helm.sh/blog/new-location-stable-incubator-charts/#:~:text=The%20new%20location%20for%20the,use%20before%20November%2013%2C%202020.)
to their repositories. Update your ~/.okctl/helm/repositories.yaml, and update URLs from

| Name     | Old Location                                               | New Location                     |
| -------- | ---------------------------------------------------------- | -------------------------------- |
stable	   | https://kubernetes-charts.storage.googleapis.com           | https://charts.helm.sh/stable    |
incubator  | https://kubernetes-charts-incubator.storage.googleapis.com | https://charts.helm.sh/incubator |

**ArgoCD fails first run**

Workaround: re-run create command.

**Service quota check will check even if cluster is already running**

Workaround: If you already created a cluster, but need to re-run the command if for example ArgoCD failed. You will be warned that there are not enough resources. Continue anyway.
