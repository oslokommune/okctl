# Thinks worth taking note of
Just a dump of information that can be worth taking note of, as a new developer on the team it can be smart to look through some of these resources.

## Security groups for pods

There are some things worth considering about security groups for pods, the following is a small sample of the information you can find [here](https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html):

- Security groups for pods can't be used with pods deployed to Fargate.
- Not all EC2 instance types support security groups for, a complete list can be found [here](https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html#supported-instance-types)
- If you're also using pod security policies to restrict access to pod mutation, then the eks-vpc-resource-controller and vpc-resource-controller Kubernetes service accounts must be specified in the Kubernetes ClusterRoleBinding for the the Role that your [pod security policy (PSP)](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) is assigned to.
- Pods using security groups must contain terminationGracePeriodSeconds in their pod spec. This is because the Amazon EKS VPC CNI plugin queries the API server to retrieve the pod IP address before deleting the pod network on the host. Without this setting, the plugin doesn't remove the pod network on the host.

### Issues and their solution

Had an issue where a Pod with a Security Group attached wasn't resolving DNS correctly. After reading thrugh this: https://aws.amazon.com/premiumsupport/knowledge-center/eks-dns-failure/ discussion, the point on selecting the `dnsPolicy` for the Pod proved to be the correct solution. By changing from `ClusterFirst` (the actual default) to `Default`, the problem was resolved. With the `Default` setup the Pod inherits the DNS configuration from the Node it is running on, when we want to resolve RDS hostnames, this appears to be a winning strategy.