package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	image := NewImage(
		"UBUNTU_24_04_64BIT",
		"name",
		"version",
		"family",
		"flavour",
		[]string{"marketApp"},
		[]string{"storageType"})

	assert.Equal(t, "UBUNTU_24_04_64BIT", image.Id)
	assert.Equal(t, "name", image.Name)
	assert.Equal(t, "version", image.Version)
	assert.Equal(t, "family", image.Family)
	assert.Equal(t, "flavour", image.Flavour)
	assert.Equal(t, []string{"marketApp"}, image.MarketApps)
	assert.Equal(t, []string{"storageType"}, image.StorageTypes)
}
