package dynamodb_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/dynamodb"

	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:    "DynamoDBTable",
			Golden:  "dynamodbtable.yaml",
			Content: dynamodb.New("myResource", "my-table", "id"),
		},
	}

	tstr.RunTests(t, testCases, nil)
}
