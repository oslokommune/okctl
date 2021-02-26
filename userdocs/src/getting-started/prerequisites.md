The following is a guide for how to create an environment that contains the elements described in
[Compare and contrast](/#compare-and-contrast).

## 1. Ensure the following requirements

### GPG key

Make sure you have your own GPG key. To list existing keys, run `gpg --list-keys`

If you don't have any keys, create a new one with `gpg --gen-key`

## 2. Create a new github repository

Go to github.com/oslokommune and create a new **private** git repository.

:information_source: This repository will be used by okctl to store infrastructure-as-code, which means files containing
various configuration for your up-and-coming cluster.

Now, run:

```bash

# Clone the repository you just made
git clone git@github.com:oslokommune/<the new repository>.git
cd <the new repository>
```

## 3. Prepare information for cluster creation

:information_source: A "cluster" is a Kubernetes cluster in an environment as described in [Compare and contrast](/#compare-and-contrast).

You will soon be creating a cluster, which will ask you for the following information:

* Username and password: This is your Oslo Kommune AD organization username (oooXXXXX) and its password.
* Multi factor token (MFA): The same one you use to login to
[AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws). If you haven't set up MFA yet, you
can do that [here](https://login.oslo.kommune.no/auth/realms/AD/account/totp).
* AWS account ID: This identifies which account you want to use. You can see which accounts you have access to just
after logging in to [AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws):

<span style="display:block;text-align:center">![okctl](../img/aws-account-id.png)</span>

## 4. Create a cluster

### Run scaffold cluster command

```bash
# Scaffold a cluster. Format:
# okctl scaffold cluster <cluster name> <environment name>
#
# <cluster name>        should be something representing your team or your product
# <environment name>    can be for instance "prod" or "test".
#
# Example:
okctl scaffold cluster kjoremiljo production > cluster.yaml
```

### Run apply cluster command

```bash
# Apply a cluster declaration. Format:
# okctl apply cluster -f <cluster declaration path>
#
# <cluster name>        should be something representing your team or your product
#
# Example:
okctl apply cluster -f cluster.yaml
```

When done, verify that you have a working cluster by running

```bash
okctl venv production
kubectl get services
```

The last command should show something like

```bash
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.100.0.1   <none>        443/TCP   1h
```

## Common commands

```bash
# Get help for any command
okctl --help
okctl scaffold cluster --help

# Scaffold a cluster declaration template
okctl scaffold cluster
# Apply a cluster declaration file
okctl apply cluster

# Scaffold an application declaration template
okctl scaffold application
# Apply an application declaration file
okctl apply application

# Show credentials for cluster
okctl show credentials <environment name>

# Run a sub shell with environment variables from the above command and a custom command prompt (PS1)
okctl venv <environment name>

# Delete the cluster
okctl delete cluster <environment name
```
