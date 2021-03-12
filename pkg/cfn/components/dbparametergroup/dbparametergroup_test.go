package dbparametergroup_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/dbparametergroup"

	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "DB Parameter Group",
			Golden: "dbparametergroup.json",
			Content: dbparametergroup.New("myParameterGroup", map[string]string{
				"some_parameter": "some_value",
			}),
		},
	}

	tstr.RunTests(t, testCases)
}
