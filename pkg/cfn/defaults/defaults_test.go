package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVPC(t *testing.T) {
	b, err := VPC("test", "test", "192.168.0.0/20", "eu-west-1")
	assert.NoError(t, err)
	assert.NotNil(t, b)
}
