package dedicated_server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewControlPanel(t *testing.T) {
	got := NewControlPanel("id", "name")
	want := ControlPanel{Id: "id", Name: "name"}

	assert.Equal(t, want, got)
}
