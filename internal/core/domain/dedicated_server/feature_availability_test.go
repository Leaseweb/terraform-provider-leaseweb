package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewFeatureAvailability(t *testing.T) {

	got := NewFeatureAvailability(true, false, true, false, true)
	want := FeatureAvailability{
		Automation:       true,
		IpmiReboot:       false,
		PowerCycle:       true,
		PrivateNetwork:   false,
		RemoteManagement: true,
	}
	assert.Equal(t, want, got)

	got = NewFeatureAvailability(false, true, false, true, false)
	want = FeatureAvailability{
		Automation:       false,
		IpmiReboot:       true,
		PowerCycle:       false,
		PrivateNetwork:   true,
		RemoteManagement: false,
	}
	assert.Equal(t, want, got)
}
