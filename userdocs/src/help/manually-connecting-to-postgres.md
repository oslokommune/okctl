### Overview

Setting up firewall rules and networking to allow for traffic from an app to a database requires the following steps:

1. Set up a `Security Group` for your application.
2. Allow traffic from the application security group to the database's incoming security group.
3. Create and apply a `Security Group Policy` to the applications pod in Kubernetes.

### Set up an application Security Group

1. In the AWS web console, navigate to the EC2 service and select `Security Groups` in the navigation menu on the left
   side of the screen.
2. In the top left corner, press `Create security group`, enter a name and description and select the `VPC` relevant to
   your cluster.
3. Add the following inbound rules with `ClusterSharedNodeGroup` as source:
    * DNS (53) for both UDP and TCP. (required by Kubernetes)
    * The port your application listens to on the protocol your application expects. (required for traffic from other
      pods and/or internet)
4. Add the following outbound rules:
    * DNS (53) for both UDP and TCP with the destination `ClusterSharedNodeGroup`
    * PostgresQL (5432) for TCP with the destination the existing security group named with the suffix
      `RDSPostgresIncoming`.
    * Every other necessary outbound rule your application depends on. For example `443` to `0.0.0.0/0` if your
      application requires general access to the internet.

### Allow access to the database from the application security group

1. In the AWS web console, navigate to the EC2 service and select `Security Groups` in the navigation menu on the left
   side of the screen.
2. Select the database's incoming security group. This security group is named in the following format:
   `<database name><cluster name>RDSPostgresIncoming`.
3. Add a new incoming rule for PostgreSQL (5432) with the source as the security group you created in the previous
   section.

### Create and apply a `Security Group Policy`

EKS uses a resource
called [Security Group Policy](https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html)
to associate security groups with pods. A `Security group policy` is a `Custom Resource Definition` (CRD) that you apply
to your Kubernetes cluster and looks like this:

```yaml
# security-group-policy.yaml
apiVersion: vpcresources.k8s.aws/v1beta1
kind: SecurityGroupPolicy
metadata:
  name: <my-security-group-policy> # Can be anything
  namespace: <my-namespace> # Must be the same namespace as the application
spec:
  podSelector: # Other selectors are also available
    matchLabels:
      <role>: <my-role>
  securityGroups:
    groupIds:
      - <sg-abc123> # The security group ID to be associated
```

Create the security group policy CRD as a YAML file and use `kubectl apply -f <file>` to apply it to your Kubernetes
cluster.

That's it. After restarting your pod, your pod should be able to connect to the database.
