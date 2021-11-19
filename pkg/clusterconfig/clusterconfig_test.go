package clusterconfig_test

import (
	"context"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/version/developmentversion"
	"testing"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name      string
		args      *clusterconfig.Args
		golden    string
		expectErr bool
		err       string
	}{
		{
			name: "Validate cluster config",
			args: &clusterconfig.Args{
				ClusterName:            "test",
				PermissionsBoundaryARN: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
				PrivateSubnets:         mock.DefaultVpcPrivateSubnets(),
				PublicSubnets:          mock.DefaultVpcPublicSubnets(),
				Region:                 mock.DefaultRegion,
				Version:                mock.DefaultVersion,
				VpcCidr:                mock.DefaultCidr,
				VpcID:                  mock.DefaultVpcID,
			},
			golden: "clusterConfig",
		},
		{
			name:      "Invalid cluster config",
			args:      &clusterconfig.Args{},
			expectErr: true,
			err:       "ClusterName: cannot be blank; PermissionsBoundaryARN: cannot be blank; PrivateSubnets: cannot be blank; PublicSubnets: cannot be blank; Region: cannot be blank; Version: cannot be blank; VpcCidr: cannot be blank; VpcID: cannot be blank.", // nolint: lll
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.New(tc.args)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func TestNewServiceAccount(t *testing.T) {
	testCases := []struct {
		name      string
		args      *clusterconfig.ServiceAccountArgs
		golden    string
		expectErr bool
		err       string
	}{
		{
			name: "Validate service account",
			args: &clusterconfig.ServiceAccountArgs{
				ClusterName: "test",
				Labels: map[string]string{
					"label": "some-label",
				},
				Name:                   "someName",
				Namespace:              "kube-system",
				PermissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
				PolicyArn:              mock.DefaultPolicyARN,
				Region:                 mock.DefaultRegion,
			},
			golden: "serviceAccount",
		},
		{
			name:      "Invalid service account",
			args:      &clusterconfig.ServiceAccountArgs{},
			expectErr: true,
			err:       "ClusterName: cannot be blank; Labels: cannot be blank; Name: cannot be blank; Namespace: cannot be blank; PermissionsBoundaryArn: cannot be blank; PolicyArn: cannot be blank; Region: cannot be blank.", // nolint: lll
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.NewServiceAccount(tc.args)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func TestNewExternalSecretsServiceAccount(t *testing.T) {
	testCases := []struct {
		name                   string
		clusterName            string
		region                 string
		policyArn              string
		permissionsBoundaryArn string
		golden                 string
		expectErr              bool
		err                    string
	}{
		{
			name:                   "Validate service account",
			clusterName:            "test",
			region:                 mock.DefaultRegion,
			policyArn:              mock.DefaultPolicyARN,
			permissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
			golden:                 "externalSecretsServiceAccount",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.NewExternalSecretsServiceAccount(tc.clusterName, tc.region, tc.policyArn, tc.permissionsBoundaryArn)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func TestNewAlbIngressControllerServiceAccount(t *testing.T) {
	testCases := []struct {
		name                   string
		clusterName            string
		region                 string
		policyArn              string
		permissionsBoundaryArn string
		golden                 string
		expectErr              bool
		err                    string
	}{
		{
			name:                   "Validate service account",
			clusterName:            "test",
			region:                 mock.DefaultRegion,
			policyArn:              mock.DefaultPolicyARN,
			permissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
			golden:                 "aws-alb-ingres-controller",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.NewAlbIngressControllerServiceAccount(tc.clusterName, tc.region, tc.policyArn, tc.permissionsBoundaryArn)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func TestNewAWSLoadBalancerControllerServiceAccount(t *testing.T) {
	testCases := []struct {
		name                   string
		clusterName            string
		region                 string
		policyArn              string
		permissionsBoundaryArn string
		golden                 string
		expectErr              bool
		err                    string
	}{
		{
			name:                   "Validate service account",
			clusterName:            "test",
			region:                 mock.DefaultRegion,
			policyArn:              mock.DefaultPolicyARN,
			permissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
			golden:                 "aws-load-balancer-controller",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.NewAWSLoadBalancerControllerServiceAccount(tc.clusterName, tc.region, tc.policyArn, tc.permissionsBoundaryArn)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func TestNewExternalDNSServiceAccount(t *testing.T) {
	testCases := []struct {
		name                   string
		clusterName            string
		region                 string
		policyArn              string
		permissionsBoundaryArn string
		golden                 string
		expectErr              bool
		err                    string
	}{
		{
			name:                   "Validate service account",
			clusterName:            "test",
			region:                 mock.DefaultRegion,
			policyArn:              mock.DefaultPolicyARN,
			permissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
			golden:                 "external-dns",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.NewExternalDNSServiceAccount(tc.clusterName, tc.region, tc.policyArn, tc.permissionsBoundaryArn)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func TestNewAutoscalerServiceAccount(t *testing.T) {
	testCases := []struct {
		name                   string
		clusterName            string
		region                 string
		policyArn              string
		permissionsBoundaryArn string
		golden                 string
		expectErr              bool
		err                    string
	}{
		{
			name:                   "Validate service account",
			clusterName:            "test",
			region:                 mock.DefaultRegion,
			policyArn:              mock.DefaultPolicyARN,
			permissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
			golden:                 "autoscaler",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.NewAutoscalerServiceAccount(tc.clusterName, tc.region, tc.policyArn, tc.permissionsBoundaryArn)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func TestNewBlockstorageServiceAccount(t *testing.T) {
	testCases := []struct {
		name                   string
		clusterName            string
		region                 string
		policyArn              string
		permissionsBoundaryArn string
		golden                 string
		expectErr              bool
		err                    string
	}{
		{
			name:                   "Validate service account",
			clusterName:            "test",
			region:                 mock.DefaultRegion,
			policyArn:              mock.DefaultPolicyARN,
			permissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
			golden:                 "blockstorage",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.NewBlockstorageServiceAccount(tc.clusterName, tc.region, tc.policyArn, tc.permissionsBoundaryArn)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func TestNewCloudwatchDatasourceServiceAccount(t *testing.T) {
	testCases := []struct {
		name                   string
		clusterName            string
		region                 string
		policyArn              string
		namespace              string
		permissionsBoundaryArn string
		golden                 string
		expectErr              bool
		err                    string
	}{
		{
			name:                   "Validate service account",
			clusterName:            "test",
			region:                 mock.DefaultRegion,
			policyArn:              mock.DefaultPolicyARN,
			namespace:              constant.DefaultMonitoringNamespace,
			permissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID),
			golden:                 "cloudwatch-datasource",
			expectErr:              false,
			err:                    "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			mockVersion()
			defer developmentversion.Reset()

			got, err := clusterconfig.NewCloudwatchDatasourceServiceAccount(tc.clusterName, tc.region, tc.policyArn, tc.namespace, tc.permissionsBoundaryArn)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			} else {
				assert.NoError(t, err)

				d, err := got.YAML()
				assert.NoError(t, err)

				g := goldie.New(t)
				g.Assert(t, tc.golden, d)
			}
		})
	}
}

func mockVersion() {
	developmentversion.ListReleases = func(
		ctx context.Context, owner string, repo string,
	) ([]*developmentversion.RepositoryRelease, error) {
		return []*developmentversion.RepositoryRelease{
			{
				TagName: github.StringPtr("0.0.70"),
			},
		}, nil
	}
}
