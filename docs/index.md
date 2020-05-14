---
layout: default
---

# Installation

To download the latest release, run the command matching your operating system:

```bash
# Linux
curl --silent --location "https://github.com/oslokommune/okctl/releases/latest/download/okctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/okctl /usr/local/bin

# macOS
brew tap oslokommune/tap
brew install oslokommune/tap/okctl
```

# Usage

```bash
okctl --help
```

# Motivation

TL;DR a good, but not great, solution for running production code in Oslo kommune

okctl is a CLI for creating a production grade run-time environment on AWS using EKS (Kubernetes). It is highly opinionated, which means that you will get something that works well, but not necessarily exactly what you want. The idea here is to provide a set of sensible defaults that allow for modern development lifecycles. We provide defaults for:

- Continuous integration
  - Github actions
- Continuous deployment
  - Github actions
  - Github packages / AWS ECR (AWS ECR has the advantage that it supports security scans on the containers)
  - ArgoCD
- Monitoring
  - Prometheus for metrics
  - Loki for logs
  - Grafana for dashboards
- Storage
  - AWS RDS postgres

We are not attempting to write something that works for the whole world, we are entirely focused on Oslo kommune, so we are building a custom CLI only for ourselves.

We also believe that functioning as a good ORM makes a lot of sense, e.g., most of the time you should be able to achieve what you want by using the provided interfaces. However, sometimes you simply need to break out of the managed aspect and run a raw query directly towards the database.

## Principals

- User first, user first, user first.
- We are building something for other developers, not for ourselves
- Meet the users where they are

- Try to reduce the number of binary dependencies as much as possible
- Try to reduce the size of the okctl binary itself as much as possible
- Be a completely self-contained binary that is easy to distribute
- The CLI should be self-explanatory and provide all the documentation required
