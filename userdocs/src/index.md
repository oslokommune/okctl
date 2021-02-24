# `okctl` - Opinionated and effortless infrastructure and application management

`okctl` is a CLI for creating a production environment on [Amazon Web Services](https://aws.amazon.com/) (AWS) using [Elastic Kubernetes Service](https://aws.amazon.com/eks/) (EKS) as its foundation, we aim to be:

<div class="center-content-container default-spacing">
    <img src="img/logo.png" alt="Okctl logo">
</div>

<!-- <span style="display:block;text-align:center">![okctl](img/logo.png)</span> -->

- **Specialised**: by tightly integrating with the services provided by Oslo kommune
- **Highly opinionated**: you will get something that works well that may not be exactly what you want
- **Modern**: base our decisions on best-practise methodologies for a modern development lifecycle
- **Sensible**: provide a set of sensible default integrations

<div class="center-content-container default-spacing">
    <a class="okctl-button primary" href="getting-started/create-cluster">Getting started</a>
    <a class="okctl-button secondary" href="getting-started/install">Install okctl</a>
</div>

## Compare and contrast

With `okctl` we are attempting to solve the production environment setup problem. What we include within the definition of a production environment, we can see below.

| Functionality | okctl | [eksctl](https://eksctl.io) | [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) | [serverless.tf](https://serverless.tf) |
|---|---|---|---|---|
|Provide authentication and authorisation|✓||?||
|Define a continuous integration pipeline|?||||
|Assist with creating a deployment pipeline|✓|✓|||
|Create a Kubernetes cluster (run-time environment)|✓|✓|||
|Facilitate with application setup and integration|✓||✓|✓|
|Provide integrations with version control system|✓||||
|Setup monitoring with metrics, logs and traces|✓||||

## Roadmap 

We have a basic cluster up and running, and we are currently working on getting the core components of the cluster deployed:

- [x] [ExternalSecrets](https://github.com/external-secrets/kubernetes-external-secrets/) for storing secrets securely
- [x] [AWS ALB Ingress Controller](https://github.com/kubernetes-sigs/aws-alb-ingress-controller) creates load balancers for incoming traffic
- [x] [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) ensures humane DNS hostnames 
- [x] [Autoscaler](https://github.com/kubernetes/autoscaler/) for adjusting the size of pods and nodes
- [x] [Argo CD](https://github.com/argoproj/argo-cd) gives us continuous delivery
- [x] [Kubernetes Prometheus](https://github.com/prometheus-operator/kube-prometheus) for monitoring
- [x] [Amazon Elastic Block Store (EBS) CSI driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/) for block storage
- [ ] [Loki](https://github.com/grafana/loki) for consuming logs
- [ ] [Jaeger](https://www.jaegertracing.io/) for traces

## Inspiration

We stand on the shoulders of giants, we have begged, borrowed and stolen from the following great products and teams:

- [Porter](https://github.com/deislabs/porter)
- [eksctl](https://github.com/weaveworks/eksctl)
- [kops](https://github.com/kubernetes/kops)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [go-kit](https://github.com/go-kit/kit)
- [saml2aws](https://github.com/Versent/saml2aws/)
