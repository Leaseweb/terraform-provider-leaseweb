package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	got := NewImage("UBUNTU_24_04_64BIT")
	want := Image{Id: "UBUNTU_24_04_64BIT"}

	assert.Equal(t, want, got)
}
