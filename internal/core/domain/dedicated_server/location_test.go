package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLocation(t *testing.T) {

	got := NewLocation("rack", "site", "suite", "unit")
	want := Location{
		Rack:  "rack",
		Site:  "site",
		Suite: "suite",
		Unit:  "unit",
	}
	assert.Equal(t, want, got)
}
