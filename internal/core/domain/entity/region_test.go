package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegion(t *testing.T) {
	want := Region{Name: "name", Location: "location"}
	got := NewRegion("name", "location")

	assert.Equal(t, want, got)

}
