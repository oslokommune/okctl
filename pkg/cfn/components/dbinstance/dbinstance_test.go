package dbinstance_test

import (
	"testing"

	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/dbinstance"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "Template for DBInstance",
			Golden: "postgres.json",
			Content: dbinstance.New(
				"myPostgresDatabase",
				"databaseIdentifier",
				"databaseName",
				"myDbSubnetGroupName",
				tstr.NewNameReferencer("myDbParameterGroup"),
				tstr.NewNameReferencer("myMonitoringRole"),
				tstr.NewNameReferencer("myAdminSecret"),
				tstr.NewNameReferencer("mySecurityGroup"),
			),
		},
	}

	tstr.RunTests(t, testCases)
}
