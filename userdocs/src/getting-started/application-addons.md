To help manage the application lifecycle we provide addons that aim to make your life easier as a developer. While there aren't many of these available at the moment, we expect this to change as new requirements in the product teams surface.

## Postgres

It is now possible to create an [AWS RDS Postgres](https://aws.amazon.com/rds/postgresql/) database from the cluster declaration.  By adding the following section to your scaffolded `cluster.yaml` file:

```yaml
# cluster.yaml
...
databases:
  postgres:
    - name: dbtest
      user: administrator
      namespace: dbtest
...
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

The secret and the configmap will be placed in the Kubernetes namespace defined in the cluster declaration when provisioning the
database server.

### Forwarding traffic to the database from a local machine

If you want to attach to the database from intellij or some other IDE, you can do so by setting up a port forwarder, you can do so with the following command:

```bash
# Usage
## Generate a password for the pgBouncer user
uuidgen | sed 's/-//g' > <path to store password file>
## Forward traffic from localhost:5432 to the pgBouncer
okctl forward postgres \
  --cluster-declaration <path to cluster declaration> \
  --name <database server name> \
  --username <pgBouncer username> \
  --password-file <path to store password file>

# Example
uuidgen | sed 's/-//g' > password.secret
okctl forward postgres \
  --cluster-declaration cluster.yaml \
  --name pgtest \
  --username bob \
  --password-file password.secret
```

By default this will use the postgres default port of `5432`, but you can change this to what you want.

:information_source: The username and password used in the `forward postgres` command defines the username and password
to use when connecting to the postgres client, for example psql. The user connected to Postgres will however be the admin
user defined in `cluster.yaml` when provisioning the database.

### Attaching to the database from a local machine

We also provide a convenience function for connecting to the Postgres database from your local machine with minimal effort:

```bash
# Usage
okctl --cluster-declaration <path to cluster declaration> attach postgres --name <database server name>

# Example
okctl --cluster-declaration cluster.yaml attach postgres --name dbtest
```

The above command will drop you into a `psql` shell.

### Connecting an application with a database

#### Motivation

Setting up firewall and networking rules for connecting to a database in a production environment can be challenging.

In okctl we've implemented a way to do all this for you. By specifying the
`postgres` attribute in the application declaration, okctl will handle it for you.

#### Prerequisites

:information_source: You'll need a database provisioned with `okctl`.
See [Postgres](/getting-started/application-addons/#postgres)
on how to use `okctl` to create a PostgreSQL database.

#### Usage

```shell
# Usage in application.yaml
...
postgres: <database name>
...

# Example in application.yaml
...
postgres: gopherdb
...
```

:information_source: The **database name** is the name you've defined for your database in the cluster declaration under
the
`databases` attribute.

After you've added the `postgres` attribute to the application declaration, run

```shell
# Usage
okctl -c <cluster declaration> apply application -f <application declaration>

# Example
okctl -c cluster.yaml apply application -f application.yaml
```

to apply the changes.

Okctl will then apply the changes in three steps:

1. set up and configure the required security groups
2. create a security group policy
3. add it to the application.

Remember to add and push the changes done in the git repository for `ArgoCD` to pick up the changes.

:information_source: After you've confirmed that ArgoCD has acknowledged the changes, you might have to manually restart
the application for the changes to take effect. You can restart the pod by either restarting it in the `ArgoCD` GUI or
by running

```shell 
kubectl -n <app namespace> delete pod <pod id>
```

That's it. Your application should now be able to read and write to the database.

#### Manually setting up an app to database connection

See [this](/help/manually-connecting-to-postgres) article.
