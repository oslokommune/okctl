// Package lambdapermission knows how to create lambda permission
// cloud formation resources
package lambdapermission

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/lambda"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// LambdaPermission contains the state required for creating a
// lambda permission resource
type LambdaPermission struct {
	StoredName string
	Principal  string
	Function   string
}

// Resource returns the cloud formation resource
func (l *LambdaPermission) Resource() cloudformation.Resource {
	return &lambda.Permission{
		Action:        "lambda:InvokeFunction",
		FunctionName:  l.Function,
		Principal:     l.Principal,
		SourceAccount: cloudformation.Ref("AWS::AccountId"),
	}
}

// Name returns the resource name
func (l *LambdaPermission) Name() string {
	return l.StoredName
}

// Ref returns a reference to the resource
func (l *LambdaPermission) Ref() string {
	return cloudformation.Ref(l.Name())
}

// NamedOutputs returns the named outputs
func (l *LambdaPermission) NamedOutputs() map[string]cloudformation.Output {
	return nil
}

// New returns a lambda assume permission policy
func New(resourceName, principal, function string) *LambdaPermission {
	return &LambdaPermission{
		StoredName: resourceName,
		Principal:  principal,
		Function:   function,
	}
}

// NewRotateLambdaPermission returns a lambda assume permission policy for a
// secrets manager rotation schedule
func NewRotateLambdaPermission(resourceName string, function cfn.Namer) *LambdaPermission {
	return New(
		resourceName,
		"secretsmanager.amazonaws.com",
		cloudformation.GetAtt(function.Name(), "Arn"),
	)
}
