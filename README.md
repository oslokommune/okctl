[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](code_of_conduct.md) [![codecov](https://codecov.io/gh/oslokommune/okctl/branch/master/graph/badge.svg)](https://codecov.io/gh/oslokommune/okctl) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=oslokommune_okctl&metric=security_rating)](https://sonarcloud.io/dashboard?id=oslokommune_okctl) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=oslokommune_okctl&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=oslokommune_okctl)

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
okctl --help
```

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