package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewNetworkInterface(t *testing.T) {
	got := NewNetworkInterface("mac", "ip", "gateway", "loc_id", true, Ports{Port{Name: "name"}})
	want := NetworkInterface{
		Mac:        "mac",
		Ip:         "ip",
		Gateway:    "gateway",
		LocationId: "loc_id",
		NullRouted: true,
		Ports:      Ports{Port{Name: "name"}},
	}
	assert.Equal(t, want, got)
}
