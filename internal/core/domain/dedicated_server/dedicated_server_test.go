package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewDedicatedServer(t *testing.T) {

	server := NewDedicatedServer(
		"id",
		"assetId",
		"sn",
		Rack{Id: "rid"},
		Location{Rack: "rack"},
		FeatureAvailability{Automation: true},
		Contract{Id: "cid"},
		Ports{Port{Name: "name1"}},
		PrivateNetworks{
			PrivateNetwork{Id: "pid"},
		},
		NetworkInterfaces{
			Public: NetworkInterface{Mac: "mac"},
		},
		Specs{
			Chassis: "chassis",
		},
	)

	assert.Equal(t, "id", server.Id)
	assert.Equal(t, "assetId", server.AssetId)
	assert.Equal(t, "sn", server.SerialNumber)
	assert.Equal(t, "chassis", server.Specs.Chassis)
	assert.Equal(t, "cid", server.Contract.Id)
	assert.Equal(t, "rid", server.Rack.Id)
	assert.True(t, server.FeatureAvailability.Automation)
	assert.Equal(t, "rack", server.Location.Rack)
	assert.Equal(t, "name1", server.PowerPorts[0].Name)
	assert.Equal(t, "pid", server.PrivateNetworks[0].Id)
	assert.Equal(t, "mac", server.NetworkInterfaces.Public.Mac)
}
