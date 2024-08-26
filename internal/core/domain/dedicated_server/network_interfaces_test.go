package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewNetworkInterfaces(t *testing.T) {

	public := NetworkInterface{Mac: "public"}
	internal := NetworkInterface{Mac: "internal"}
	remote := NetworkInterface{Mac: "remote"}
	networkInterfaces := NewNetworkInterfaces(public, internal, remote)
	assert.Equal(t, "public", networkInterfaces.Public.Mac)
	assert.Equal(t, "internal", networkInterfaces.Internal.Mac)
	assert.Equal(t, "remote", networkInterfaces.RemoteManagement.Mac)
}
