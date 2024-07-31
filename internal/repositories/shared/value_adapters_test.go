package shared

import (
	"fmt"
	"testing"
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestAdaptNullableStringToValue(t *testing.T) {
	t.Run("value is returned if set", func(t *testing.T) {
		val := "value"

		got := AdaptNullableStringToValue(*publicCloud.NewNullableString(&val))

		assert.Equal(t, "value", *got)
	})

	t.Run("nil is returned if not set", func(t *testing.T) {
		got := AdaptNullableStringToValue(*publicCloud.NewNullableString(nil))

		assert.Nil(t, got)
	})
}

func ExampleAdaptNullableStringToValue() {
	val := "value"
	adaptedValue := AdaptNullableStringToValue(*publicCloud.NewNullableString(&val))

	fmt.Println(*adaptedValue)
	// Output: value
}

func TestAdaptNullableTimeToValue(t *testing.T) {
	t.Run("value is returned if set", func(t *testing.T) {
		val := time.Now()

		got := AdaptNullableTimeToValue(*publicCloud.NewNullableTime(&val))

		assert.Equal(t, val, *got)
	})

	t.Run("nil is returned if not set", func(t *testing.T) {
		got := AdaptNullableTimeToValue(*publicCloud.NewNullableTime(nil))

		assert.Nil(t, got)
	})
}

func ExampleAdaptNullableTimeToValue() {
	val, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")

	adaptedValue := AdaptNullableTimeToValue(*publicCloud.NewNullableTime(&val))

	fmt.Println(*adaptedValue)
	// Output: 2019-09-08 00:00:00 +0000 UTC
}

func TestAdaptNullableInt32ToValue(t *testing.T) {
	t.Run("value is returned if set", func(t *testing.T) {
		val := int32(2)

		got := AdaptNullableInt32ToValue(*publicCloud.NewNullableInt32(&val))

		assert.Equal(t, int(val), *got)
	})

	t.Run("nil is returned if not set", func(t *testing.T) {
		got := AdaptNullableInt32ToValue(*publicCloud.NewNullableInt32(nil))

		assert.Nil(t, got)
	})
}

func ExampleAdaptNullableInt32ToValue() {
	val := int32(2)
	adaptedValue := AdaptNullableInt32ToValue(*publicCloud.NewNullableInt32(&val))

	fmt.Println(*adaptedValue)
	// Output: 2
}
