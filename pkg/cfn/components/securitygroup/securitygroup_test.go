package securitygroup_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

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

func TestPatchAppendNewIngressRule(t *testing.T) {
	testCases := []struct {
		name string

		withExistingTemplate []byte
		withRule             api.Rule
	}{
		{
			name: "Should correctly add ingress rule",

			withExistingTemplate: []byte(aTemplate),
			withRule: api.Rule{
				Description: "Allow traffic from application security group",
				FromPort:    5432,
				ToPort:      5432,
				CidrIP:      "192.168.0.1/24",
				Protocol:    "tcp",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			updatedTemplate, err := securitygroup.PatchAppendIngressRule(tc.withExistingTemplate, aResourceName, tc.withRule)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, updatedTemplate)
		})
	}
}

func TestPatchAppendNewEgressRule(t *testing.T) {
	testCases := []struct {
		name string

		withExistingTemplate []byte
		withRule             api.Rule
	}{
		{
			name: "Should correctly add egress rule",

			withExistingTemplate: []byte(aTemplate),
			withRule: api.Rule{
				Description: "Allow traffic to application security group",
				FromPort:    5432,
				ToPort:      5432,
				CidrIP:      "192.168.0.1/24",
				Protocol:    "tcp",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			updatedTemplate, err := securitygroup.PatchAppendEgressRule(tc.withExistingTemplate, aResourceName, tc.withRule)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, updatedTemplate)
		})
	}
}

func TestPatchRemoveIngressRule(t *testing.T) {
	testCases := []struct {
		name string

		withTemplate []byte
		withRule     api.Rule
	}{
		{
			name: "Should correctly remove ingress rule",

			withTemplate: []byte(aTemplate),
			withRule: api.Rule{
				FromPort: 53,
				ToPort:   53,
				CidrIP:   "192.168.0.0/20",
				Protocol: "tcp",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			updatedTemplate, err := securitygroup.PatchRemoveIngressRule(tc.withTemplate, aResourceName, tc.withRule)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, updatedTemplate)
		})
	}
}

func TestPatchRemoveEgressRule(t *testing.T) {
	testCases := []struct {
		name string

		withTemplate []byte
		withRule     api.Rule
	}{
		{
			name: "Should correctly remove egress rule",

			withTemplate: []byte(aTemplate),
			withRule: api.Rule{
				FromPort: 5432,
				ToPort:   5432,
				CidrIP:   "192.168.5.0/24",
				Protocol: "tcp",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			updatedTemplate, err := securitygroup.PatchRemoveEgressRule(tc.withTemplate, aResourceName, tc.withRule)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, updatedTemplate)
		})
	}
}

func TestPatchAppendEgressRuleIdempotency(t *testing.T) {
	testCases := []struct {
		name string

		withTemplate []byte
		withRules    []api.Rule
	}{
		{
			name: "Should only append the same egress rule once",

			withTemplate: []byte(aTemplate),
			withRules: []api.Rule{
				{
					Description: "Fancy and enlightning description",
					FromPort:    1337,
					ToPort:      1337,
					CidrIP:      "10.0.0.0/24",
					Protocol:    "tcp",
				},
				{
					Description: "Boring, but enlightning description",
					FromPort:    1337,
					ToPort:      1337,
					CidrIP:      "10.0.0.0/24",
					Protocol:    "tcp",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			var (
				updatedTemplate = tc.withTemplate
				err             error
			)

			for _, rule := range tc.withRules {
				updatedTemplate, err = securitygroup.PatchAppendEgressRule(updatedTemplate, aResourceName, rule)
				assert.NoError(t, err)
			}

			g := goldie.New(t)
			g.Assert(t, tc.name, updatedTemplate)
		})
	}
}

func TestPatchAppendIngressRuleIdempotency(t *testing.T) {
	testCases := []struct {
		name string

		withTemplate []byte
		withRules    []api.Rule
	}{
		{
			name: "Should only append the same ingress rule once",

			withTemplate: []byte(aTemplate),
			withRules: []api.Rule{
				{
					Description: "Fancy and enlightning description",
					FromPort:    1337,
					ToPort:      1337,
					CidrIP:      "10.0.0.0/24",
					Protocol:    "tcp",
				},
				{
					Description: "Boring, but enlightning description",
					FromPort:    1337,
					ToPort:      1337,
					CidrIP:      "10.0.0.0/24",
					Protocol:    "tcp",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			var (
				updatedTemplate = tc.withTemplate
				err             error
			)

			for _, rule := range tc.withRules {
				updatedTemplate, err = securitygroup.PatchAppendIngressRule(updatedTemplate, aResourceName, rule)
				assert.NoError(t, err)
			}

			g := goldie.New(t)
			g.Assert(t, tc.name, updatedTemplate)
		})
	}
}

const (
	aResourceName = "SG"
	aTemplate     = `AWSTemplateFormatVersion: 2010-09-09
Outputs:
  SG:
    Export:
      Name:
        Fn::Sub: ${AWS::StackName}-SG
    Value:
      Ref: SG
  SGGroupId:
    Export:
      Name:
        Fn::Sub: ${AWS::StackName}-SGGroupId
    Value:
      Fn::GetAtt:
      - SG
      - GroupId
Resources:
  SG:
    Properties:
      GroupDescription: Defines network access regarding dbapp
      GroupName: dbappjulius-devdbapp
      SecurityGroupEgress:
      - CidrIp: 192.168.8.0/24
        Description: Allow Postgres traffic to database subnet
        FromPort: 5432
        IpProtocol: tcp
        ToPort: 5432
      - CidrIp: 192.168.5.0/24
        Description: Allow Postgres traffic to database subnet
        FromPort: 5432
        IpProtocol: tcp
        ToPort: 5432
      - CidrIp: 192.168.2.0/24
        Description: Allow Postgres traffic to database subnet
        FromPort: 5432
        IpProtocol: tcp
        ToPort: 5432
      SecurityGroupIngress:
      - CidrIp: 192.168.0.0/20
        Description: Required DNS/tcp entrypoint for control plane
        FromPort: 53
        IpProtocol: tcp
        ToPort: 53
      - CidrIp: 192.168.0.0/20
        Description: Required DNS/udp entrypoint for control plane
        FromPort: 53
        IpProtocol: udp
        ToPort: 53
      VpcId: vpc-074eaaf3bcc642368
    Type: AWS::EC2::SecurityGroup`
)
