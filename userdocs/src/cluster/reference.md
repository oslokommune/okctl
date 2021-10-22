<style>
sub.required {
  color: #ff8274;
}

span.default,
span.example,
span.type {
  font-weight: bold;
}
</style>

The following attributes are available for modifying an okctl environment

## `metadata`

<div>
    <span class="type">Type</span>:
    Object
</div>

### `accountID` <sub class="required">required</sub>

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="example">Example</span>:
    123456789123
</div>

The AWS Account ID where your okctl environment will live. This is the account that will own all resources that `okctl`
sets up.

### `name` <sub class="required">required</sub>

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="example">Example</span>:
    kjoremiljo-prod
</div>

Name can be anything, but should define the scope of the cluster. Meaning if the cluster is scoped to one product, you
might want to name it the name of the product. If the cluster contains all services and products owned by a team, the
team name might be more fitting.

### `region`

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="default">Default</span>:
    eu-west-1
</div>

Region defines the AWS region to prefer when creating resources

## `clusterRootDomain` <sub class="required">required</sub>

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="example">Example</span>:
    kjoremiljo-prod.oslo.systems
</div>

The cluster root domain defines the domain of which to create services beneath. For example; `okctl` will setup ArgoCD
which has a frontend. The frontend will be available at https://argocd.&lt;clusterRootDomain&gt;.

For example, for Cognito it will be https://auth.&lt;clusterRootDomain&gt;

## `github`

<div>
    <span class="type">Type</span>:
    Object
</div>

For `okctl` to be able to setup ArgoCD correctly for you, it needs to know what repository on Github that will contain
your infrastructure.

### `repository` <sub class="required">required</sub>

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="example">Example</span>:
    okctl
</div>

The name of the repository

### `organization`

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="default">Default</span>:
    oslokommune
</div>

The organization that owns the repository

### `outputPath`

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="default">Default</span>:
    infrastructure
</div>

The folder to place infrastructure files

## `users`

<div>
    <span class="type">Type</span>:
    List
</div>

Defines which users can access everything connected to Cognito. Applications connected to Cognito include: ArgoCD UI,
Grafana.

### `email`

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="example">Example</span>:
    olly@okctl.io
</div>

The email of the user to add as an admin of the `okctl` environment

## `integrations`

<div>
    <span class="type">Type</span>:
    Object
</div>

### `argoCD`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

ArgoCD is a service that watches a repository for Kubernetes charts and ensures the defined resources are running as
declared in the cluster

### `autoscaler`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

Autoscaler automatically adjusts the size of pods and nodes in your cluster depending on load

### `awsLoadBalancerController`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

AWS Load Balancer Controller handles routing from the internet to your application running inside your `okctl`
Kubernetes cluster. If you want your applications and services accessible from the internet, this needs to be enabled

### `blockstorage`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

Block storage provides persistent storage for your cluster (Persistent Volumes)

### `cognito`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

Cognito is an authentication provider that okctl uses to control access to different resources, like ArgoCD and Grafana

### `externalDNS`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

External DNS handles defining the necessary DNS records required to route traffic to your defined service or application

### `externalSecrets`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

External Secrets fetches secrets from external sources and exposes them as native Kubernetes secrets inside the cluster

### `kubePromStack`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

KubePromStack enables Prometheus and Grafana for metrics

### `promtail`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

Promtail scrapes logs from pods and feeds them to Loki

### `loki`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

Loki collects logs and exposes them as a data source in Grafana

### `tempo`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

Tempo collects traces and exposes them as a data source in Grafana. Supports formats like jaeger, zipkin, open telemetry

## `databases`

<div>
    <span class="type">Type</span>:
    Object
</div>

### `postgres`

<div>
    <span class="type">Type</span>:
    List of objects
</div>

#### `name`

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="example">Example</span>:
    okctl-db
</div>

Name defines the name of the database to provision

#### `namespace`

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="example">Example</span>:
    okctl
</div>

Namespace defines which namespace to place the database information in. There will be created a Kubernetes Secret
containing the administrator credentials and a config map for connection details.

#### `user`

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="example">Example</span>:
    admin
</div>

User defines which admin user to provision for administrative operations

## `vpc`

<div>
    <span class="type">Type</span>:
    Object
</div>

`okctl` creates a Virtual Private Cloud for you which it organizes all the intended resources that require networking. A
VPC is mandatory, but can be configured by the following attributes.

### `cidr`

<div>
    <span class="type">Type</span>:
    String
</div>
<div>
    <span class="default">Default</span>:
    192.168.0.0/20
</div>

CIDR defines the VPC IP range. Leave this be if you don't know what it is/does

### `highAvailability`

<div>
    <span class="type">Type</span>:
    Boolean
</div>
<div>
    <span class="default">Default</span>:
    true
</div>

HighAvailability means we create redundancy in the network setup. If set to true we will create a NAT gateway per public
subnet, instead of routing all traffic through one.
