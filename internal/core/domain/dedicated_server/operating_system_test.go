package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	got := NewOperatingSystem("id", "name")
	want := OperatingSystem{Id: "id", Name: "name"}

	assert.Equal(t, want, got)
}
