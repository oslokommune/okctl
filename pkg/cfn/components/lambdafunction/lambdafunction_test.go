package lambdafunction_test

import (
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

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
				tstr.NewNameReferencer("myVPCEndpoint"),
			),
		},
	}

	var template []byte

	cb := func(tmpl []byte) {
		template = tmpl
	}

	tstr.RunTests(t, testCases, cb)

	patched, err := lambdafunction.PatchRotateLambda("myRotater", "myVPCEndpoint", template)
	assert.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "rotate-lambda-patched.yaml", patched)
}
