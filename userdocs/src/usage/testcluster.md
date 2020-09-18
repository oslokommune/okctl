Sometimes it can be useful to setup a minimal cluster for testing or experimentation, which doesn't consume large amounts of resources. With `okctl` we provide such functionality through the `testcluster`.

## Command

To create a `testcluster` run the following command and follow the guide:

```bash
$ okctl create testcluster exp 123456789012
```

## Details

The `testcluster` reduces overhead by omitting a number of integrations. This means that we don't integrate with [Github](../components/github.md) or setup [Argo CD](../deployment/argocd.md), but we do configure the following integrations:

- [AWS ALB Ingress Controller](../components/kubernetes.md#aws-alb-ingress-controller) for routing traffic from the internet into a pod
- [ExternalDNS](../components/kubernetes.md#externaldns) for human consumable hostnames/URLs
- [Kubernetes External Secrets](../components/kubernetes.md#kubernetes-external-secrets) for an SSM Parameter Store integration

The `testcluster` command also uses the smaller `t2.medium` instances when creating a [NodeGroup](https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html). In addition, we only configure one [NAT Gateway](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html) for all the private subnets, which means we don't have a high availability setup, but this shouldn't be necessary for a test cluster and it reduces the number of [EIPs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html) we need.

## Usage

After the cluster has been created you can save the following definitions under a directory, e.g, `experiments/2048/`.

**NB:** Remember to change the values in brackets `{...}` so they match your cluster setup.

`experiments/2048/2048-game.yml`:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: "2048-game"
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: "2048-ingress"
  namespace: "2048-game"
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
  labels:
    app: 2048-ingress
spec:
  rules:
    - host: 2048-game.{hosted_zone}.oslo.systems
      http:
        paths:
          - path: /*
            backend:
              serviceName: "service-2048"
              servicePort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: "service-2048"
  namespace: "2048-game"
spec:
  ports:
    - port: 80
      targetPort: 80
      protocol: TCP
  type: NodePort
  selector:
    app: "2048"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: "2048-deployment"
  namespace: "2048-game"
spec:
  selector:
    matchLabels:
      app: "2048"
  replicas: 5
  template:
    metadata:
      labels:
        app: "2048"
    spec:
      containers:
      - image: alexwhen/docker-2048
        imagePullPolicy: Always
        name: "2048"
        ports:
        - containerPort: 80
```

Then you can fetch your credentials using `okctl`:

```bash
$ okctl show credentials {env}
```

This command will output a number of environment variables you can export. You need to use the ones for `AWS` and `KUBECONFIG` and export them in a shell.

Then you can apply the manifest you saved earlier to the cluster:

```bash
$ kubectl apply -f experiments/2048/2048-game.yml
```

Give `AWS ALB Ingress Controller` and `ExternalDNS` some time to work their magic, and eventually you should be able to access and play your game at the URL: `2048-game.{hosted_zone}.oslo.systems`.

When you are done, simply delete the manifest:

```bash
$ kubectl delete -f experiments/2048/2048-game.yml
```

And delete the cluster, if you feel done:

```bash
$ okctl delete testcluster {env}
```