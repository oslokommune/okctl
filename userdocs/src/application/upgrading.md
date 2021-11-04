Sometimes newer versions of okctl updates the way `okctl apply application` works. For example, the latest okctl version
may create an application with better health checks than the old version.

To apply these changes, you will need to upgrade the application. This guide describes how to do this.

Some motivations to upgrade an application include

* the newest okctl has fixed a bug from a previous version
* you want to stay updated with how okctl works
* Kubernetes manifests are being deprecated (for instance, newer versions of Kubernetes require Ingress versions to
  be `apiVersion: networking.k8s.io/v1` instead of `apiVersion: networking.k8s.io/v1beta1`)

# Method

`okctl apply application` writes Kubernetes manifests from scratch. It does not care about existing Kubernetes
manifests, it simply overwrites any existing files. Since you probably have made some changes to these files, simply
running `okctl apply application` again would overwrite these changes.

In order to update the manifests, you need to do a manual merge. You can do this in any way you want, but in the steps
below, we describe how to do this with git. The principle is to keep your changes, but add the stuff that is new from
okctl.

# Steps

Let's say you have previously scaffolded and applied an application like this:

```shell
okctl scaffold application > my-app.yaml
okctl -c my-cluster.yaml apply application -f my-app.yaml
```

:tip: If you don't have the application manifest (`my-app-yaml`), you can attempt to re-create it by
running `okctl scaffold application` and edit the values to match your existing application.

To upgrade, start by downloading the latest version of okctl. See [install instructions](/getting-started/install.md).

Ensure you have a clean git state:

```shell
git status -s # should return no output
```

Then delete existing application manifests. This ensures that we delete files that okctl apply application doesn't
produce anymore. You _can_ skip this command to get a simpler git diff below, but you should only do so if you know what
you are doing (i.e. you know that `okctl apply application` doesn't produce fewer files than it did last time you
ran `okctl apply application`).

```shell
rm -rf infrastructure/application/my-app
```

Re-create application manifests from scratch.

```shell
okctl -c my-cluster.yaml apply application -f my-app.yaml
```

Now, use your favourite git diff tool to get back what is specific for your application.

In the example screenshot below, we're using Intellij to do the diff, and the newest version of okctl has
changed `initialDelaySeconds` from `3` to `5`. So, you would apply those two changes. Additionaly, my application uses
an environment variable `DSN`, so you would apply that change as well.

![okctl](/img/application-upgrade-diff.png)

When you have reapplied all changes, all that is left is committing and pushing the changes.

```shell
git add .
git commit -m "Upgrade my-app to okctl 0.0.60"
git push
```

You have now successfully upgraded your application to the latest version of okctl.
