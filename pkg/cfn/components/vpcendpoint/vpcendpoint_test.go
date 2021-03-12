package vpcendpoint_test

import (
	"testing"

	tstr "github.com/oslokommune/okctl/pkg/cfn/components/testing"
	"github.com/oslokommune/okctl/pkg/cfn/components/vpcendpoint"
)

func TestNew(t *testing.T) {
	testCases := []tstr.CloudFormationTemplateTestCase{
		{
			Name:   "SecretsManagerVPCEndpoint",
			Golden: "sm-vpcendpoint.json",
			Content: vpcendpoint.NewSecretsManager(
				"mySecretsManagerEndpoint",
				tstr.NewNameReferencer("mySecurityGroup"),
				"vpcid-30fejkjhs",
				[]string{"subnetid-0ei0fgi", "subnetid-0ie0fie"},
			),
		},
	}

	tstr.RunTests(t, testCases)
}
