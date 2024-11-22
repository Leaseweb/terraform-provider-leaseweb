package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestAdaptInt32PointerValueToNullableInt32(t *testing.T) {
	t.Run("passing nil returns nil", func(t *testing.T) {
		got := AdaptInt32PointerValueToNullableInt32(
			basetypes.NewInt32PointerValue(nil),
		)
		assert.Nil(t, got)
	})

	t.Run("passing value returns int32 variant", func(t *testing.T) {
		value := int32(1)
		int32Value := basetypes.NewInt32PointerValue(&value)
		got := AdaptInt32PointerValueToNullableInt32(int32Value)

		assert.Equal(t, int32(1), *got)
	})

	t.Run("passing unknown returns nil", func(t *testing.T) {
		got := AdaptInt32PointerValueToNullableInt32(basetypes.NewInt32Unknown())
		assert.Nil(t, got)
	})
}

func ExampleAdaptInt32PointerValueToNullableInt32() {
	value := int32(3)
	int32Value := basetypes.NewInt32PointerValue(&value)
	adaptedValue := AdaptInt32PointerValueToNullableInt32(int32Value)

	fmt.Println(*adaptedValue)
	// Output: 3
}

func ExampleAdaptInt32PointerValueToNullableInt32_second() {
	int32Value := basetypes.NewInt32PointerValue(nil)
	adaptedValue := AdaptInt32PointerValueToNullableInt32(int32Value)

	fmt.Println(adaptedValue)
	// Output: <nil>
}

type mockDomainEntity struct {
}

type mockModel struct {
	Value string `tfsdk:"value"`
}

func TestAdaptNullableTimeToStringValue(t *testing.T) {
	value, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")

	type args struct {
		value *time.Time
	}
	tests := []struct {
		name string
		args args
		want basetypes.StringValue
	}{
		{
			name: "time is not set",
			args: args{value: nil},
			want: basetypes.NewStringNull(),
		},
		{
			name: "time is set",
			args: args{value: &value},
			want: basetypes.NewStringValue("2019-09-08 00:00:00 +0000 UTC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, AdaptNullableTimeToStringValue(
				tt.args.value,
			), "AdaptNullableTimeToStringValue(%v)", tt.args.value)
		})
	}
}

func TestAdaptDomainEntityToResourceObject(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run("attributeTypes are incorrect", func(t *testing.T) {
		got, err := AdaptSdkModelToResourceObject(
			entity,
			map[string]attr.Type{},
			context.TODO(),
			func(entity mockDomainEntity) (model mockModel) {
				return mockModel{}
			},
		)

		assert.Equal(t, types.ObjectUnknown(map[string]attr.Type{}), got)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Value Conversion Error")
	})

	t.Run("sdkModel is processed properly", func(t *testing.T) {
		got, diags := AdaptSdkModelToResourceObject(
			entity,
			map[string]attr.Type{"value": types.StringType},
			context.TODO(),
			func(entity mockDomainEntity) mockModel {
				return mockModel{Value: "tralala"}
			},
		)

		assert.Nil(t, diags)
		assert.Equal(t, "\"tralala\"", got.Attributes()["value"].String())
	})
}

func TestAdaptDomainSliceToListValue(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run(
		"slice can successfully be converted into a ListValue",
		func(t *testing.T) {
			got, diags := AdaptSdkModelsToListValue(
				[]mockDomainEntity{entity},
				map[string]attr.Type{"value": types.StringType},
				context.TODO(),
				func(entity mockDomainEntity) mockModel {
					return mockModel{Value: "tralala"}
				},
			)

			assert.Nil(t, diags)
			assert.Len(t, got.Elements(), 1)
			assert.Equal(
				t,
				"{\"value\":\"tralala\"}",
				got.Elements()[0].String(),
			)
		},
	)

	t.Run(
		"error is returned if passed attributeTypes are incorrect",
		func(t *testing.T) {
			_, err := AdaptSdkModelsToListValue(
				[]mockDomainEntity{entity},
				map[string]attr.Type{},
				context.TODO(),
				func(entity mockDomainEntity) mockModel {
					return mockModel{Value: "tralala"}
				},
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "Value Conversion Error")
		},
	)
}

