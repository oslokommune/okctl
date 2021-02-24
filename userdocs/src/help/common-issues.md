**Device authentication flow fails (Linux only)**

If you are unable to complete device authentication against github, you need to install `pass`
This is because github token is stored in an encrypted keyring on your device.

Install pass as described in the [Prerequisites](../../getting-started/prerequisites/#pass-linux-only).

**okctl keeps trying to do the Github Device Authentication Flow while trying to do \<any action\>**

This is known to happen if `pass init <gpg-key-id>` has not been run after installing `pass`.

Initialize `pass` as described in [Prerequisites](../../getting-started/prerequisites/#initialize-pass).

**On `okctl delete cluster`, some resources are not deleted (automatic deletion is coming in a later version)**

Workaround: manually delete the following resources:

* Hosted zone and records (unless you specify "i-know-what-i-am-doing-delete-hosted-zone-and-records" flag to "true")
* It is recommended to delete the infrastructure/<env> directory and .okctl.yaml file upon successful delete of cluster, as the last manual step.

**ArgoCD doesn't show my apps**

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

**okctl create cluster: Create identitypool fails / Re-create cluster within short timespan fails**

If you do the following:

* Create a cluster
* Delete it
* Create a new cluster with the same domain name (e.g. whatever.oslo.systems)

This might fail if you do these steps within 15 minutes. This is due to DNS resolvers caching NS server records.
More details: https://github.com/oslokommune/okctl/pull/231

Workaround: Wait for up to 15 minutes before creating cluster again.

15 minutes is the TTL (Time to live, i.e. cache expiry) of the NS record. You can see this value in
Route 53 -> Hosted zones -> Your domain -> NS record for your top domain -> Edit -> See TTL field.

**okctl create cluster: Failed to create external secrets helm chart**

You get the following error (shortened):

```
..  creating: external-secrets (elapsed: 1 second 76 microseconds)WARN[0007] failed to process request, because: failed to create external secrets helm chart: failed to update repository: failed to fetch https://kubernetes-charts-incubator.storage.googleapis.com/index.yaml : 403 Forbidden  endpoint=create service=helm/externalSecrets
✓   creating
Error:
....
request failed with Internal Server Error, because: failed to create external secrets helm chart: failed to update repository: failed to fetch https://kubernetes-charts-incubator.storage.googleapis.com/index.yaml : 403 Forbidden
```

This happens because Helm
[changed URLs](https://helm.sh/blog/new-location-stable-incubator-charts/#:~:text=The%20new%20location%20for%20the,use%20before%20November%2013%2C%202020.)
to their repositories. Update your ~/.okctl/helm/repositories.yaml, and update URLs from

| Name     | Old Location                                               | New Location                     |
| -------- | ---------------------------------------------------------- | -------------------------------- |
stable	   | https://kubernetes-charts.storage.googleapis.com           | https://charts.helm.sh/stable    |
incubator  | https://kubernetes-charts-incubator.storage.googleapis.com | https://charts.helm.sh/incubator |

**ArgoCD fails first run**

Workaround: re-run create command.

**Service quota check will check even if cluster is already running**

Workaround: If you already created a cluster, but need to re-run the command if for example ArgoCD failed. You will be warned that there are not enough resources. Continue anyway.
