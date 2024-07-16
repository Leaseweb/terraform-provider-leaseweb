package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func TestNewHealthCheck(t *testing.T) {
	t.Run("required values are set properly", func(t *testing.T) {
		got := NewHealthCheck(
			enum.MethodPost,
			"uri",
			22,
			OptionalHealthCheckValues{},
		)

		assert.Equal(t, enum.MethodPost, got.Method)
		assert.Equal(t, "uri", got.Uri)
		assert.Equal(t, 22, got.Port)

		assert.Nil(t, got.Host)
	})

	t.Run("optional values are set properly", func(t *testing.T) {
		host := "host"

		got := NewHealthCheck(
			enum.MethodPost,
			"",
			22,
			OptionalHealthCheckValues{Host: &host},
		)

		assert.Equal(t, "host", *got.Host)
	})
}
