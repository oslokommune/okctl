`okctl` relies on services in AWS and Github to provide its functionality. In the following sections we describe some of the core services we use from the cloud provider.

## Cloud components

Cloud providers offer a vast array of functionality for:

 - Networking
 - Computation
 - DNS
 - Certificates
 - Databases
 - Block storage
 - Artificial intelligence

 In `okctl` we use a subset of this functionality to provide a platform for running production workloads:

- [Amazon Web Services](#amazon-web-services-aws) as the cloud provider
- [Virtual Private Cloud](#virtual-private-cloud-vpc) for network isolation
- [Elastic Kubernetes Service](#elastic-kubernetes-service-eks) for deploying and running applications
- [Route53](#aws-route53-route53) for DNS
- [Certificate Manager](#aws-certificate-manager-acm) for issuing SSL/TLS certificates for secure communication
- [Systems Manager Parameter Store](#aws-systems-manager-amazon-ssm-parameter-store) for storing secrets

### Amazon Web Services (AWS)

With `okctl` we use [AWS](https://aws.amazon.com/) as our cloud operator; there is no particular reason for preferring AWS over other cloud vendors, such as [Microsoft Azure](https://azure.microsoft.com/) or [Google Cloud](https://cloud.google.com/). In Oslo kommune, we can use any of these, but there are a number of teams that have greater experience with AWS.

### Virtual Private Cloud (VPC)

[Amazon Virtual Private Cloud](https://aws.amazon.com/vpc/) (Amazon VPC) lets you provision a logically isolated section of the AWS Cloud where you can launch AWS resources in a virtual network that you define. You have complete control over your virtual networking environment, including selection of your own IP address range, creation of subnets, and configuration of route tables and network gateways. You can use both IPv4 and IPv6 in your VPC for secure and easy access to resources and applications.

### Elastic Kubernetes Service (EKS)
![kubernetes](../img/kubernetes.png){: style="height:150px;width=auto;aligned:center;margin-left:auto;margin-right:auto;display:block;"}

[Amazon Elastic Kubernetes Service](https://aws.amazon.com/eks/) (Amazon EKS) is a fully managed [Kubernetes](https://kubernetes.io/) (k8s) service. Being a fully managed service, AWS ensures that the control plane is secure, reliable and scalable. This allows us to focus more on application security.

K8s is an open-source system for automating deployment, scaling, and management of containerized applications. It provides a powerful platform to build applications on top of.

### AWS Route53 (Route53)

[AWS Route53](https://aws.amazon.com/route53/) (Route53) is a highly available and scalable [Domain Name System](https://en.wikipedia.org/wiki/Domain_Name_System) (DNS) web service. It is designed to give developers and businesses an extremely reliable and cost effective way to route end users to Internet applications by translating names like www.example.com into the numeric IP addresses like 192.0.2.1 that computers use to connect to each other.

```bash
dig test.oslo.systems NS +short

ns-327.awsdns-40.com.
ns-612.awsdns-12.net.
ns-1706.awsdns-21.co.uk.
ns-1322.awsdns-37.org.
```

### AWS Certificate Manager (ACM)

[AWS Certificate Manager](https://aws.amazon.com/certificate-manager/) (ACM) lets you easily provision, manage, and deploy public and private Secure Sockets Layer/Transport Layer Security (SSL/TLS) certificates for use with AWS services, and your internal connected resources. SSL/TLS certificates secure network communication and establish the identity of websites over the Internet as well as resources on private networks.

```bash
curl -vvI https://argocd.veiviser.oslo.systems

<snip>
* Server certificate:
*  subject: CN=argocd.veiviser.oslo.systems
*  start date: Jun 17 00:00:00 2020 GMT
*  expire date: Jul 17 12:00:00 2021 GMT
*  subjectAltName: host "argocd.veiviser.oslo.systems" matched cert's "argocd.veiviser.oslo.systems"
*  issuer: C=US; O=Amazon; OU=Server CA 1B; CN=Amazon
*  SSL certificate verify ok.
<snip>
```

### AWS Systems Manager (Amazon SSM) Parameter Store

[AWS Systems Manager](https://aws.amazon.com/systems-manager/) (Amazon SSM) gives you visibility and control of your infrastructure on AWS. Systems Manager provides a unified user interface so you can view operational data from multiple AWS services and allows you to automate operational tasks across your AWS resources.

The parameter store provides centralized storage and management of secrets and configuration data such as passwords, database strings, and license codes. We can encrypt values, or store as plain text, and secure access at every level.
