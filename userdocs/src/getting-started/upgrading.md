After installing your cluster, the Okctl team will eventually update the `okctl` tool in order to bring new features and
bug fixes to the infrastructure and applications that `okctl` sets up. Downloading the new version of `okctl` ensures
that _new_ clusters and resources will be up-to-date and working properly.

However, _existing_ clusters and their resources need to be upgraded, which is what `okctl upgrade` does. By resource,
we mean everything that is managed by okctl, including applications.

For instance, a new version of `okctl` installs a newer version of ArgoCD when applying a cluster from scratch. For
existing clusters, however, since ArgoCD already exists, running `okctl apply cluster` won't change anything. To upgrade
ArgoCD and other resources that are lagging behind, you can run `okctl upgrade`.

## Upgrade

```bash
# Usage
okctl --cluster-declaration <path to cluster declaration> upgrade

# Example:
okctl --cluster-declaration cluster.yaml upgrade
```

You will be prompted before continuing.

Example output:

```shell
$ okctl2 -c cluster-dev.yaml upgrade
Found 2 applicable upgrade(s):
0.0.50, 0.0.51

Simulating upgrades...

--- Simulating upgrade: okctl-upgrade_0.0.50 ---
Simulating the upgrade, not doing any changes.
Upgrading ArgoCD
Upgrading ArgoCD complete!
--- Simulating upgrade: okctl-upgrade_0.0.51 ---
Simulating the upgrade, not doing any changes.
Upgrading Some Other Component
Upgrading Some Other Component complete!

Simulating upgrades complete.

? This will upgrade your okctl cluster, are you sure you want to continue? Yes

--- Running upgrade: okctl-upgrade_0.0.66 ---
Upgrading ArgoCD
Upgrading ArgoCD complete!
--- Running upgrade: okctl-upgrade_0.0.67 ---
Upgrading Some Other Component
Upgrading Some Other Component complete!

Upgrade complete! Cluster version is now 0.0.70. Remember to commit and push changes with git.

```

## Questions and answers

**Is there any way of knowing more in detail what okctl upgrade will do?**

Yes, go to https://github.com/oslokommune/okctl-upgrade/releases and read about each upgrade.

**Can I downgrade after upgrading?**

No. If you want to know that everything works as expected after upgrading, we advise you to create an additional cluster
for testing. This test cluster should be as similar as possible to your production cluster, running the same
applications and with the same configuration. Run `okctl upgrade` on this cluster to see that everything works as
expected after upgrading.
