package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRack(t *testing.T) {
	got := NewRack("id", "cap", "type")
	want := Rack{
		Id:       "id",
		Capacity: "cap",
		Type:     "type",
	}
	assert.Equal(t, want, got)
}
