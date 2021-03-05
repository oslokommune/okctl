`okctl` relies on services in AWS and Github to provide its functionality. In the following sections we describe some of the core services we install into the Kubernetes cluster.

## Kubernetes components

Kubernetes is [highly configurable and extensible](https://kubernetes.io/docs/concepts/extend-kubernetes/). We make use of this fact with `okctl` and bind our Kubernetes cluster setup tighter together with AWS. We do this by deploying a variety of applications into the cluster to manage various aspects of an application setup for us.

- [Kubernetes External Secrets](#kubernetes-external-secrets)
- [AWS ALB Ingress Controller](#aws-alb-ingress-controller)
- [ExternalDNS](#externaldns)
- [Cluster Autoscaler](#cluster-autoscaler)
- [Kubernetes Prometheus](#kubernetes-prometheus)
- [Loki](#loki)
- [Tempo](#tempo)

### Kubernetes External Secrets

[Kubernetes External Secrets](https://github.com/external-secrets/kubernetes-external-secrets/) allows you to use external secret management systems, like AWS Secrets Manager or HashiCorp Vault, to securely add secrets in Kubernetes.

We have installed external secrets and configured it to use [SSM Parameter store](#aws-systems-manager-amazon-ssm-parameter-store) as a backend. This means that we can [store secrets in SSM](https://github.com/external-secrets/kubernetes-external-secrets/#add-a-secret) and eventually have them made available as a Kubernetes [secret](https://kubernetes.io/docs/concepts/configuration/secret/) resource that we can reference in our [deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) manifests.

#### Example

Create an SSM parameter:

```bash
aws ssm put-parameter --name "/postgres/adminpass" --value "P@sSwW)rd" --type "SecureString"
```

Kubernetes External Secrets adds a [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) (CRD). We use this CRD to make use of the correct backend; SSM (systemManager) in this case and the path to the SSM secret.

```yaml
apiVersion: 'kubernetes-client.io/v1'
kind: ExternalSecret
metadata:
  name: postgres-config
spec:
  backendType: systemManager
  data:
    - key: /postgres/adminpass
      name: admin_password
```

When this definition is applied to the cluster with `kubectl apply -f {secret.yaml} --namespace {your-namespace}` it
will result in a Kubernets secret being created.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: postgres-config
type: Opaque
data:
  admin_password: ...
```

Inspect the secrets you just applied with:

```bash
kubectl get secrets --namespace {your-namespace}
kubectl describe secret postgres-config --namespace {your-namespace}
```

If you have a postgres `Deployment`, add the following `env` configuration to let your pods get the `POSTGRES_PASSWORD`
environment variable set on startup (the deployment needs to be in the same namespace as the secret):

```yaml
spec:
  containers:
  - image: your-docker-image:0.0.1
    name: application-name
    env:
      - name: POSTGRES_PASSWORD
        valueFrom:
          secretKeyRef:
            name: postgres-config
            key: admin_password
```

### AWS Load Balancer Controller

[AWS Load Balancer Controller](https://github.com/kubernetes-sigs/aws-load-balancer-controller) 


satisfies the Kubernetes [ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) resources by provisioning [Application Load Balancers](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) (ALB). An ALB  functions at the application layer, the seventh layer of the Open Systems Interconnection (OSI) model. After the load balancer receives a request, it evaluates the listener rules in priority order to determine which rule to apply, and then selects a target from the target group for the rule action. We use ALBs, among other things, to route traffic from the internet into a [pod](https://kubernetes.io/docs/concepts/workloads/pods/) (container).

We have configured AWS Load Balancer Controller to work with a cluster, which means that a user doesn't have to manage the life-cycle of an ALB outside of their cluster. Instead, one can attach [annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) to the ingress object and have the [controller create these for one](https://kubernetes-sigs.github.io/aws-alb-ingress-controller/guide/controller/config/).

#### Example

The following ingress resource will result in the creation of a public ALB. In this example, we only use a subset of the [available annotations](https://kubernetes-sigs.github.io/aws-alb-ingress-controller/guide/ingress/annotation/), but they demonstrate how you can:

 1. Configure HTTP to HTTPS redirect
    ```
        alb.ingress.kubernetes.io/actions.ssl-redirect: \
            '{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}'
    ```
2. Define the ports to listen on
    ```
        alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS":443}]'
    ```
3. The TLS certificate to use
    ```
        alb.ingress.kubernetes.io/certificate-arn: arn:::certificate/...
    ```
4. Define a custom healthcheck endpoint on the service
    ```
        alb.ingress.kubernetes.io/healthcheck-path: /health
    ```

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: test-backend-ingress
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: instance
    alb.ingress.kubernetes.io/healthcheck-path: /health
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS":443}]'
    alb.ingress.kubernetes.io/actions.ssl-redirect: '{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}'
    alb.ingress.kubernetes.io/certificate-arn: arn:::certificate/...
  labels:
    app: test-backend
spec:
  rules:
    - host: test-backend.oslo.systems
      http:
        paths:
          # A path like this is required for the SSL redirect to function
          # This rule must also be the first in the list.
          - path: /*
            backend:
              serviceName: ssl-redirect
              servicePort: use-annotation
          - path: /*
            backend:
              serviceName: test-backend-service
              servicePort: 80
```

### ExternalDNS

[ExternalDNS](https://github.com/kubernetes-sigs/external-dns) synchronizes exposed Kubernetes Services and Ingresses with DNS providers.

We have configured ExternalDNS to work with [Route53](#aws-route53-route53), which looks at the host part of an ingress definition and creates a DNS entry for that host. This functionality works in tandem with the [AWS ALB Ingress Controller](#aws-alb-ingress-controller) so traffic is routed to the correct deployment.

#### Example

To understand ExternalDNS, we can simplify the previous example, and focus on the following line:

`/spec/rules/0/host: test-backend.oslo.systems`

ExternalDNS will simply look at the defined `host` and create a Route53 DNS entry that it associates with the ALB created in the [AWS ALB Ingress Controller](#aws-alb-ingress-controller) example.

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: test-backend-ingress
  labels:
    app: test-backend
spec:
  rules:
    - host: test-backend.oslo.systems
      http:
        paths:
          - path: /*
            backend:
              serviceName: test-backend-service
              servicePort: 80
```

### Cluster Autoscaler 

[Cluster Autoscaler](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md) is a tool that automatically adjusts the size of the Kubernetes cluster. On AWS it does this by adjusting the size of the [Autoscaling Group](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html) of the cluster's nodegroup. We have configured the cluster autoscaler to automatically adjust the default nodegroup that we add to a cluster. Starting at 1 node it can scale up and down to a maximum of 10 nodes.

### Amazon Elastic Block Store (EBS) CSI driver

[Amazon Elastic Block Store (EBS) CSI driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver) provides a CSI interface used by Container Orchestrators to manage the lifecycle of Amazon EBS volumes. This makes it possible to create `PersistentVolumeClaims` in the Kubernetes cluster.

### Kubernetes Prometheus

[Kubernetes Prometheus](https://github.com/prometheus-operator/kube-prometheus) is a stack meant for cluster monitoring, so it is pre-configured to collect metrics from all Kubernetes components. In addition to that it delivers a default set of dashboards and alerting rules. For a better idea of what possibilities this stack offers take a look at some [examples](https://github.com/prometheus-operator/kube-prometheus/tree/master/examples).

### Loki

[Loki](https://grafana.com/oss/loki/) is a horizontally-scalable, highly-available, multi-tenant log aggregation system inspired by Prometheus. It is designed to be very cost effective and easy to operate. It does not index the contents of the logs, but rather a set of labels for each log stream.

### Tempo

[Tempo](https://grafana.com/oss/tempo/)  is an open source, easy-to-use and high-scale distributed tracing backend. Tempo is cost-efficient, requiring only object storage to operate, and is deeply integrated with Grafana, Prometheus, and Loki. Tempo can be used with any of the open source tracing protocols, including Jaeger, Zipkin, and OpenTelemetry.