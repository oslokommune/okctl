
## Running a Docker image in your cluster

Running a Docker image in your cluster depends on where you store it. Two common places to store Docker images are

1. Github Container Registry (GHCR) ([official documentation](https://docs.github.com/en/packages/guides/pushing-and-pulling-docker-images))
1. Elastic Container Registry (ECR) ([official documentation](https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html))

## Access container registries

### Access Github Container Registry (GHCR) images

For Kubernetes to be able to download Docker images from GHCR, it needs to have the necessary credentials. This credential is called a pull secret.

To create a Kubernetes pull secret, first go to your [github account settings](https://github.com/settings/profile) and select `Developer settings > Personal access tokens`. Here you can create a token Kubernetes can use to access GHCR.

For Kubernetes to be able to read packages (docker images), it needs the `read:packages` scope. For us to be able to push to GHCR, it needs the `write:packages` scope.

Copy this token and run 

```shell
kubectl create secret docker-registry regcred \
  --docker-server=ghcr.io \
  --docker-username=<your-name> \
  --docker-password=<enter-token-here> \
  --docker-email=<your-email>
```

### Access Elastic Container Registry (ECR)

Hopefully nothing is needed to be done.

If you are suffering from imagePullBackoff, A detailed article about applying the correct policy can be found [here](https://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_on_EKS.html).


## Push a Docker image to the registry of choice

To push a Docker image to a Docker registry, you need to do the following (these steps are described in detail below):

1. Log in to your registry of choice
1. Tag the image you are going to push. The tag needs to be prefixed with the host of the registry.
1. Push image

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
docker tag 9df7297819f7 ghcr.io/oslokommune/gatekeeper:1.0.41
```

More information can be found [here](https://docs.github.com/en/packages/guides/pushing-and-pulling-docker-images).

#### Push image

```shell
# Usage
# docker push TAG
#
# TAG     the full tag the image was tagged with in the previous step
#
# Example
docker push ghcr.io/oslokommune/gatekeeper:1.0.41
```

### Push a Docker image to the Amazon Elastic Container Registry (ECR)

Before you start, you need an ECR repository. It can be created in the [AWS Console](https://eu-west-1.console.aws.amazon.com/ecr/create-repository?region=eu-west-1) ([official documentation](https://docs.aws.amazon.com/AmazonECR/latest/userguide/repository-create.html)).

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
# docker tag SOURCE_IMAGE AWS_ACCOUNT_ID.dkr.ecr.REGION.amazonaws.com/IMAGE:VERSION
#
# SOURCE_IMAGE    The tag of a previously built or downloaded image. Can also be the image SHA.
# AWS_ACCOUNT_ID  The AWS account ID representing the account that owns the ECR
# REGION          The region where the ECR
# VERSION         The image version
#
# Example
docker tag 9df7297819f7 123456789012.dkr.ecr.eu-west1.amazonaws.com/gatekeeper:1.0.41
```

#### Push image

```shell
# Usage
# docker push TAG
#
# TAG     the full tag the image was tagged with in the previous step
#
# Example
docker push 123456789012.dkr.ecr.eu-west1.amazonaws.com/gatekeeper:1.0.41
```