func TestAdaptStringPointerValueToNullableString(t *testing.T) {
	t.Run("returns nil when value is unknown", func(t *testing.T) {
		value := basetypes.NewStringUnknown()
		assert.Nil(t, AdaptStringPointerValueToNullableString(value))
	})

	t.Run("returns pointer when value is set", func(t *testing.T) {
		target := "tralala"
		value := basetypes.NewStringPointerValue(&target)

		assert.Equal(t, target, *AdaptStringPointerValueToNullableString(value))
	})

	t.Run("returns nil when value is not set", func(t *testing.T) {
		value := basetypes.NewStringPointerValue(nil)

		assert.Nil(t, AdaptStringPointerValueToNullableString(value))
	})
}

func ExampleAdaptNullableTimeToStringValue() {
	nullableTime, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	value := AdaptNullableTimeToStringValue(&nullableTime)

	fmt.Println(value)
	// Output: "2019-09-08 00:00:00 +0000 UTC"
}

func ExampleAdaptNullableTimeToStringValue_second() {
	value := AdaptNullableTimeToStringValue(nil)

	fmt.Println(value)
	// Output: <null>
}

func ExampleAdaptSdkModelToResourceObject() {
	type Image struct {
		Id types.String `tfsdk:"id"`
	}

	resourceModel, _ := AdaptSdkModelToResourceObject(
		publicCloud.Image{Id: "imageId"},
		map[string]attr.Type{
			"id": types.StringType,
		},
		context.TODO(),
		func(image publicCloud.Image) Image {
			return Image{
				Id: basetypes.NewStringValue(image.Id),
			}
		},
	)

	fmt.Println(resourceModel)
	// Output: {"id":"imageId"}
}

func ExampleAdaptSdkModelsToListValue() {
	type Ip struct {
		Ip types.String `tfsdk:"ip"`
	}

	listValue, _ := AdaptSdkModelsToListValue(
		[]publicCloud.Ip{{Ip: "1.2.3.4"}},
		map[string]attr.Type{
			"ip": types.StringType,
		},
		context.TODO(),
		func(ip publicCloud.Ip) Ip {
			return Ip{
				Ip: basetypes.NewStringValue(ip.Ip),
			}
		},
	)

	fmt.Println(listValue)
	// Output: [{"ip":"1.2.3.4"}]
}

func TestReturnError(t *testing.T) {
	t.Run("diagnostics contain errors", func(t *testing.T) {
		diags := diag.Diagnostics{}
		diags.AddError("summary", "detail")

		got := ReturnError("functionName", diags)
		want := `functionName: "summary" "detail"`

		assert.Error(t, got)
		assert.Equal(t, want, got.Error())
	})

	t.Run("diagnostics do not contain errors", func(t *testing.T) {
		diags := diag.Diagnostics{}

		got := ReturnError("functionName", diags)

		assert.NoError(t, got)
	})
}

func ExampleAdaptStringPointerValueToNullableString() {
	value := "tralala"
	terraformStringPointerValue := basetypes.NewStringPointerValue(&value)

	convertedValue := AdaptStringPointerValueToNullableString(terraformStringPointerValue)

	fmt.Println(*convertedValue)
	// Output: tralala
}

func ExampleAdaptStringPointerValueToNullableString_second() {
	terraformStringPointerValue := basetypes.NewStringPointerValue(nil)

	convertedValue := AdaptStringPointerValueToNullableString(terraformStringPointerValue)

	fmt.Println(convertedValue)
	// Output: <nil>
}

func ExampleReturnError() {
	diags := diag.Diagnostics{}
	diags.AddError("summary", "detail")

	returnedErrors := ReturnError("functionName", diags)

	fmt.Println(returnedErrors)
	// Output:  functionName: "summary" "detail"
}

