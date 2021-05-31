
## Motivation

Deploying an application in Kubernetes can be challenging, especially if one is new to the declarative paradigm. More so 
when dealing with the concept of infrastructure as code through a technology like ArgoCD for the first time. Add a new 
kind of Ingress controller combined with the AWS platform for provisioning certificates, and you have yourself a jolly
good time in research land. 

Our deployment process aims to fix that. 

1. First you create and configure an application declaration. This declaration can be deleted after applying it. You'll
be making changes to the Kubernetes and ArgoCD resources from that point on.
2. We apply the application declaration to the infrastructure as code repository, which converts the application
declaration to resources understood by Kubernetes and ArgoCD.

## Deploy an application

:information_source: If you need information on how to make your Docker images available in the cluster, see [Running a Docker image in your cluster](/help/docker-registry)

### Scaffold the application declaration
To scaffold an application.yaml template, run the following command:

```bash
# Syntax
okctl --cluster-declaration=<relevant cluster.yaml> scaffold application

# Example
okctl --cluster-declaration=cluster.yaml scaffold application > application.yaml
```

This creates an application declaration in application.yaml.

After configuring the application.yaml file, you turn it into Kubernetes and ArgoCD resources by running:

### Apply the application declaration

```bash
# Syntax
okctl --cluster-declaration=<relevant cluster.yaml> apply application -f <relevant application.yaml>

# Example
okctl --cluster-declaration=cluster.yaml apply application -f application.yaml
```

This command will create the following content in the infrastructure folder:

1. `infrastructure/applications/APP-NAME/base`
    * Contains Kubernetes resources common for your application regardless of environment
2. `infrastructure/applications/APP-NAME/overlays`
    * Contains (Kustomize) patches containing environment specific adjustments to the common resources.

After that, the following manual steps remain:

1. (optional) Create the namespace you specified in the application declaration(application.yaml), i.e.:

    ```bash
    kubectl create namespace <name of namespace>
    ```

    This is only needed if the namespace you specified in the application declaration is not pre-existing.

2. Commit and push the changes done by `okctl apply application` to your infrastructure as code repository remote making
   them accessible for ArgoCD.

3. Apply the ArgoCD resource to the cluster:

    ```
    kubectl apply -f infrastructure/applications/<app-name>/overlays/<cluster name>/argocd-application.yaml
    ```
   
The application.yaml declaration is used as an alternative to a wizard or numerous command flags. After running 

```shell
okctl apply application
```

there is no longer any point to keep it around. Feel free to delete it.

## Setup an ArgoCD application manually

This section describes what happens behind the scenes in the process described above. You can skip this entirely, it's
just for educational purposes.

This demonstrates how you can use Argo CD in conjunction with some other components we have installed in the base cluster. In your git repository, save the following definitions under a directory, e.g, `{git_root}/deployment/2048/`.

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



