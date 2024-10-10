package shared

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	dataSourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
	"github.com/stretchr/testify/assert"
)

type mockDomainEntity struct {
}

type mockModel struct {
	Value string `tfsdk:"value"`
}

func TestAdaptNullableIntToInt64Value(t *testing.T) {
	value := 1234

	type args struct {
		value *int
	}
	tests := []struct {
		name string
		args args
		want basetypes.Int64Value
	}{
		{
			name: "value has been set to nil",
			args: args{value: nil},
			want: basetypes.NewInt64Null(),
		},
		{
			name: "value has been set",
			args: args{value: &value},
			want: basetypes.NewInt64Value(1234),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				AdaptNullableIntToInt64Value(tt.args.value),
				"AdaptNullableIntToInt64Value(%v)",
				tt.args.value,
			)
		})
	}
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

func TestAdaptNullableDomainEntityToDatasourceModel(t *testing.T) {
	entity := mockDomainEntity{}
	mockGenerator := func(domainEntity mockDomainEntity) *string {
		value := "tralala"
		return &value
	}

	t.Run("value is nil", func(t *testing.T) {
		got := AdaptNullableDomainEntityToDatasourceModel(nil, mockGenerator)
		assert.Nil(t, got)
	})
	t.Run("value is set", func(t *testing.T) {
		got := AdaptNullableDomainEntityToDatasourceModel(&entity, mockGenerator)
		assert.Equal(t, "tralala", *got)
	})
}

func TestAdaptNullableDomainEntityToResourceObject(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run("value is nil", func(t *testing.T) {
		got, gotDiags := AdaptNullableDomainEntityToResourceObject(
			nil,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				entity mockDomainEntity,
			) (model *mockModel, err error) {

				return &mockModel{Value: "tralala"}, nil
			},
		)
		assert.Nil(t, gotDiags)
		assert.Equal(t, types.ObjectNull(map[string]attr.Type{}), got)
	})

	t.Run("value is set", func(t *testing.T) {
		got, gotDiags := AdaptNullableDomainEntityToResourceObject(
			&entity,
			map[string]attr.Type{"value": types.StringType},
			context.TODO(),
			func(
				ctx context.Context,
				entity mockDomainEntity,
			) (model *mockModel, err error) {

				return &mockModel{Value: "tralala"}, nil
			},
		)

		assert.Nil(t, gotDiags)
		assert.Equal(t, "\"tralala\"", got.Attributes()["value"].String())
	})

	t.Run("generateTerraformModel returns an error", func(t *testing.T) {
		got, err := AdaptNullableDomainEntityToResourceObject(
			&entity,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				entity mockDomainEntity,
			) (model *mockModel, err error) {
				return nil, errors.New("tralala")
			},
		)

		assert.Equal(t, types.ObjectUnknown(map[string]attr.Type{}), got)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}

func TestAdaptDomainEntityToResourceObject(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run("generateTerraformModel returns an error", func(t *testing.T) {
		got, err := AdaptDomainEntityToResourceObject(
			entity,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				entity mockDomainEntity,
			) (model *mockModel, err error) {
				return nil, errors.New("tralala")
			},
		)

		assert.Equal(t, types.ObjectUnknown(map[string]attr.Type{}), got)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("attributeTypes are incorrect", func(t *testing.T) {
		got, err := AdaptDomainEntityToResourceObject(
			entity,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				entity mockDomainEntity,
			) (model *mockModel, err error) {

				return &mockModel{}, nil
			},
		)

		assert.Equal(t, types.ObjectUnknown(map[string]attr.Type{}), got)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Value Conversion Error")
	})

	t.Run("sdkModel is processed properly", func(t *testing.T) {
		got, diags := AdaptDomainEntityToResourceObject(
			entity,
			map[string]attr.Type{"value": types.StringType},
			context.TODO(),
			func(
				ctx context.Context,
				entity mockDomainEntity,
			) (*mockModel, error) {

				return &mockModel{Value: "tralala"}, nil
			},
		)

		assert.Nil(t, diags)
		assert.Equal(t, "\"tralala\"", got.Attributes()["value"].String())
	})
}