func TestAdaptStringTypeArrayToStringArray(t *testing.T) {
	want := []string{"HOURLY", "MONTHLY"}
	got := AdaptStringTypeArrayToStringArray(publicCloud.AllowedContractTypeEnumValues)

	assert.Equal(t, want, got)
}

func ExampleAdaptStringTypeArrayToStringArray() {
	type customType string
	customTypes := []customType{customType("value")}

	stringTypes := AdaptStringTypeArrayToStringArray(customTypes)

	fmt.Println(stringTypes)
	// Output: [value]
}

func TestAdaptBoolPointerValueToNullableBool(t *testing.T) {
	t.Run("returns nil when value is unknown", func(t *testing.T) {
		value := basetypes.NewBoolUnknown()
		assert.Nil(t, AdaptBoolPointerValueToNullableBool(value))
	})

	t.Run("returns pointer when value is set", func(t *testing.T) {
		target := true
		value := basetypes.NewBoolPointerValue(&target)

		assert.Equal(t, target, *AdaptBoolPointerValueToNullableBool(value))
	})

	t.Run("returns nil when value is not set", func(t *testing.T) {
		value := basetypes.NewBoolPointerValue(nil)

		assert.Nil(t, AdaptBoolPointerValueToNullableBool(value))
	})
}

func ExampleAdaptBoolPointerValueToNullableBool() {
	value := true
	terraformBoolPointerValue := basetypes.NewBoolPointerValue(&value)

	convertedValue := AdaptBoolPointerValueToNullableBool(terraformBoolPointerValue)

	fmt.Println(*convertedValue)
	// Output: true
}

func TestAdaptNullableSdkModelToResourceObject(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run("attributeTypes are incorrect", func(t *testing.T) {
		got, err := AdaptNullableSdkModelToResourceObject(
			&entity,
			map[string]attr.Type{},
			context.TODO(),
			func(entity mockDomainEntity) (model mockModel) {
				return mockModel{}
			},
		)

		assert.Equal(t, types.ObjectUnknown(map[string]attr.Type{}), got)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Value Conversion Error")
	})

	t.Run("sdkModel is processed properly", func(t *testing.T) {
		got, diags := AdaptNullableSdkModelToResourceObject(
			&entity,
			map[string]attr.Type{"value": types.StringType},
			context.TODO(),
			func(entity mockDomainEntity) mockModel {
				return mockModel{Value: "tralala"}
			},
		)

		assert.Nil(t, diags)
		assert.Equal(t, "\"tralala\"", got.Attributes()["value"].String())
	})

	t.Run("passing nil returns a Null object", func(t *testing.T) {
		got, diags := AdaptNullableSdkModelToResourceObject(
			nil,
			map[string]attr.Type{"value": types.StringType},
			context.TODO(),
			func(entity mockDomainEntity) mockModel {
				return mockModel{Value: "tralala"}
			},
		)

		assert.Nil(t, diags)
		assert.True(t, got.IsNull())
	})
}

func ExampleAdaptNullableSdkModelToResourceObject() {
	type Image struct {
		Id types.String `tfsdk:"id"`
	}

	resourceModel, _ := AdaptNullableSdkModelToResourceObject(
		&publicCloud.Image{Id: "imageId"},
		map[string]attr.Type{
			"id": types.StringType,
		},
		context.TODO(),
		func(image publicCloud.Image) Image {
			return Image{
				Id: basetypes.NewStringValue(image.Id),
			}
		},
	)

	fmt.Println(resourceModel)
	// Output: {"id":"imageId"}
}

func ExampleAdaptNullableSdkModelToResourceObject_second() {
	type Image struct {
		Id types.String `tfsdk:"id"`
	}

	resourceModel, _ := AdaptNullableSdkModelToResourceObject(
		nil,
		map[string]attr.Type{
			"id": types.StringType,
		},
		context.TODO(),
		func(image publicCloud.Image) Image {
			return Image{
				Id: basetypes.NewStringValue(image.Id),
			}
		},
	)

	fmt.Println(resourceModel)
	// Output: <null>
}
