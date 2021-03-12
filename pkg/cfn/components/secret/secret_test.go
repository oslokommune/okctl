package secret_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/secret"
	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:    "RDSInstanceSecret",
			Golden:  "rdsinstance-secret.json",
			Content: secret.NewRDSInstanceSecret("myAdminSecret", "admin"),
		},
	}

	tstr.RunTests(t, testCases)
}
