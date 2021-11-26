The Okctl developers is continously updating the `okctl` tool in order to bring new features and bug fixes to the infrastructure
and applications that `okctl` sets up. When you create a cluster and its components with okctl, the cluster is pinned to the
version of `okctl` you are using. This means that after a new version of `okctl` is released, _existing_ clusters should be
upgraded to stay up to date, which is what `okctl upgrade` does.

For instance, a new version of `okctl` installs a newer version of ArgoCD when applying a cluster from scratch. For existing
clusters, however, since ArgoCD already exists, running `okctl apply cluster` won't change anything. To upgrade ArgoCD and other
resources that are lagging behind, you can run `okctl upgrade`.

## Concepts

When you create a cluster with `okctl` version 0.0.50, the cluster version is set to 0.0.50. To upgrade, you download a newer
version of okctl, let's say 0.0.60, and run `okctl upgrade`. When this command completes, your cluster version will be 0.0.60.

Important commands such as `okctl apply cluster` and `okctl apply application` requires the okctl binary and the cluster version
to be the same in order to work. This means that after an upgrade, everybody that uses this cluster, must get that version of
okctl to run those commands. Some commands are a bit less restrictive, and only require the binary version to be equal or higher
than the cluster version. However, we recommend using the same version of the okctl binary as the cluster version.

## Upgrading

First, download a new version of okctl - see [Install](/getting-started/install).

Then run an upgrade command.

```bash
# Usage
okctl --cluster-declaration <path to cluster declaration> upgrade

# Example:
okctl --cluster-declaration cluster.yaml upgrade
```

You will be prompted before continuing.

Example output:

```shell
$ okctl -c cluster-dev.yaml upgrade
Found 2 applicable upgrade(s):
0.0.50.argocd, 0.0.51.some-other-component

preloading missing binary: okctl-upgrade_0.0.50.argocd (0.0.50.argocd)
preloading missing binary: okctl-upgrade_0.0.51.some-other-component (0.0.51.some-other-component)

Simulating upgrades (we're not doing any actual changes yet, just printing what's going to happen)...

--- Simulating upgrade: okctl-upgrade_0.0.50.argocd ---
Upgrading ArgoCD
Upgrading ArgoCD complete!
--- Simulating upgrade: okctl-upgrade_0.0.51.some-other-component ---
Upgrading Some Other Component
Upgrading Some Other Component complete!

Simulating upgrades complete.

? This will upgrade your okctl cluster, are you sure you want to continue? Yes

--- Running upgrade: okctl-upgrade_0.0.50.argocd ---
Upgrading ArgoCD
Upgrading ArgoCD complete!
--- Running upgrade: okctl-upgrade_0.0.51.some-other-component ---
Upgrading Some Other Component
Upgrading Some Other Component complete!

Upgrade complete! Cluster version is now 0.0.70. Remember to commit and push changes with git.
```

## Questions and answers

**Is there any way of knowing more in detail what okctl upgrade will do?**

Yes, go to the [upgrade release page](https://github.com/oslokommune/okctl-upgrade/releases) and read about each upgrade.

Also, as you can see from the example output above, upgrades are being simulated before actually applied.

**Can I downgrade after upgrading?**

No. If you want to know that everything works as expected after upgrading, we advise you to create an additional cluster for
testing. This test cluster should be as similar as possible to your production cluster, running the same applications and with the
same configuration. Run `okctl upgrade` on this cluster to see that everything works as expected after upgrading.

**Okctl says something about initializing parts of the cluster state. What does it mean?**

When we released the upgrade feature, parts of the state (the state.db file) in okctl needed to be initialized before the user
could run the `upgrade` command. Okctl prints this information because user action is needed, that is, commit the state.db file
and push it.

Specficially, the cluster initialization means setting up data structures in the state for storing cluster versions, which is
needed by upgrade to calculate which upgrades to apply.
