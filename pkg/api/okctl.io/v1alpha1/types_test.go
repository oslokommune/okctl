package v1alpha1_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewClusterConfig(t *testing.T) {
	cfg := v1alpha1.NewClusterConfig()
	got, err := cfg.YAML()
	assert.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "clusterConfig", got)
}
