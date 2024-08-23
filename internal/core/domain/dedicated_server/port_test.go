package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPort(t *testing.T) {
	got := NewPort("name", "1111")
	want := Port{
		Name: "name",
		Port: "1111",
	}
	assert.Equal(t, want, got)
}
