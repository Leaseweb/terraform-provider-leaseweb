package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPrivateNetwork(t *testing.T) {
	privateNetwork := NewPrivateNetwork("id", "status", "subnet")

	assert.Equal(t, "id", privateNetwork.Id)
	assert.Equal(t, "status", privateNetwork.Status)
	assert.Equal(t, "subnet", privateNetwork.Subnet)
}
