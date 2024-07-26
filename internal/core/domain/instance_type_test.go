package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/enum"
)

func TestNewInstanceType(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		resources := Resources{Cpu: Cpu{Unit: "unit"}}
		prices := Prices{Compute: Price{HourlyPrice: "11"}}

		optional := OptionalInstanceTypeValues{}

		got := NewInstanceType("name", resources, prices, optional)

		assert.Equal(t, "name", got.Name)
		assert.Equal(t, resources, got.Resources)
		assert.Equal(t, prices, got.Prices)

		assert.Nil(t, got.StorageTypes)
	})

	t.Run("optional values are set", func(t *testing.T) {
		storageTypes := StorageTypes{enum.RootDiskStorageTypeCentral}
		optional := OptionalInstanceTypeValues{StorageTypes: &storageTypes}

		got := NewInstanceType("name", Resources{}, Prices{}, optional)

		assert.Equal(t, storageTypes, *got.StorageTypes)
	})
}

func TestInstanceType_String(t *testing.T) {
	instanceType := InstanceType{Name: "tralala"}
	got := instanceType.String()
	want := "tralala"

	assert.Equal(t, want, got)
}
