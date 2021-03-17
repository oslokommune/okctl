# Secrets Manager Rotation Schedule

We need to create a custom AWS Lambda Function, because of the permissions boundary set on our AWS accounts.

The code for the lambda requires some python modules, so we use a docker container to fetch the dependencies and zip this down into an archive that is compatible with the AWS Lambda framework.

## Usage

1. Makes modifications to the code or structure
2. Run: `make package`
3. This should result in a new lambda_function.zip with the new content
4. The zip file will be embedded into the rotater.go file
5. Success

## Code

We have slightly modified the code found at this location: https://github.com/aws-samples/aws-secrets-manager-rotation-lambdas/blob/master/SecretsManagerRDSPostgreSQLRotationSingleUser