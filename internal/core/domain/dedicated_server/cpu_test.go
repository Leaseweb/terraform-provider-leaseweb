package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCpu(t *testing.T) {

	got := NewCpu(12, "type")
	want := Cpu{
		Quantity: 12,
		Type:     "type",
	}
	assert.Equal(t, want, got)
}
