The following is a guide for how to create an environment that contains the elements described in
[Compare and contrast](#compare-and-contrast).

## 1. Create a new github repository

Go to github.com/oslokommune and create a new **private** git repository.

:information_source: This repository will be used by okctl to store infrastructure-as-code, which means files containing
various configuration for your up-and-coming cluster.

Now, run:

```bash

# Clone the repository you just made
git clone git@github.com:oslokommune/<the new repository>.git
cd <the new repository>
```

## 2. Create a cluster

A "cluster" is a Kubernetes cluster in an environment as described in [Compare and contrast](#compare-and-contrast).

### Gather input data

You will soon be creating a cluster, which will ask you for the following information:

* Username and password: This is your Oslo Kommune AD organization username (oooXXXXX) and its password.
* Multi factor token (MFA): The same one you use to login to
[AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws). If you haven't set up MFA yet, you
can do that [here](https://login.oslo.kommune.no/auth/realms/AD/account/totp).
* AWS account ID: This identifies which account you want to use. You can see which accounts you have access to just
after logging in to [AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws):

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
