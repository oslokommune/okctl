
## Motivation

Something catchy

## Commands

To create an application.yaml template, run the following command:

`okctl create application > application.yaml`

This creates a minimal application.yaml declaration in ./application.yaml. For a more extensive template, add the `-f/--full` flag.

After configuring the application.yaml file, you turn it into Kubernetes and ArgoCD resources by running:

`okctl apply application -f application.yaml`

This command will create the following two files in the ./deployment folder:
1. `<app-name>.yaml` containing all of the Kubernetes resources
2. `<app-name>-application.yaml` containing the ArgoCD Application declaration.

Both is needed by ArgoCD to deploy your application or service. Read more about ArgoCD [here](https://okctl.io/deployment/argocd/).
