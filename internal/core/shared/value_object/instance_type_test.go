package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInstanceType(t *testing.T) {
	t.Run("valid instanceType is accepted", func(t *testing.T) {
		got, err := NewInstanceType("tralala", []string{"tralala"})
		want := InstanceType{Type: "tralala"}

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("invalid instanceType returns an error", func(t *testing.T) {
		got, err := NewInstanceType("tralala", []string{"blaat"})

		assert.Nil(t, got)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func TestInstanceType_String(t *testing.T) {
	instanceType := InstanceType{Type: "tralala"}
	got := instanceType.String()
	want := "tralala"

	assert.Equal(t, want, got)
}

func TestNewUnvalidatedInstanceType(t *testing.T) {
	got := NewUnvalidatedInstanceType("tralala")
	want := InstanceType{Type: "tralala"}

	assert.Equal(t, want, got)
}
