package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceTypes_Contains(t *testing.T) {
	t.Run("return true if instanceType exists", func(t *testing.T) {
		instanceTypes := InstanceTypes{"tralala"}
		assert.True(t, instanceTypes.Contains("tralala"))
	})

	t.Run("return false if instanceType does not exist", func(t *testing.T) {
		instanceTypes := InstanceTypes{"piet"}
		assert.False(t, instanceTypes.Contains("tralala"))
	})
}
