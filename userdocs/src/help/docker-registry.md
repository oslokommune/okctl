
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

## Access Elastic Container Registry (ECR)

Hopefully nothing is needed to be done.

If you are suffering from imagePullBackoff, A detailed article about applying the correct policy can be found [here](https://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_on_EKS.html)


## Push a Docker image to the registry of choice

To push a Docker image to a Docker registry, you need to do the following:

1. Log in to your registry of choice with `docker login`.
1. Tag the image you are going to push with `docker tag`. The tag needs to be prefixed with the host of the registry.
1. Run `docker push <image tag>`.

### Push a Docker image to the Github Container Registry (GHCR)

#### Log in

To login in to GHCR, you need a Github personal access token (PAT). Instructions on obtaining a PAT can be found [here](#access-github-container-registry-ghcr-images).

```shell
# Command format
# docker login ghcr.io -u GITHUB_USERNAME -p <your PAT>
```

#### Tag image

```shell
# Usage
# docker tag SOURCE_IMAGE ghcr.io/OWNER/IMAGE_NAME:VERSION
#
# SOURCE_IMAGE    The tag of a previously built or downloaded image. Can also be the image SHA.
# OWNER           Owner is either an organization name or a username.
# IMAGE_NAME      A name representing the dockerized application
# VERSION         The version of the dockerized application
# 
# Example
docker tag faea735f5ca00d0c84cbe72ac6d17522cce1e37ac9fe49ba5e3db149d55e193b ghcr.io/oslokommune/gatekeeper:1.0.41
```

More information can be found [here](https://docs.github.com/en/packages/guides/pushing-and-pulling-docker-images)

#### Push image

```shell
# Usage
# docker push IMAGE
#
# IMAGE     the full tag the image was tagged with in the previous step
#
# Example
docker push ghcr.io/oslokommune/gatekeeper:1.0.41
```

### Push a Docker image to the Amazon Elastic Container Registry (ECR)

#### Log in

```shell
# Usage
# aws ecr get-login-password --region region | docker login --username AWS --password-stdin aws_account_id.dkr.ecr.region.amazonaws.com
#
# Example
aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin 123456789012.dkr.ecr.eu-west-1.amazonaws.com
```

More information can be found [here](https://docs.aws.amazon.com/AmazonECR/latest/userguide/registry_auth.html).

#### Tag image

```shell
# Usage
# docker tag SOURCE_IMAGE AWS_ACCOUNT_ID.dkr.ecr.REGION.amazonaws.com/IMAGE
#
# SOURCE_IMAGE    The tag of a previously built or downloaded image. Can also be the image SHA.
# AWS_ACCOUNT_ID  The AWS account ID representing the account that owns the ECR
# REGION          The region where the ECR
```
