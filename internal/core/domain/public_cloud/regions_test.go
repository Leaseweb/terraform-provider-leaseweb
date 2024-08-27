package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegions_Contains(t *testing.T) {
	t.Run("contains find result", func(t *testing.T) {
		regions := Regions{{Name: "region"}}

		assert.True(t, regions.Contains("region"))
	})

	t.Run("contains does not result", func(t *testing.T) {
		regions := Regions{{Name: "region"}}

		assert.False(t, regions.Contains("tralala"))
	})

}

func TestRegions_ToArray(t *testing.T) {
	regions := Regions{{Name: "region"}}
	got := regions.ToArray()
	want := []string{"region"}

	assert.Equal(t, want, got)
}

func TestRegions_GetByName(t *testing.T) {
	t.Run("returns error when region cannot be found", func(t *testing.T) {
		regions := Regions{{Name: "region"}}
		got, err := regions.GetByName("nonexistent")

		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("returns region when region can be found", func(t *testing.T) {
		want := Region{Name: "region"}
		regions := Regions{want}
		got, err := regions.GetByName("region")

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})
}
