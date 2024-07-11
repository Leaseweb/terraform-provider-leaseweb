package value_object

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUuid(t *testing.T) {
	t.Run("valid uuid is accepted", func(t *testing.T) {
		want := "8df3c996-dcd3-45bd-8bbc-90a90f36137f"
		got, err := NewUuid(want)

		assert.NoError(t, err)
		assert.Equal(t, want, got.Uuid.String())
	})

	t.Run("valid uuid is not accepted", func(t *testing.T) {
		_, err := NewUuid("tralala")

		assert.ErrorIs(t, ErrCouldNotConvertValueIntoUUID, err)
	})
}

func TestUuid_String(t *testing.T) {
	want := "8df3c996-dcd3-45bd-8bbc-90a90f36137f"
	got, _ := NewUuid(want)

	assert.Equal(t, want, got.String())
}

func TestNewGeneratedUuid(t *testing.T) {
	got := NewGeneratedUuid()

	assert.NoError(t, uuid.Validate(got.String()))

}
