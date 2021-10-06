package kubeconfig_test

import (
	"encoding/base64"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"

	mock2 "github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/kubeconfig"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func testConfig(t *testing.T) *v1alpha5.ClusterConfig {
	conf, err := clusterconfig.New(&clusterconfig.Args{
		ClusterVersionInfo:     mock.DefaultVersionInfo(),
		ClusterName:            "okctl-pro",
		PermissionsBoundaryARN: v1alpha1.PrincipalARN(mock.DefaultAWSAccountID),
		PrivateSubnets:         mock.DefaultVpcPrivateSubnets(),
		PublicSubnets:          mock.DefaultVpcPublicSubnets(),
		Region:                 "eu-west-1",
		Version:                mock.DefaultVersion,
		VpcCidr:                mock.DefaultCidr,
		VpcID:                  mock.DefaultVpcID,
	})
	assert.NoError(t, err)

	conf.Status = &v1alpha5.ClusterStatus{
		Endpoint:                 "https://some-endpoint",
		CertificateAuthorityData: []byte(base64.StdEncoding.EncodeToString([]byte("something"))),
		ARN:                      "arn:///something",
		StackName:                mock.DefaultClusterName,
	}

	return conf
}

func TestCreate(t *testing.T) {
	testCases := []struct {
		name   string
		cfg    *v1alpha5.ClusterConfig
		golden string
	}{
		{
			name:   "Should work",
			cfg:    testConfig(t),
			golden: "kubeconf",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := clientcmd.Write(kubeconfig.Create("someuser", tc.cfg))
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		name      string
		cfg       *v1alpha5.ClusterConfig
		provider  v1alpha1.CloudProvider
		expect    interface{}
		expectErr bool
		golden    string
	}{
		{
			name:     "Should work",
			cfg:      testConfig(t),
			provider: mock2.NewGoodCloudProvider(),
			golden:   "valid-config",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := kubeconfig.New(tc.cfg, tc.provider).Get()
			if tc.expectErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				data, err := got.Bytes()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, data)
			}
		})
	}
}
