package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"time"
)

func GetStringValue(hasValue bool, value string) basetypes.StringValue {
	if hasValue {
		return basetypes.NewStringValue(value)
	}

	return basetypes.NewStringNull()
}

func GetBoolValue(hasValue bool, value bool) basetypes.BoolValue {
	if hasValue {
		return basetypes.NewBoolValue(value)
	}

	return basetypes.NewBoolNull()
}

func GetIntValue(hasValue bool, value int32) basetypes.Int64Value {
	if hasValue {
		return basetypes.NewInt64Value(int64(value))
	}

	return basetypes.NewInt64Null()
}

func GetFloatValue(hasValue bool, value float32) basetypes.Float64Value {
	if hasValue {
		return basetypes.NewFloat64Value(float64(value))
	}

	return basetypes.NewFloat64Null()
}

func GetDateTime(value time.Time) basetypes.StringValue {
	if value.IsZero() {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(value.String())
}
