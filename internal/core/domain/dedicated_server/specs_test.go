package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSpecs(t *testing.T) {
	got := NewSpecs(
		"chassis",
		true,
		Cpu{Quantity: 1},
		Ram{Size: 2},
		Hdds{Hdd{Id: "hid"}},
		PciCards{PciCard{"d"}},
	)
	want := Specs{
		Chassis:             "chassis",
		HardwareRaidCapable: true,
		Cpu:                 Cpu{Quantity: 1},
		Ram:                 Ram{Size: 2},
		Hdds:                Hdds{Hdd{Id: "hid"}},
		PciCards:            PciCards{PciCard{"d"}},
	}
	assert.Equal(t, want, got)
}
