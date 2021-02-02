Some people enjoy a lengthy CLI wizard. Others prefer to declare their resources and let the CLI tool do its thing.

Introducing declarative cluster.

To create a cluster using the declarative approach:

1. Create a declaration and edit it to suit your needs
2. Apply the declaration


## Commands

To scaffold a cluster declaration, run the following:

```bash
# Usage
okctl scaffold cluster TEAM-/PRODUCT_NAME ENVIRONMENT_NAME > path

# Example
okctl scaffold cluster utviklerportalen production > cluster.yaml
```

To create a cluster based on the declaration, run the following:

```bash
# Usage
okctl apply cluster -f PATH

# Example
okctl apply cluster -f cluster.yaml
```

## Details

The declarative approach does exactly the same as `okctl create cluster`, except that it gets all its user input up
front. That way you can easily take down and recreate clusters as you experiment. You create a declaration once and 
apply it as many times as you need.

## Usage

When running 

```bash
okctl scaffold cluster utviklerportalen production > cluster.yaml
``` 

you'll end up with a file looking like this:

```yaml
apiVersion: okctl.io/v1alpha1
kind: Cluster

metadata:
  # Account ID is your AWS account ID.
  accountID: "123456789123"
  # Environment is the name you use to identify the type of cluster it is. Common names are production, test, staging
  environment: production
  # Name can be anything, but should define the scope of the cluster. Meaning if the cluster is scoped to one product,
  # you might want to name it the name of the product. If the cluster contains all services and products owned by a
  # team, the team name might be more fitting.
  name: utviklerportalen
  # Region defines the AWS region to prefer when creating resources
  region: eu-west-1

# The primary DNS zone defines the domain of which to create services beneath. For example; okctl will setup ArgoCD
# which has a frontend. The frontend will be available at argocd.<parentDomain>. For Cognito it will be 
# auth.<parentDomain>
primaryDNSZone:
  managedZone: false
  parentDomain: utviklerportalen-production.oslo.systems

# For okctl to be able to setup ArgoCD correctly for you, it needs to know what repository on Github that will contain
# your infrastructure
github:
  # The organization that owns the repository
  organisation: oslokommune
  # The name of the repository
  repository: my_iac_repo_name
  # The folder to place (and look for) application declarations
  outputPath: infrastructure
  # The team defined in Github that should be allowed access to the created resources
  team: my_team

# okctl creates a Virtual Private Cloud for you which it organizes all the intended resources that require networking.
vpc:
  # CIDR defines the VPC IP range. Leave this be if you don't know what it is/does
  cidr: 192.168.0.0/20
  highAvailability: true

integrations:
  # ALB Ingress Controller handles routing from the internet to your application running inside your okctl Kubernetes
  # cluster. If you want your applications and services accessible from the internet, this needs to be enabled
  albIngressController: true
  # ArgoCD is a service that watches a repository for Kubernetes charts and ensures the defined resources are running
  # as declared in the cluster
  argoCD: true
  # Cognito is an authentication provider that okctl uses to control access to different resources, like ArgoCD
  cognito: true
  # External DNS handles defining the necessary DNS records required to route traffic to your defined service or 
  # application
  externalDNS: true
  # External Secrets handles 
  externalSecrets: true
```

Modify the declaration to suit your situation and needs, then use

```bash
okctl apply cluster -f cluster.yaml
```

to have okctl generate a cluster based on the declaration.

That's it. Sit back and enjoy or go do something else while `okctl` does its thing.
