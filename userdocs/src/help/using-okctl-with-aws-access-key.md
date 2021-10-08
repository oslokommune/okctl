## Using okctl with an access key from AWS

If you want to use okctl with an access key from AWS instead of okctl's default authentication, follow these steps.

### Get an access key

* Log in to your AWS console ➡ IAM ➡ Users ➡ Pick one ➡ Security credentials ➡ Create access key. Note down access key ID and secret.

### Log in to your cluster

Run

```shell
export AWS_ACCESS_KEY_ID=<your access key id>
export AWS_SECRET_ACCESS_KEY=<your access key secret>

okctl -a access-key <the okctl command you want to run> 
```

For instance:

```shell
export AWS_ACCESS_KEY_ID=someid
export AWS_SECRET_ACCESS_KEY=somesecret

okctl -a access-key venv -c my-cluster.yaml 
```

Reference:

* [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)

### Get a team member to give you access

Someone with access to the okctl cluster must give you access by following the steps below.

#### Step 2

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

#### Step 3

Create a file `mycluster-access-list.yaml` with the contents (or edit a file, if you have run this step before):

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
    name: someone@email.com
  - apiGroup: rbac.authorization.k8s.io
    kind: User
    name: someone@email.com

```

Under `subjects`, add a list of all the users who shall have access to the cluster. This should be the same list as in `mapUsers` mentioned above.

Then run

```yaml
kubectl apply -f mycluster-access-list.yaml
```

Reference:

* [AWS documentation](https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html)

### Verify

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
