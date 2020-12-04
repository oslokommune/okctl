
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

## Commands

To create an application.yaml template, run the following command:

```bash
okctl create application ENV
```

This creates an application declaration in ./application.yaml.

After configuring the application.yaml file, you turn it into Kubernetes and ArgoCD resources by running:

```bash
okctl apply application prod -f application.yaml
```

This command will create the following files in the ./infrastructure folder:

1. ./infrastructure/base/applications/<app-name>
    * `<app-name>.yaml` containing all the Kubernetes resources.
    * `<app-name>-application.yaml` containing the ArgoCD Application declaration.
2. ./infrastructure/<env>/certificates/<application.url>
    * The certificate declaration for the URL specified in the application.yaml.

Both files in 1. is needed by ArgoCD to deploy your application or service. Read more about ArgoCD
[here](https://okctl.io/deployment/argocd/).

After that, the following manual steps remain:

1. (optional) Create the namespace you specified in the application declaration(application.yaml), i.e.:
`kubectl create namespace <name of namespace>`. This is only needed if the namespace you specified in the application
declaration is not pre-existing.
2. Commit and push the changes done by `okctl apply` to your infrastructure as code repository remote accessible for
ArgoCD.
3. Apply the ArgoCD resource to the cluster: `kubectl apply -f ./infrastructure/base/applications/<app-name>/<app-name>-application.yaml`

The application.yaml declaration can be neat to create resource files from for a while, but as soon as you find yourself
adding configuration directly to the Kubernetes resource file or the ArgoCD resource file, there is no longer any point
to keep it around. Feel free to delete it.  
