To help manage the application lifecycle we provide addons that aim to make your life easier as a developer. While there aren't many of these available at the moment, we expect this to change as new requirements in the product teams surface.

## Postgres

It is now possible to create an [AWS RDS Postgres](https://aws.amazon.com/rds/postgresql/) database from the cluster declaration.  By adding the following section to your scaffolded `cluster.yaml` file:

```yaml
metadata:
  environment: test # Just here to clarify the example
databases:
  postgres:
    - name: dbtest
      user: administrator
      namespace: dbtest
```

Followed by applying the updated declaration with:

```bash
$ okctl apply cluster -f cluster.yaml
```

This will create a Postgres database in your AWS account together with a number of related resources. We require the user to provide a Kubernetes Namespace, the reason for this is because several of these resources need to be created in the Kubernetes cluster itself.

- A secrets manager secret for the administrator password
- Lambda function that rotates the administrator password every 30 days
- A security group that allows Pods that have that association to connect with the RDS Postgres database
- Create a Secret and ConfigMap in the namespace that provide the parameters required to connect to the database

### Available parameters

We currently provide the following parameters through the Secret:

- PGUSER
- PGPASSWORD

and through the ConfigMap:

- PGHOST
- PGPORT
- PGDATABASE
- PGSSLMODE

When combined these can be used to create a valid postgresql connection string.

### DNS Policy

In Kubernetes, it is possible to define a [DNS Policy](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy) for a pod. For a Pod to be able to connect to the Postgres database it needs permissions via an attached security group, but also need to be able to resolve the PGHOST. To resolve the PGHOST, it might be necessary to set the `dnsPolicy` of the pod to `Default`, which means that the Pod inherits the node's DNS policy. 

### Attaching to the database from a local machine

We also provide a convenience function for connecting to the Postgres database from your local machine with minimal effort:

```bash
$ okctl attach postgres --environment test --name dbtest
```

### Security Group Policy

Per the documentation on [security groups for pods](https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html), the way we associate a VPC security group in AWS with a Kubernetes Pod is by declaring a `SecurityGroupPolicy` manifest.

```yaml
apiVersion: vpcresources.k8s.aws/v1beta1
kind: SecurityGroupPolicy
metadata:
  name: <my-security-group-policy>
  namespace: <my-namespace>
spec:
  <podSelector>: 
    matchLabels:
      <role>: <my-role>
  securityGroups:
    groupIds:
      - <sg-abc123>
```

In the `SecurityGroupPolicy` manifest, one can select which Pods should be associated with the security group by using the `podSelector` or `serviceAccountSelector`. Either will match on labels associated with the service accounts or pods in question.

### Security groups for Pods and Fargate

Per now, it is not possible to associate security groups with pods running in fargate, this is on the roadmap however, the following [issue](https://github.com/aws/containers-roadmap/issues/625) is tracking the progress.