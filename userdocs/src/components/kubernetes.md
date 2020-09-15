`okctl` relies on services in AWS and Github to provide its functionality. In the following sections we describe some of the core services we install into the Kubernetes cluster.

## Kubernetes components

Kubernetes is [highly configurable and extensible](https://kubernetes.io/docs/concepts/extend-kubernetes/). We make use of this fact with `okctl` and bind our Kubernetes cluster setup tighter together with AWS. We do this by deploying a variety of applications into the cluster to manage various aspects of an application setup for us. 

- [Kubernetes External Secrets](#kubernetes-external-secrets)
- [AWS ALB Ingress Controller](#aws-alb-ingress-controller)
- [ExternalDNS](#externaldns)

### Kubernetes External Secrets

[Kubernetes External Secrets](https://github.com/godaddy/kubernetes-external-secrets/) allows you to use external secret management systems, like AWS Secrets Manager or HashiCorp Vault, to securely add secrets in Kubernetes.

We have installed external secrets and configured it to use [SSM Parameter store](#aws-systems-manager-amazon-ssm-parameter-store) as a backend. This means that we can [store secrets in SSM](https://github.com/godaddy/kubernetes-external-secrets/#add-a-secret) and eventually have them made available as a Kubernetes [secret](https://kubernetes.io/docs/concepts/configuration/secret/) resource that we can reference in our [deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) manifests.

### AWS ALB Ingress Controller

[AWS ALB Ingress Controller](https://github.com/kubernetes-sigs/aws-alb-ingress-controller) satisfies the Kubernetes [ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) resources by provisioning [Application Load Balancers](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) (ALB). An ALB  functions at the application layer, the seventh layer of the Open Systems Interconnection (OSI) model. After the load balancer receives a request, it evaluates the listener rules in priority order to determine which rule to apply, and then selects a target from the target group for the rule action. We use ALBs, among other things, to route traffic from the internet into a [pod](https://kubernetes.io/docs/concepts/workloads/pods/) (container).

We have configured AWS ALB Ingress Controller to work with a cluster, which means that a user doesn't have to manage the life-cycle of an ALB outside of their cluster. Instead, one can attach [annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) to the ingress object and have the [controller create these for one](https://kubernetes-sigs.github.io/aws-alb-ingress-controller/guide/controller/config/).

### ExternalDNS

[ExternalDNS](https://github.com/kubernetes-sigs/external-dns) synchronizes exposed Kubernetes Services and Ingresses with DNS providers.

We have configured ExternalDNS to work with [Route53](#aws-route53-route53), which looks at the host part of an ingress definition and creates a DNS entry for that host. This functionality works in tandem with the [AWS ALB Ingress Controller](#aws-alb-ingress-controller) so traffic is routed to the correct deployment.
