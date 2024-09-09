package to_domain_entity

import (
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdaptNullableStringToValue(t *testing.T) {
	t.Run("value is returned if set", func(t *testing.T) {
		val := "value"

		got := AdaptNullableStringToValue(*dedicatedServer.NewNullableString(&val))

		assert.Equal(t, "value", *got)
	})

	t.Run("nil is returned if not set", func(t *testing.T) {
		got := AdaptNullableStringToValue(*dedicatedServer.NewNullableString(nil))

		assert.Nil(t, got)
	})
}
