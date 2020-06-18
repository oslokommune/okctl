package vpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVPC(t *testing.T) {
	err := New("test", "test", "192.168.0.0/20", "eu-west-1").Build()
	assert.NoError(t, err)
}
