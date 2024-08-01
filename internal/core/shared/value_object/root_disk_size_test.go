package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRootDiskSize(t *testing.T) {
	t.Run("values are set", func(t *testing.T) {
		got, err := NewRootDiskSize(10)

		assert.NoError(t, err)
		assert.Equal(t, 10, got.Value)
	})

	t.Run("return error when rootDiskSize is too small", func(t *testing.T) {
		_, err := NewRootDiskSize(4)

		assert.Error(t, err)
	})

	t.Run("return error when rootDiskSize is too large", func(t *testing.T) {
		_, err := NewRootDiskSize(1001)

		assert.Error(t, err)
	})
}
