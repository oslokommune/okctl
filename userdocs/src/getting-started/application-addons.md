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

This will create a Postgres database in your AWS account together with a number of related resources:

- A secrets manager secret for the administrator password
- Lambda function that rotates the administrator password every 30 days
- A security group that allows Pods that have that association to connect with the RDS Postgres database
- Create an ExternalSecret and ConfigMap in the namespace that provide all the information required to connect to the database

### Attaching to the database

We also provide a convenience function for connecting to the Postgres database from your local machine with minimal effort:

```bash
$ okctl attach postgres --environment test --name dbtest
```
