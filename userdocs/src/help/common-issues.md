## Device authentication flow fails (Linux only)

If you are unable to complete device authentication against github, you need to install `pass`
This is because github token is stored in an encrypted keyring on your device.

Install pass as described in the [Prerequisites](../../getting-started/prerequisites/#pass-linux-only).

## okctl keeps trying to do the Github Device Authentication Flow while trying to do \<any action\>

This is known to happen if `pass init <gpg-key-id>` has not been run after installing `pass`.

Initialize `pass` as described in [Prerequisites](../../getting-started/prerequisites/#initialize-pass).

## okctl delete cluster fails

In some cases, it's possible to fix the issue and run `okctl delete cluster` again. We cannot give a guide general
guide for fixing this here, but ask us if you want support with this.

If not, you can proceed to delete the cluster manually, as described below.

### Delete application resources

* If you are able to use the cluster with kubectl, run

```shell
kubectl delete namespace monitoring
kubectl delete namespace argocd
```

Delete any namespace containing applications you have made

```shell
kubectl delete namespace myapp
```

* If you are not able to use the cluster with kubectl, open the AWS console
  * EC2 -> Load Balancers -> Delete any load balancers with tags matching the name of your cluster (try the filter
    `tag:elbv2.k8s.aws/cluster : mycluster-myenv`)
  * EC2 -> Auto Scaling Groups -> Delete any auto scaling groups with a name matching your cluster

### Delete Fargate Profile

Open the AWS console
* EKS -> Clusters -> Open your cluster -> Configuration -> Compute -> Delete all fargate profiles
  (usually one, `fp-default`) 

### Delete all CloudFormation stacks

Open the AWS console
* CloudFormation -> Enter a filter for your cluster name to only see the stacks for your cluster 
* Then for every stack repeat the following
  * Select the stack on the top (which is the newst)
  * Delete it
  * Refresh the page, wait until the stack has been deleted (it is gone)
    * (We have to wait because some stacks may have dependencies on other stacks) 

### Delete files in your infrastructure-as-code repository

* In your infrastructure-as-code repository, delete
  ** infrastructure/<env>
  ** .okctl.yaml, the block containing your environment (or the whole file if you only have one environment)


## okctl create cluster: Create identitypool fails / Re-create cluster within short timespan fails

If you do the following:

* Create a cluster
* Delete it
* Create a new cluster with the same domain name (e.g. whatever.oslo.systems)

This might fail if you do these steps within 15 minutes. This is due to DNS resolvers caching NS server records.
More details: https://github.com/oslokommune/okctl/pull/231

Workaround: Wait for up to 15 minutes before creating cluster again.

15 minutes is the TTL (Time to live, i.e. cache expiry) of the NS record. You can see this value in
Route 53 -> Hosted zones -> Your domain -> NS record for your top domain -> Edit -> See TTL field.

## okctl create cluster: Failed to create external secrets helm chart

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

