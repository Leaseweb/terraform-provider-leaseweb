package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAutoScalingGroupReference(t *testing.T) {
	t.Run("valid reference is set", func(t *testing.T) {
		got, err := NewAutoScalingGroupReference("value")

		assert.Nil(t, err)
		assert.Equal(t, "value", got.value)
	})

	t.Run("error is returned if reference is too long", func(t *testing.T) {
		_, err := NewAutoScalingGroupReference("zgbnjmbdyquzkyzhkuvdgkkaxfxnwbrmhrdiiqutrgqpmymwykettmnhbnnfvnpxziqebmnjybhqfnnvraqqwyueqmddxfnxgprxjyqmdrmcgpmptuvprhhezwextdnrcudpfkmqfrxjcyjvxgbamxgjvyhtzgdvfdvrzzabmviehkyzcdikyumcwyqgkuegvwdfmjehnnwunuhztcfyarpzbmifkfdpwuyprdgegazdbruyftvpejgdrwcaaaaa")

		assert.NotNil(t, err)
	})

	t.Run("string value is correct", func(t *testing.T) {
		got, err := NewAutoScalingGroupReference("value")

		assert.Nil(t, err)
		assert.Equal(t, "value", got.String())
	})

}
