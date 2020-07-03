package api_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewClusterConfig(t *testing.T) {
	cfg := api.NewClusterConfig()
	got, err := cfg.YAML()
	assert.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "clusterConfig", got)
}
