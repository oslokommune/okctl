# Namespaces/

To remove ownership of a namespace from an individual application, we've set up this folder and made ArgoCD
automatically track changes to it. This folder contains all namespace manifests. Adding manifests to this folder will
automatically apply them to your cluster.

When running the apply application command, okctl will place a namespace manifest here.