func TestAdaptNullableStringToStringValue(t *testing.T) {
	value := "tralala"

	type args struct {
		value *string
	}
	tests := []struct {
		name string
		args args
		want basetypes.StringValue
	}{
		{
			name: "value has been set to nil",
			args: args{value: nil},
			want: basetypes.NewStringNull(),
		},
		{
			name: "value has been set",
			args: args{value: &value},
			want: basetypes.NewStringValue("tralala"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				AdaptNullableStringToStringValue(tt.args.value),
				"AdaptNullableStringToStringValue(%v)",
				tt.args.value,
			)
		})
	}
}

func TestAdaptDomainSliceToListValue(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run(
		"slice can successfully be converted into a ListValue",
		func(t *testing.T) {
			got, diags := AdaptEntitiesToListValue(
				[]mockDomainEntity{entity},
				map[string]attr.Type{"value": types.StringType},
				context.TODO(),
				func(
					ctx context.Context,
					entity mockDomainEntity,
				) (*mockModel, error) {

					return &mockModel{Value: "tralala"}, nil
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
		"error is returned if list element cannot be converted",
		func(t *testing.T) {
			_, err := AdaptEntitiesToListValue(
				[]mockDomainEntity{entity},
				map[string]attr.Type{"value": types.StringType},
				context.TODO(),
				func(
					ctx context.Context,
					entity mockDomainEntity,
				) (*mockModel, error) {
					return nil, errors.New("tralala")
				},
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"error is returned if passed attributeTypes are incorrect",
		func(t *testing.T) {
			_, err := AdaptEntitiesToListValue(
				[]mockDomainEntity{entity},
				map[string]attr.Type{},
				context.TODO(),
				func(
					ctx context.Context,
					entity mockDomainEntity,
				) (*mockModel, error) {

					return &mockModel{Value: "tralala"}, nil
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

func ExampleAdaptNullableIntToInt64Value() {
	nullableInt := 64
	value := AdaptNullableIntToInt64Value(&nullableInt)

	fmt.Println(value)
	// Output: 64
}

func ExampleAdaptNullableIntToInt64Value_second() {
	value := AdaptNullableIntToInt64Value(nil)

	fmt.Println(value)
	// Output: <null>
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

func ExampleAdaptNullableStringToStringValue() {
	nullableString := "tralala"
	value := AdaptNullableStringToStringValue(&nullableString)

	fmt.Println(value)
	// Output: "tralala"
}

func ExampleAdaptNullableStringToStringValue_second() {
	value := AdaptNullableStringToStringValue(nil)

	fmt.Println(value)
	// Output: <null>
}

func ExampleAdaptNullableDomainEntityToDatasourceModel() {
	image := public_cloud.NewImage(
		"imageId",
		"",
		"",
		"",
		false,
	)

	datasourceModel := AdaptNullableDomainEntityToDatasourceModel(
		&image,
		func(image public_cloud.Image) *dataSourceModel.Image {
			return &dataSourceModel.Image{
				Id: basetypes.NewStringValue(image.Id),
			}
		},
	)

	fmt.Println(datasourceModel)
	// Output: &{"imageId"}
}

func ExampleAdaptNullableDomainEntityToDatasourceModel_second() {
	datasourceModel := AdaptNullableDomainEntityToDatasourceModel(
		nil,
		func(image public_cloud.Image) *dataSourceModel.Image {
			return &dataSourceModel.Image{
				Id: basetypes.NewStringValue(image.Id),
			}
		},
	)

	fmt.Println(datasourceModel)
	// Output: <nil>
}

func ExampleAdaptNullableDomainEntityToResourceObject() {
	image := public_cloud.NewImage(
		"imageId",
		"",
		"",
		"",
		false,
	)

	datasourceModel, _ := AdaptNullableDomainEntityToResourceObject(
		&image,
		map[string]attr.Type{
			"id": types.StringType,
		},
		context.TODO(),
		func(ctx context.Context, image public_cloud.Image) (*model.Image, error) {
			return &model.Image{
				Id: basetypes.NewStringValue(image.Id),
			}, nil
		},
	)

	fmt.Println(datasourceModel)
	// Output: {"id":"imageId"}
}

func ExampleAdaptNullableDomainEntityToResourceObject_second() {
	datasourceModel, _ := AdaptNullableDomainEntityToResourceObject(
		nil,
		map[string]attr.Type{
			"id": types.StringType,
		},
		context.TODO(),
		func(ctx context.Context, image public_cloud.Image) (*model.Image, error) {
			return &model.Image{
				Id: basetypes.NewStringValue(image.Id),
			}, nil
		},
	)

	fmt.Println(datasourceModel)
	// Output: <null>
}

func ExampleAdaptDomainEntityToResourceObject() {
	datasourceModel, _ := AdaptDomainEntityToResourceObject(
		public_cloud.NewImage(
			"imageId",
			"",
			"",
			"",
			false,
		),
		map[string]attr.Type{
			"id": types.StringType,
		},
		context.TODO(),
		func(ctx context.Context, image public_cloud.Image) (*model.Image, error) {
			return &model.Image{
				Id: basetypes.NewStringValue(image.Id),
			}, nil
		},
	)

	fmt.Println(datasourceModel)
	// Output: {"id":"imageId"}
}

func ExampleAdaptEntitiesToListValue() {
	listValue, _ := AdaptEntitiesToListValue(
		public_cloud.Ips{public_cloud.NewIp(
			"1.2.3.4",
			"",
			2,
			false,
			true,
			enum.NetworkTypeInternal,
			public_cloud.OptionalIpValues{},
		)},
		map[string]attr.Type{
			"ip": types.StringType,
		},
		context.TODO(),
		func(ctx context.Context, entity public_cloud.Ip) (*model.Ip, error) {
			return &model.Ip{
				Ip: basetypes.NewStringValue(entity.Ip),
			}, nil
		},
	)

	fmt.Println(listValue)
	// Output: [{"ip":"1.2.3.4"}]
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

func TestAdaptIntArrayToInt64(t *testing.T) {
	want := []int64{5}
	got := AdaptIntArrayToInt64Array([]int{5})

	assert.Equal(t, want, got)
}

func ExampleAdaptIntArrayToInt64Array() {
	convertedValue := AdaptIntArrayToInt64Array([]int{5})

	fmt.Println(convertedValue)
	// Output: [5]
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

func ExampleReturnError() {
	diags := diag.Diagnostics{}
	diags.AddError("summary", "detail")

	returnedErrors := ReturnError("functionName", diags)

	fmt.Println(returnedErrors)
	// Output:  functionName: "summary" "detail"
}

func TestAdaptNullableBoolToBoolValue(t *testing.T) {
	t.Run(
		"returns a boolean when the value is not nil",
		func(t *testing.T) {
			value := true
			got := AdaptNullableBoolToBoolValue(&value)

			assert.True(t, got.ValueBool())
			assert.False(t, got.IsNull())
		},
	)

	t.Run(
		"returns null when the value is nil",
		func(t *testing.T) {
			got := AdaptNullableBoolToBoolValue(nil)

			assert.True(t, got.IsNull())
		},
	)
}

func ExampleAdaptNullableBoolToBoolValue() {
	value := true
	adaptedValue := AdaptNullableBoolToBoolValue(&value)

	fmt.Println(adaptedValue)
	// Output: true
}

func ExampleAdaptNullableBoolToBoolValue_second() {
	adaptedValue := AdaptNullableBoolToBoolValue(nil)

	fmt.Println(adaptedValue)
	// Output: <null>
}

func TestAdaptBoolToBoolValue(t *testing.T) {
	t.Run("works with false", func(t *testing.T) {
		got := AdaptBoolToBoolValue(false)

		assert.False(t, got.ValueBool())
	})

	t.Run("works with true", func(t *testing.T) {
		got := AdaptBoolToBoolValue(true)

		assert.True(t, got.ValueBool())
	})
}

func ExampleAdaptBoolToBoolValue() {
	adaptedValue := AdaptBoolToBoolValue(true)

	fmt.Println(adaptedValue)
	// Output: true
}
