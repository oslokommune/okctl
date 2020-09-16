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

### Using Argo CD

This demonstrates how you can use Argo CD in conjunction with some of the other components we have installed in the base cluster. In your git repository, save the following definitions under a directory, e.g, `{git_root}/deployment/2048/`.

**NB:** Remember to change the values in brackets `{...}` so they match your cluster setup.

`{git_root}/deployment/2048-game.yml`:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: "2048-game"
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: "2048-ingress"
  namespace: "2048-game"
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
  labels:
    app: 2048-ingress
spec:
  rules:
    - host: 2048-game.{team/product}.oslo.systems
      http:
        paths:
          - path: /*
            backend:
              serviceName: "service-2048"
              servicePort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: "service-2048"
  namespace: "2048-game"
spec:
  ports:
    - port: 80
      targetPort: 80
      protocol: TCP
  type: NodePort
  selector:
    app: "2048"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: "2048-deployment"
  namespace: "2048-game"
spec:
  selector:
    matchLabels:
      app: "2048"
  replicas: 5
  template:
    metadata:
      labels:
        app: "2048"
    spec:
      containers:
      - image: alexwhen/docker-2048
        imagePullPolicy: Always
        name: "2048"
        ports:
        - containerPort: 80
```

Then create an Argo CD application definition:

`{git_root}/deployment/2048-application.yaml`:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: "2048-application"
spec:
  destination:
    namespace: "2048-game"
    server: 'https://kubernetes.default.svc'
  source:
    path: deployment/2048/
    repoURL: 'git@github.com:{organisation}/{repository}.git'
    targetRevision: HEAD
  project: default
  syncPolicy:
    automated:
      prune: false
      selfHeal: false
```

Commit these definitions to your git repository and push them to master. Then you can apply the application definition with `kubectl apply -f deployment/2048-application.yaml`.  If everything is configured correctly, you should now see the "2048-application" in the Argo CD UI. It should also give you a status on the deployment. Argo CD will ensure that the definition stored under the provided path is applied to your cluster. Eventually you should be able to reach your 2048 game from the host you defined: `http://2048-game.{team/product}.oslo.systems`.