
## Running a Docker image in your cluster

Running a Docker image in your cluster depends on where you store it. Two common places to store Docker images are

1. Github Container Registry (GHCR)
1. Elastic Container Registry (ECR)

## Access Github Container Registry (GHCR) images

For Kubernetes to be able to download Docker images from GHCR, it needs to have the necessary credentials. This
credential is called a pull secret. 

To create a Kubernetes pull secret, first go to your Github account's settings and open the developer tab. Here you can
choose the Personal Access Token tab, which will let you create a token Kubernetes can use to access GHCR.

This token needs the `read:packages` scope.

Copy this token and run 

```shell
kubectl create secret docker-registry regcred \
  --docker-server=ghcr.io \
  --docker-username=<your-name> \
  --docker-password=<enter-token-here> \
  --docker-email=<your-email>
```

## Access Elastic Container Registry

Hopefully nothing is needed to be done.

If you are suffering from imagePullBackoff, A detailed article about applying the correct policy can be found [here](https://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_on_EKS.html)


## Push a Docker image to the registry of choice

To push a Docker image 
