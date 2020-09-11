package state_test

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/config/state"
)

// nolint: funlen
func TestData(t *testing.T) {
	testCases := []struct {
		name   string
		data   *state.Repository
		golden string
	}{
		{
			name: "Should work",
			data: &state.Repository{
				Metadata: &state.Metadata{
					Name:      "okctl",
					Region:    "eu-west-1",
					OutputDir: "infrastructure",
				},
				Clusters: map[string]*state.Cluster{
					"pro": {
						Name:         "okctl-pro",
						Environment:  "pro",
						AWSAccountID: "123456789012",
						HostedZone: map[string]*state.HostedZone{
							"test.oslo.systems": {
								IsDelegated: true,
								Domain:      "test.oslo.systems",
								FQDN:        "test.oslo.systems",
								NameServers: []string{
									"ns1.aws.com",
									"ns2.aws.com",
								},
							},
						},
						VPC: &state.VPC{
							Subnets: map[string][]*state.VPCSubnet{
								state.SubnetTypePublic: {
									{
										CIDR:             "192.168.0.0/24",
										AvailabilityZone: "eu-west-1a",
									},
								},
								state.SubnetTypePrivate: {
									{
										CIDR:             "192.168.10.0/24",
										AvailabilityZone: "eu-west-1c",
									},
								},
							},
							VpcID: "3456ygfghj",
							CIDR:  "192.168.0.0/20",
						},
						Certificates: map[string]*state.Certificate{
							"argocd.test.oslo.systems": {

								Domain: "argocd.test.oslo.systems",
								ARN:    "arn:::cert/something",
							},
						},
						Github: &state.Github{
							Organisation: "oslokommune",
							OauthApp: map[string]*state.GithubOauthApp{
								"okctl-kjøremlijø-pro": {
									Team:     "kjøremiljø",
									Name:     "okctl-kjøremiljø-pro",
									ClientID: "asdfg123456",
									ClientSecret: &state.ClientSecret{
										Name:    "argocd-client-secret",
										Path:    "/something/argocd",
										Version: 1,
									},
								},
							},
							Repositories: map[string]*state.GithubRepository{
								"oslokommune/okctl-iac": {
									Name:   "okctl-iac",
									Types:  []string{"infrastructure"},
									GitURL: "git@github.com/oslokommune/okctl-iac",
									DeployKey: &state.DeployKey{
										Title:     "okctl-kjøremlijø-pro",
										ID:        23456865,
										PublicKey: "ssh-rsa 098f09ujf9rewjvjlejf3jf933",
										PrivateKeySecret: &state.PrivateKeySecret{
											Name:    "okctl-kjøremiljø-pro",
											Path:    "/something/privatekey",
											Version: 1,
										},
									},
								},
							},
						},
						ArgoCD: &state.ArgoCD{
							SiteURL: "https://argocd.oslo.systems",
							Domain:  "argocd.oslo.systems",
							SecretKey: &state.SecretKeySecret{
								Name:    "something",
								Path:    "/some/path",
								Version: 1,
							},
						},
					},
				},
			},
			golden: "basic",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := yaml.Marshal(tc.data)
			assert.NoError(t, err)
			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
