package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/enum"
)

func TestNewImage(t *testing.T) {
	image := NewImage(
		enum.Ubuntu240464Bit,
		"name",
		"version",
		"family",
		"flavour",
		[]string{"marketApp"},
		[]string{"storageType"})

	assert.Equal(t, enum.Ubuntu240464Bit, image.Id)
	assert.Equal(t, "name", image.Name)
	assert.Equal(t, "version", image.Version)
	assert.Equal(t, "family", image.Family)
	assert.Equal(t, "flavour", image.Flavour)
	assert.Equal(t, []string{"marketApp"}, image.MarketApps)
	assert.Equal(t, []string{"storageType"}, image.StorageTypes)
}
