// Package lambdafunction builds a cloud formation resource for Lambda Functions
package lambdafunction

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/lambda"
	"github.com/gosimple/slug"
)

// LambdaFunction contains all the required state for creating a Lambda Function
type LambdaFunction struct {
	StoredName string
	Opts       Opts
}

// Opts contains the configurable options
type Opts struct {
	FunctionName    string
	Handler         string
	Runtime         string
	Bucket          string
	Key             string
	Env             map[string]string
	Role            cfn.Namer
	SecurityGroupID cfn.Namer
	SubnetIds       []string
}

const (
	timeoutInSeconds = 30
)

// Resource returns the cloud formation resource
func (l *LambdaFunction) Resource() cloudformation.Resource {
	return &lambda.Function{
		Code: &lambda.Function_Code{
			S3Bucket: l.Opts.Bucket,
			S3Key:    l.Opts.Key,
		},
		Environment: &lambda.Function_Environment{
			Variables: l.Opts.Env,
		},
		FunctionName: l.Opts.FunctionName,
		Handler:      l.Opts.Handler,
		Role:         cloudformation.GetAtt(l.Opts.Role.Name(), "Arn"),
		Runtime:      "python3.7",
		Timeout:      timeoutInSeconds,
		VpcConfig: &lambda.Function_VpcConfig{
			SecurityGroupIds: []string{
				cloudformation.GetAtt(l.Opts.SecurityGroupID.Name(), "GroupId"),
			},
			SubnetIds: l.Opts.SubnetIds,
		},
	}
}

// Name returns the name of the resource
func (l *LambdaFunction) Name() string {
	return l.StoredName
}

// Ref returns a reference to the resource
func (l *LambdaFunction) Ref() string {
	return cloudformation.Ref(l.Name())
}

// NamedOutputs returns the outputs
func (l *LambdaFunction) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(l.Name(), cloudformation.GetAtt(l.Name(), "Arn")).NamedOutputs()
}

// New returns an initialised lambda function resource
func New(resourceName string, opts Opts) *LambdaFunction {
	return &LambdaFunction{
		StoredName: resourceName,
		Opts:       opts,
	}
}

// NewRotateLambda returns an initialised lambda for rotating secrets
func NewRotateLambda(
	resourceName, bucket, key string,
	role cfn.Namer,
	securityGroup cfn.Namer,
	subnetIDs []string,
) *LambdaFunction {
	return New(resourceName, Opts{
		FunctionName:    fmt.Sprintf("%s-Rotater", slug.Make(resourceName)),
		Handler:         "lambda_function.lambda_handler",
		Runtime:         "python3.7",
		Bucket:          bucket,
		Key:             key,
		Role:            role,
		SecurityGroupID: securityGroup,
		SubnetIds:       subnetIDs,
	})
}
