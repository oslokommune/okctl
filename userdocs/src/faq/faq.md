
## How do I grant access to a cluster to a team member?

### Prerequisites

Your team member needs the following privileges

1. Access to the AWS account where the cluster runs
1. They need access to the IAC repo on Github that you run `okctl` in

### Run `kubectl` commands

Your team member can now

1. Clone and enter the IAC repository
1. Run `okctl venv <environment name>`

Your team member should now be able to use `kubectl` commands

## How do i store my docker-images in aws or github?

## Which sub-domain names are set up?

We reserve some of the subdomain of the primary hosted zone, e.g., "argocd.<primary>" and "grafana.<primary>"

## How can I use the ESB CSI driver?

## How can I create service monitors and dashboards in prometheus?