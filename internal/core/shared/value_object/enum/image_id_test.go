package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageId_String(t *testing.T) {
	got := Almalinux864Bit.String()

	assert.Equal(t, "ALMALINUX_8_64BIT", got)
}

func TestNewImageId(t *testing.T) {
	want := Debian1164Bit
	got, err := NewImageId("DEBIAN_11_64BIT")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
