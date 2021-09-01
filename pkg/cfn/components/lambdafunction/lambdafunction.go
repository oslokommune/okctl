// Package lambdafunction builds a cloud formation resource for Lambda Functions
package lambdafunction

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	patchpkg "github.com/oslokommune/okctl/pkg/jsonpatch"

	jsonpatch "github.com/evanphx/json-patch/v5"

	"sigs.k8s.io/yaml"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/lambda"
)

// LambdaFunction contains all the required state for creating a Lambda Function
type LambdaFunction struct {
	StoredName string
	Opts       Opts
}

// Opts contains the configurable options
type Opts struct {
	Description     string
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
		Description: l.Opts.Description,
		Handler:     l.Opts.Handler,
		Role:        cloudformation.GetAtt(l.Opts.Role.Name(), "Arn"),
		Runtime:     "python3.8",
		Timeout:     timeoutInSeconds,
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
	secretsManagerVPCEndpoint cfn.Namer,
) *LambdaFunction {
	return New(resourceName, Opts{
		Description: "RDS Postgres Rotater",
		Handler:     "lambda_function.lambda_handler",
		Runtime:     "python3.7",
		Bucket:      bucket,
		Key:         key,
		Role:        role,
		Env: map[string]string{
			"SECRETS_MANAGER_ENDPOINT": cloudformation.Select(
				0,
				[]string{cloudformation.GetAtt(secretsManagerVPCEndpoint.Name(), "DnsEntries")},
			),
		},
		SecurityGroupID: securityGroup,
		SubnetIds:       subnetIDs,
	})
}

// PatchRotateLambda patches the rotater lambda
// nolint: lll
func PatchRotateLambda(lambdaResourceName, secretsManagerVPCEndpointName string, template []byte) ([]byte, error) {
	patchJSON, err := patchpkg.New().Add(patchpkg.Operation{
		Type: patchpkg.OperationTypeReplace,
		Path: fmt.Sprintf("/Resources/%s/Properties/Environment/Variables/SECRETS_MANAGER_ENDPOINT", lambdaResourceName),
		Value: &patchpkg.Inline{Data: []byte(
			fmt.Sprintf(`{"Fn::Join":["/",["https:/",{"Fn::Select":["1",{"Fn::Split":[":",{"Fn::Select":["0",{"Fn::GetAtt":["%s","DnsEntries"]}]}]}]}]]}`, secretsManagerVPCEndpointName),
		)},
	}).MarshalJSON()
	if err != nil {
		return nil, err
	}

	jsonData, err := yaml.YAMLToJSON(template)
	if err != nil {
		return nil, fmt.Errorf(constant.ConvertJsonToYamlError, err)
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return nil, fmt.Errorf(constant.DecodePatchError, err)
	}

	modified, err := patch.Apply(jsonData)
	if err != nil {
		return nil, fmt.Errorf(constant.ApplyPatchError, err)
	}

	return yaml.JSONToYAML(modified)
}
