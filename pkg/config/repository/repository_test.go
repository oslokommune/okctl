package repository_test

import (
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/config/repository"
)

// nolint: funlen
func TestData(t *testing.T) {
	testCases := []struct {
		name   string
		data   *repository.Data
		golden string
	}{
		{
			name: "Should work",
			data: &repository.Data{
				Name:      "okctl",
				Region:    "eu-west-1",
				OutputDir: "infrastructure",
				Clusters: map[string]*repository.Cluster{
					"pro": {
						Name:         "okctl-pro",
						Environment:  "pro",
						AWSAccountID: "123456789012",
						HostedZone: map[string]*repository.HostedZone{
							"test.oslo.systems": {
								IsDelegated: true,
								IsCreated:   false,
								Domain:      "test.oslo.systems",
								FQDN:        "test.oslo.systems",
								NameServers: []string{
									"ns1.aws.com",
									"ns2.aws.com",
								},
							},
						},
						VPC: &repository.VPC{
							VpcID: "3456ygfghj",
							CIDR:  "192.168.0.0/20",
						},
						Certificates: map[string]string{
							"argocd.test.oslo.systems": "arn:::cert/something",
						},
						Github: &repository.Github{
							Organisation: "oslokommune",
							OauthApp: map[string]*repository.OauthApp{
								"okctl-kjøremlijø-pro": {
									Team:     "kjøremiljø",
									Name:     "okctl-kjøremiljø-pro",
									ClientID: "asdfg123456",
									ClientSecret: &repository.ClientSecret{
										Name:    "argocd-client-secret",
										Path:    "/something/argocd",
										Version: 1,
									},
								},
							},
							Repositories: map[string]*repository.Repository{
								"oslokommune/okctl-iac": {
									Name:   "okctl-iac",
									Types:  []string{"infrastructure"},
									GitURL: "git@github.com/oslokommune/okctl-iac",
									DeployKey: &repository.DeployKey{
										Title:     "okctl-kjøremlijø-pro",
										ID:        23456865,
										PublicKey: "ssh-rsa 098f09ujf9rewjvjlejf3jf933",
										PrivateKeySecret: &repository.PrivateKeySecret{
											Name:    "okctl-kjøremiljø-pro",
											Path:    "/something/privatekey",
											Version: 1,
										},
									},
								},
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
			got, err := tc.data.YAML()
			assert.NoError(t, err)
			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
