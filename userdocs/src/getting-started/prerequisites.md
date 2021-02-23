The following is a guide for how to create an environment that contains the elements described in
[Compare and contrast](#compare-and-contrast).

## 1. Satisfy requirements

### GPG key

Make sure you have your own GPG key. To list existing keys, run `gpg --list-keys`

If you don't have any keys, create a new one with `gpg --gen-key`

### Pass (Linux only)

#### Install `pass`

Follow the instructions [here](https://www.passwordstore.org/#download) to install `pass`

#### Initialize `pass`

Run `pass ls` to see if you have a password store. If it returns `Error: password store is empty. Try "pass init".`
then run `pass init <GPG ID>` where GPG ID is the ID of your GPG key - use `gpg --list-keys` to see your keys (both
e-mail and the ID are valid values, but if you have several keys for the same e-mail, use the ID).

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

:information_source: A "cluster" is a Kubernetes cluster in an environment as described in [Compare and contrast](#compare-and-contrast).

You will soon be creating a cluster, which will ask you for the following information:

* Username and password: This is your Oslo Kommune AD organization username (oooXXXXX) and its password.
* Multi factor token (MFA): The same one you use to login to
[AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws). If you haven't set up MFA yet, you
can do that [here](https://login.oslo.kommune.no/auth/realms/AD/account/totp).
* AWS account ID: This identifies which account you want to use. You can see which accounts you have access to just
after logging in to [AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws):

<span style="display:block;text-align:center">![okctl](../img/aws-account-id.png)</span>

## 4. Run create command

```bash
# Create a cluster. Format:
# okctl create cluster <environment name> <AWS account ID>
#
# <environment name>    can be for instance "prod" or "test".
# <AWS account ID>      is the account ID described in the above
#
# Example:
okctl create cluster prod 123456789012
```

Follow the instructions. If something fails, it's safe to just run the same command again (it's idempotent).

When done, verify that you have a working cluster by running

```bash
okctl venv prod
kubectl get service

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
okctl create cluster --help

# Show credentials for cluster
okctl show credentials prod

# Run a sub shell with environment variables from the above command and a custom command prompt (PS1)
okctl venv prod

# Delete the cluster
okctl delete cluster prod
```
