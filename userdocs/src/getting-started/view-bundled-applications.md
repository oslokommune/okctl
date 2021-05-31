Okctl, if enabled, will provide several bundled applications. Here you can learn more about them.

## Access the bundled applications

Some of our bundled applications require you to log in to be able to use. To create an account which can be used to log
in to these applications, declare the users in your cluster declaration:

```bash
# cluster.yaml
...

users:
- email: user.email@emailprovider.org

...
```

This creates a user in `Cognito` which can be used to log into the following bundled applications:

* ArgoCD UI
* Grafana

## ArgoCD

The bundled application ArgoCD has a web GUI where you can administrate and observe your deployed applications. ArgoCD
is by default made available at 

`argocd.<cluster root domain>`

For example, if the cluster in question has the name `citykey-test`, ArgoCD will be made available at

`argocd.citykey-test.oslo.systems`

:information_source: More information on ArgoCD can be found [here](/buildingblocks/argocd/).

## Prometheus ❤ Grafana

The bundled application Prometheus has a web GUI where you can administrate and observe metrics regarding your cluster
and applications. This web GUI is called Grafana and is by default made available at

`grafana.<cluster root domain>`

For example, if the cluster in question has the name `dataplatform-test`, ArgoCD will be made available at

`grafana.dataplatform-test.oslo.systems`
