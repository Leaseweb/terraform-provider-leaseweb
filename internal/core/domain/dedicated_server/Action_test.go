package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewAction(t *testing.T) {
	got := NewAction(
		"EMAIL",
		OptionalActionValues{},
	)
	want := Action{
		Type: "EMAIL",
	}
	assert.Equal(t, want, got)
}
