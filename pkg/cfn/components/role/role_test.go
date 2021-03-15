package role_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/components/policydocument"
	"github.com/oslokommune/okctl/pkg/cfn/components/role"
	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "IAM Role",
			Golden: "role.json",
			Content: role.New(
				"myIAMRole",
				v1alpha1.PermissionsBoundaryARN("123456789012"),
				[]string{"arn:::policy/some-policy"},
				policydocument.PolicyDocument{
					Version: policydocument.Version,
					Statement: []policydocument.StatementEntry{
						{
							Effect: policydocument.EffectTypeAllow,
							Action: []string{
								"someAction",
							},
							Resource: []string{
								"*",
							},
						},
					},
				},
				nil,
			),
		},
	}

	tstr.RunTests(t, testCases)
}
