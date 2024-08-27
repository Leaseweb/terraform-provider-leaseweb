package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRam(t *testing.T) {
	got := NewRam(12, "gb")
	want := Ram{
		Size: 12,
		Unit: "gb",
	}
	assert.Equal(t, want, got)
}
