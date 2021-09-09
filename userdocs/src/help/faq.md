
## How do I grant access to a cluster to a team member?

### Prerequisites

Your team member needs the following privileges

1. Access to the AWS account where the cluster runs
1. They need access to the IAC repo on GitHub that you run `okctl` in

### Run `kubectl` commands

Your team member can now

1. Clone and enter the IAC repository
1. Run `okctl venv -c <cluster declaration`

Your team member should now be able to use `kubectl` commands

## Which sub-domain names are set up?

Run this to find out:

```bash
kubectl get ingress --all-namespaces
```

## Okctl is expecting an oslokommune-boundary to be present, but it's missing

You're probably trying to create an okctl cluster on a Crayon account. We've yet to adapt okctl to work on the new
accounts, so until then you can run the following command to create a dummy policy in the new account.

```shell
aws iam create-policy \
  --policy-name oslokommune-boundary \
  --path /oslokommune/ \
  --policy-document "{\"Version\": \"2012-10-17\", \"Statement\": [ {\"Sid\": \"AllowAccessToAllServices\", \"Effect\": \"Allow\", \"NotAction\": [\"iam:CreateUser\"], \"Resource\": \"*\"}]}"
```
