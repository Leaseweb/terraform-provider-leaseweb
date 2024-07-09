package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageId_String(t *testing.T) {
	got := Almalinux864Bit.String()

	assert.Equal(t, "ALMALINUX_8_64BIT", got)
}
