package shared

import (
	"time"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

func AdaptNullableStringToValue(nullableString publicCloud.NullableString) *string {
	return nullableString.Get()
}

func AdaptNullableTimeToValue(nullableTime publicCloud.NullableTime) *time.Time {
	return nullableTime.Get()
}

func AdaptNullableInt32ToValue(nullableInt publicCloud.NullableInt32) *int {
	if nullableInt.Get() == nil {
		return nil
	}

	value := int(*nullableInt.Get())
	return &value
}
