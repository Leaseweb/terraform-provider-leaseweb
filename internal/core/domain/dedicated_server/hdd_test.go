package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewHdd(t *testing.T) {

	got := NewHdd("id", "type", "unit", "performanceType", 12, 34)
	want := Hdd{
		Id:              "id",
		Type:            "type",
		Unit:            "unit",
		PerformanceType: "performanceType",
		Amount:          12,
		Size:            34,
	}
	assert.Equal(t, want, got)
}
