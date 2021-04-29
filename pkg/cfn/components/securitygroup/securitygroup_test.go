package securitygroup_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/securitygroup"
	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "PostgresIncoming",
			Golden: "postgres-incoming-sg.json",
			Content: securitygroup.NewPostgresIncoming(
				"myIncomingPG",
				"myIncomingPG",
				"vpcid-r3ufh3",
				tstr.NewNameReferencer("mySourceSecurityGroup"),
			),
		},
		{
			Name:   "PostgresOutgoing",
			Golden: "postgres-outgoing-sg.json",
			Content: securitygroup.NewPostgresOutgoing(
				"myIncomingPG",
				"myIncomingPG",
				"vpcid-wof03ef3",
				[]string{"192.168.1.0/20", "192.168.2.0/20"},
			),
		},
	}

	tstr.RunTests(t, testCases, nil)
}
