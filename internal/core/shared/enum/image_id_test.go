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

func TestImageId_Values(t *testing.T) {
	want := []string{
		Almalinux864Bit.String(),
		Almalinux964Bit.String(),
		ArchLinux64Bit.String(),
		Centos764Bit.String(),
		Debian1264Bit.String(),
		Debian1064Bit.String(),
		Debian1164Bit.String(),
		Debian1264Bit.String(),
		Freebsd1364Bit.String(),
		Freebsd1464Bit.String(),
		RockyLinux864Bit.String(),
		RockyLinux964Bit.String(),
		Ubuntu200464Bit.String(),
		Ubuntu220464Bit.String(),
		Ubuntu240464Bit.String(),
		WindowsServer2016Standard64Bit.String(),
		WindowsServer2019Standard64Bit.String(),
		WindowsServer2022Standard64Bit.String(),
	}
	got := Debian1264Bit.Values()

	assert.Equal(t, want, got)
}
