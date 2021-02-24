# Outdated common issues
Here you will find bugs that are fixed in newer versions, but that still might cause problems for clusters created with 
older versions of okctl.

##ArgoCD doesn't show my apps

This is due to an authorization bug.

Workaround:
```bash
kubectl edit configmap argocd-rbac-cm --namespace argocd
```

Add a new line after `g, admins, role:admin`, so it becomes something like this:
```
policy.csv: |
  g, admins, role:admin
  g, my.email@mail.com, role:admin
```

##ArgoCD fails first run

Workaround: re-run create command.

##Service quota check will check even if cluster is already running

Workaround: If you already created a cluster, but need to re-run the command if for example ArgoCD failed. You will be warned that there are not enough resources. Continue anyway.
