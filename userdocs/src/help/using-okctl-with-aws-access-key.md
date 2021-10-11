## Motivation

If you want to use okctl with an access key from AWS instead of okctl's default authentication, follow these steps.

## Get an access key

Log in to your AWS console ➡ IAM ➡ Users ➡ Pick one ➡ Security credentials ➡ Create access key.

Note down access key ID and secret.

<span style="display:block;text-align:center">![create-access-key](/img/create-access-key.png)</span>

## Log in to your cluster

Run

```shell
# Usage
# export AWS_ACCESS_KEY_ID=<your access key id>
# export AWS_SECRET_ACCESS_KEY=<your access key secret>
#
# okctl -a access-key venv --cluster-declaration <cluster declaration 

# Example
export AWS_ACCESS_KEY_ID=someid
export AWS_SECRET_ACCESS_KEY=somesecret

okctl -a access-key venv -c my-cluster.yaml
```

Reference: 
  [AWS documentation - CLI configure envvars](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)

## Get a team member to give you access

Someone with access to the okctl cluster must give you access by following the steps below.

### Step 1

Log in to the cluster with `okctl venv` as described above.

```shell
kubectl edit configmap -n kube-system aws-auth
```

Add an element to the `mapUsers` field, so it looks like this:

```yaml
apiVersion: v1
data:
  mapRoles: |
    ...
  mapUsers: |
    - userarn: arn:aws:iam::123456789012:user/clustowner@email.com
          username: clustowner@email.com
          groups:
          - system:master
    - userarn: arn:aws:iam::123456789012:user/someone@email.com
          username: someone@email.com
          groups:
          - system:master
```

Replace

* `123456789012` with AWS account number
* `clustowner@email.com` with the e-mail of the user giving access 
* `someone@email.com` with the e-mail of the user who wants access 

There could be more users listed here.

### Step 2

Create a file `mycluster-access-list.yaml` with the contents (or edit the file, if you have run this step before):

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: mycluster-access-list
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: User
    name: clustowner@email.com
  - apiGroup: rbac.authorization.k8s.io
    kind: User
    name: someone@email.com
```

Under `subjects`, edit the list so that it contains all the users who shall have access to the cluster. This should be the same list as in `mapUsers` mentioned above.

Then run

```yaml
kubectl apply -f mycluster-access-list.yaml
```

Reference: [AWS documentation - Add user role](https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html)

## Verify

Now, the user who wants access can verify that things work by running

```shell
export AWS_ACCESS_KEY_ID=someid
export AWS_SECRET_ACCESS_KEY=somesecret

okctl -a access-key venv -c my-cluster.yaml

kubectl get pods
```

This should give no errors - either a list of pods, or just the message

```
No resources found in default namespace.
```

That's it. Now you are able to run all okctl commands with the `-a access-key` option, which tells okctl to use the provided access key instead of using the default authentication method.
