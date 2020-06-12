package vpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVPC(t *testing.T) {
	got, err := New("test", "test", "192.168.0.0/20", "eu-west-1").Build()
	assert.NotNil(t, got)
	assert.NoError(t, err)
}
