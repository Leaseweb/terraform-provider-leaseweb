package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAutoScalingGroupReference(t *testing.T) {
	t.Run("valid reference is set", func(t *testing.T) {
		got, err := NewAutoScalingGroupReference("Value")

		assert.Nil(t, err)
		assert.Equal(t, "Value", got.value)
	})

	t.Run("error is returned if reference is too long", func(t *testing.T) {
		_, err := NewAutoScalingGroupReference("zgbnjmbdyquzkyzhkuvdgkkaxfxnwbrmhrdiiqutrgqpmymwykettmnhbnnfvnpxziqebmnjybhqfnnvraqqwyueqmddxfnxgprxjyqmdrmcgpmptuvprhhezwextdnrcudpfkmqfrxjcyjvxgbamxgjvyhtzgdvfdvrzzabmviehkyzcdikyumcwyqgkuegvwdfmjehnnwunuhztcfyarpzbmifkfdpwuyprdgegazdbruyftvpejgdrwcaaaaa")

		assert.NotNil(t, err)
	})

	t.Run("string Value is correct", func(t *testing.T) {
		got, err := NewAutoScalingGroupReference("Value")

		assert.Nil(t, err)
		assert.Equal(t, "Value", got.String())
	})

}
