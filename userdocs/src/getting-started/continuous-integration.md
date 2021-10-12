#Continuous Integration

The following guide is to help set up continuous integration for an application running on a cluster set up with okctl.
For this example and in the reference app we will be using github actions

## Prerequisites

It is assumed that you already have set up a cluster, and that you have applied your application so that it runs there
Now you are ready to set up continuous integration, so that every push to main will deploy your app
to your dev cluster, and a tag to main will deploy your app to your production cluster.

You need a IAM user with credentials that you can use for github actions. You need to create credentials for both your
and prod environment: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey

Github repositories used in this guide as a working example
1. Application repository (https://github.com/oslokommune/okctl-reference-app)
2. IAC repo (https://github.com/oslokommune/okctl-reference-iac)


Generate a secret key in your IAC repo
```bash
# We will generate a key inside a separate secret directory to ensure it will not get mixed up with anything else
mkdir secret
# Gitignore that directory, so tehre will be no accidental commits
echo secret/* >> .gitignore

cd secret
# We use stronger encryption than a default key for added security
ssh-keygen -t rsa -b 4096 -f cluster_deploy_key -C git@github.com:oslokommune/okctl-reference-iac.git
```

**NOTE:** The -C parameter of the ssh-keygen command, which is the comment for the public-key, needs to be the git@github address for your IAC repo, this is because it will be needed by github actions later.

1. Go to your settings/secrets under your application
https://github.com/oslokommune/okctl-reference-app/settings/secrets/actions

1. Add **Repository**secret (it will be the same for dev and produdction) named CLUSTER_DEPLOY_KEY, paste the contnetns of
cluster_deploy_key (the private key) that you generated earlier

1. Create a **dev** and a **prod** environment (only available in github enterprise):
https://github.com/oslokommune/okctl-reference-app/settings/environments

1. Create *Environment*secrets for
   * AWS_ECR_ACCESS_KEY_ID
   * AWS_ECR_ACCESS_KEY_SECRET

    Here you will place values for AWS_ACCESS_KEY_ID and AWS_ACCESS_KEY_SECRET that you created earlier for each account
1. Go to your IAC repo (https://github.com/oslokommune/okctl-reference-iac/settings/keys)
1. Add a new deploy key called cluster_deploy_key, using the value in `cluster_deploy_key.pub`, that you generated earlier NOTE: Make sure you check the `Allow write access` checkbox

## Setup github actions workflow files

Use templates found here: https://github.com/oslokommune/okctl-reference-app/tree/main/.github/workflows

You need to edit the following:

* jobs -> docker-build-push -> steps[0] -> with -> aws-region `if you run somewhere else than Ireland`
* jobs -> docker-build-push -> steps[1] -> env -> ECR_REPOSITORY `name of your ecr_repositroy, i.e. okctl-reference-app`
* jobs -> update-tag -> steps[0] -> with -> repository `your iac-repository, i.e oslokommune/okctl-reference-iac`
* jobs -> update-tag -> steps[0] -> env -> CONTAINER_NAME `name of the container, i.e kotlin-test-app`
* jobs -> update-tag -> steps[0] -> env -> DEPLOYMENT_YAML_FILE `location of overlay deployment patch file: i.e infrastructure/applications/okctl-reference-app/overlays/okctl-reference-dev/deployment-patch.json`

