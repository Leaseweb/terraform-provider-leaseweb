package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceTypes_ToArray(t *testing.T) {
	instanceTypes := InstanceTypes{InstanceType{Name: "tralala"}}
	got := instanceTypes.ToArray()
	want := []string{"tralala"}

	assert.Equal(t, want, got)
}

func TestInstanceTypes_ContainsName(t *testing.T) {
	t.Run("return true if name exists", func(t *testing.T) {
		instanceTypes := InstanceTypes{InstanceType{Name: "tralala"}}
		assert.True(t, instanceTypes.ContainsName("tralala"))
	})

	t.Run("return false if name does not exist", func(t *testing.T) {
		instanceTypes := InstanceTypes{InstanceType{Name: "piet"}}
		assert.False(t, instanceTypes.ContainsName("tralala"))
	})
}
