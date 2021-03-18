package lambdapermission_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/lambdapermission"
	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "LambdaPermission",
			Golden: "lambda-permission.yaml",
			Content: lambdapermission.NewRotateLambdaPermission(
				"myPermission",
				tstr.NewNameReferencer("myFunction"),
			),
		},
	}

	tstr.RunTests(t, testCases, nil)
}
