package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPrivateNetwork(t *testing.T) {
	got := NewPrivateNetwork("id", "status", "subnet", "vlanid", 12)
	want := PrivateNetwork{
		Id:        "id",
		Status:    "status",
		Subnet:    "subnet",
		VlanId:    "vlanid",
		LinkSpeed: 12,
	}
	assert.Equal(t, want, got)
}
