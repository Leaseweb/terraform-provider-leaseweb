package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImages_FilterById(t *testing.T) {
	t.Run("image is found", func(t *testing.T) {
		images := Images{Image{Id: "tralala"}}

		got, err := images.FilterById("tralala")
		want := Image{Id: "tralala"}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("image is not found", func(t *testing.T) {
		images := Images{Image{Id: "tralala"}}
		_, err := images.FilterById("blaat")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "blaat")
	})
}
