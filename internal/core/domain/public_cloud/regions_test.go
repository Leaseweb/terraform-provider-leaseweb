package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegions_Contains(t *testing.T) {
	t.Run("contains find result", func(t *testing.T) {
		regions := Regions{"region"}

		assert.True(t, regions.Contains("region"))
	})

	t.Run("contains does not find result", func(t *testing.T) {
		regions := Regions{"region"}

		assert.False(t, regions.Contains("tralala"))
	})

}
