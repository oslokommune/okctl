package rotationschedule_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/rotationschedule"
	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "PostgresRotationSchedule",
			Golden: "postgres-rotationschedule.json",
			Content: rotationschedule.NewPostgres(
				"myPostgresRotationSchedule",
				tstr.NewNameReferencer("mySecret"),
				tstr.NewNameReferencer("mySecretAttachment"),
				[]string{"subnetid-3fjij3", "subnetid-jfe29md"},
				tstr.NewNameReferencer("mySecurityGroup"),
			),
		},
	}

	tstr.RunTests(t, testCases)
}
