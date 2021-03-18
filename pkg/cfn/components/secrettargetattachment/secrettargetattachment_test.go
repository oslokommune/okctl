package secrettargetattachment_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/secrettargetattachment"
	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "RDSInstanceSecretTargetAttachment",
			Golden: "rdsinstance-secretattachment.json",
			Content: secrettargetattachment.NewRDSDBInstance(
				"myAttachment",
				tstr.NewNameReferencer("mySecret"),
				tstr.NewNameReferencer("myRDSInstance"),
			),
		},
	}

	tstr.RunTests(t, testCases, nil)
}
