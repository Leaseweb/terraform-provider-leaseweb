package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegion(t *testing.T) {
	want := Region{Name: "name", Location: "location"}
	got := NewRegion("name", "location")

	assert.Equal(t, want, got)

}

func TestRegion_String(t *testing.T) {
	want := "name"
	got := NewRegion("name", "location").String()

	assert.Equal(t, want, got)
}
