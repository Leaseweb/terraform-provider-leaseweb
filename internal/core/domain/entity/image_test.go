package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func TestNewImage(t *testing.T) {
	image := NewImage(
		enum.UBUNTU_24_04_64_BIT,
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"marketApp"},
		[]string{"storageType"})

	assert.Equal(t, enum.UBUNTU_24_04_64_BIT, image.Id)
	assert.Equal(t, "name", image.Name)
	assert.Equal(t, "version", image.Version)
	assert.Equal(t, "family", image.Family)
	assert.Equal(t, "flavour", image.Flavour)
	assert.Equal(t, "architecture", image.Architecture)
	assert.Equal(t, []string{"marketApp"}, image.MarketApps)
	assert.Equal(t, []string{"storageType"}, image.StorageTypes)
}
