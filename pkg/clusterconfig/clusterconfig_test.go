package clusterconfig_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewClusterConfig(t *testing.T) {
	cfg := clusterconfig.New("test").
		PermissionsBoundary(v1alpha1.PermissionsBoundaryARN(mock.DefaultAWSAccountID)).
		Region(mock.DefaultRegion).
		Vpc(mock.DefaultVpcID, mock.DefaultCidr).
		Subnets(mock.DefaultVpcPublicSubnets(), mock.DefaultVpcPrivateSubnets()).
		Build()
	got, err := cfg.YAML()
	assert.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "clusterConfig", got)
}
