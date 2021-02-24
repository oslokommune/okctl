> **NB:** While this describes what Argo CD is and how it can be used, the intent is to faciliate our users in the creation of these resources.

## What value does Argo CD provide

With `okctl` we use [Argo CD](https://argoproj.github.io/argo-cd/) for managing application deployment. Argo CD is based on a [gitops](https://www.weave.works/technologies/gitops/) model:

> GitOps is a way to do Kubernetes cluster management and application delivery.  It works by using Git as a single source of truth for declarative infrastructure and applications. With GitOps, the use of software agents can alert on any divergence between Git with what's running in a cluster, and if there's a difference, Kubernetes reconcilers automatically update or rollback the cluster depending on the case. With Git at the center of your delivery pipelines, developers use familiar tools to make pull requests to accelerate and simplify both application deployments and operations tasks to Kubernetes.

The primary advantages of using a gitops model are:

1. The entire system described declaratively.
2. The canonical desired system state versioned in Git.
3. Approved changes that can be automatically applied to the system.  
4. Software agents to ensure correctness and alert on divergence.

### Setup of Argo CD

When you create a full cluster we configure Argo CD for you, this configuration consists of a couple of operations:
 
1. Add a Github [deploy key](https://docs.github.com/en/developers/overview/managing-deploy-keys) to your Github repository with `read-only` access. This means that Argo CD can poll for changes to your git repository and see if it needs to deploy a new version or not.

2. Configure an ingress resource and make the Argo CD UI available on the following URL: `https://argocd.{product/team}.oslo.systems`.

3. Setup a Github [Oauth app](https://docs.github.com/en/developers/apps/building-oauth-apps) to manage access to the Argo CD UI. Only the configured Github organisation and team that are set during the configuration are given access to this UI.
