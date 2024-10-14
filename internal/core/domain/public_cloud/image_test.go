package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	got := NewImage(
		"UBUNTU_24_04_64BIT",
		"name",
		"family",
		"flavour",
		false,
	)
	want := Image{
		Id:      "UBUNTU_24_04_64BIT",
		Name:    "name",
		Family:  "family",
		Flavour: "flavour",
		Custom:  false,
	}

	assert.Equal(t, want, got)
}
