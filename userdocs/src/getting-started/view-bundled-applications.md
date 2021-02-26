
## Access the bundled applications

Some of our bundled applications require you to log in to be able to use. To create an account which can be used to log
in to these applications, use the following command:

```bash
# Usage
# okctl adduser ENV EMAIL
#
# ENV     Name of the environment to create the user in
# EMAIL   The email of the user to add
#
# Example
okctl adduser production jane.doe@origo.oslo.kommune.no
```

This creates a user in `Cognito` which can be used to log into the following bundled applications:

* ArgoCD UI
* Grafana

### ArgoCD

The bundled application ArgoCD has a web GUI where you can administrate and observe your deployed applications. ArgoCD
is by default made available at 

`argocd.<cluster primary hosted zone>`

For example, if the cluster in question has the name `citykey` and the environment `test`, ArgoCD will be made available
at

`argocd.citykey-test.oslo.systems`

:information_source: More information on ArgoCD can be found [here](/buildingblocks/argocd/).

### Prometheus ❤ Grafana

The bundled application Prometheus has a web GUI where you can administrate and observe metrics regarding your cluster
and applications. This web GUI is called Grafana and is by default made available at

`grafana.<cluster primary hosted zone>`

For example, if the cluster in question has the name `dataplatform` and the environment `test`, ArgoCD will be made available
at

`grafana.dataplatform-test.oslo.systems`
