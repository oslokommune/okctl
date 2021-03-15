package lambdafunction_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/lambdafunction"
	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "RotateLambda",
			Golden: "rotate-lambda.yaml",
			Content: lambdafunction.NewRotateLambda(
				"myRotater",
				"mybucket",
				"mykey",
				tstr.NewNameReferencer("myRole"),
				tstr.NewNameReferencer("mySecGroup"),
				[]string{"subnet893u290uf", "subnet90uf03j"},
			),
		},
	}

	tstr.RunTests(t, testCases)
}
