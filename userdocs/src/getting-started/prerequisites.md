The following prerequisites are required to initiate and complete the cluster creation process.

## GPG key

Make sure you have your own GPG key. To list existing keys, run `gpg --list-keys`

If you don't have any keys, create a new one with `gpg --gen-key`

## Create a new GitHub repository

Go to github.com/oslokommune and create a new **private** git repository.

:information_source: This repository will be used by okctl to store infrastructure-as-code, which means files containing
various configuration for your up-and-coming cluster.

Now, run:

```bash

# Clone the repository you just made
git clone git@github.com:oslokommune/<the new repository>.git
cd <the new repository>
```

## Prepare information for cluster creation

:information_source: A "cluster" is a Kubernetes cluster in an environment as described in [Compare and contrast](/#compare-and-contrast).

You will soon be creating a cluster, which will ask you for the following information:

* Username and password: This is your Oslo Kommune AD organization username (oooXXXXX) and its password.
* Multi factor token (MFA): The same one you use to login to
[AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws). If you haven't set up MFA yet, you
can do that [here](https://login.oslo.kommune.no/auth/realms/AD/account/totp).
* AWS account ID: This identifies which account you want to use. You can see which accounts you have access to just
after logging in to [AWS](https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws):

<span style="display:block;text-align:center">![okctl](../img/aws-account-id.png)</span>
