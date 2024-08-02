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
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	dataSourceModel "terraform-provider-leaseweb/internal/provider/data_sources/public_cloud/model"
	"terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
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
	iso := domain.NewIso("id", "name")

	datasourceModel := AdaptNullableDomainEntityToDatasourceModel(
		&iso,
		func(iso domain.Iso) *dataSourceModel.Iso {
			return &dataSourceModel.Iso{
				Id:   basetypes.NewStringValue(iso.Id),
				Name: basetypes.NewStringValue(iso.Name),
			}
		},
	)

	fmt.Println(datasourceModel)
	// Output: &{"id" "name"}
}

func ExampleAdaptNullableDomainEntityToDatasourceModel_second() {
	datasourceModel := AdaptNullableDomainEntityToDatasourceModel(
		nil,
		func(iso domain.Iso) *dataSourceModel.Iso {
			return &dataSourceModel.Iso{
				Id:   basetypes.NewStringValue(iso.Id),
				Name: basetypes.NewStringValue(iso.Name),
			}
		},
	)

	fmt.Println(datasourceModel)
	// Output: <nil>
}

func ExampleAdaptNullableDomainEntityToResourceObject() {
	iso := domain.NewIso("id", "name")

	datasourceModel, _ := AdaptNullableDomainEntityToResourceObject(
		&iso,
		map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
		context.TODO(),
		func(ctx context.Context, iso domain.Iso) (*model.Iso, error) {
			return &model.Iso{
				Id:   basetypes.NewStringValue(iso.Id),
				Name: basetypes.NewStringValue(iso.Name),
			}, nil
		},
	)

	fmt.Println(datasourceModel)
	// Output: {"id":"id","name":"name"}
}

func ExampleAdaptNullableDomainEntityToResourceObject_second() {
	datasourceModel, _ := AdaptNullableDomainEntityToResourceObject(
		nil,
		map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
		context.TODO(),
		func(ctx context.Context, iso domain.Iso) (*model.Iso, error) {
			return &model.Iso{
				Id:   basetypes.NewStringValue(iso.Id),
				Name: basetypes.NewStringValue(iso.Name),
			}, nil
		},
	)

	fmt.Println(datasourceModel)
	// Output: <null>
}

func ExampleAdaptDomainEntityToResourceObject() {

	datasourceModel, _ := AdaptDomainEntityToResourceObject(
		domain.NewIso("id", "name"),
		map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
		context.TODO(),
		func(ctx context.Context, iso domain.Iso) (*model.Iso, error) {
			return &model.Iso{
				Id:   basetypes.NewStringValue(iso.Id),
				Name: basetypes.NewStringValue(iso.Name),
			}, nil
		},
	)

	fmt.Println(datasourceModel)
	// Output: {"id":"id","name":"name"}
}

func ExampleAdaptEntitiesToListValue() {
	listValue, _ := AdaptEntitiesToListValue(
		domain.Ips{domain.NewIp(
			"1.2.3.4",
			"prefixLength",
			2,
			false,
			true,
			enum.NetworkTypeInternal,
			domain.OptionalIpValues{},
		)},
		map[string]attr.Type{
			"ip":             types.StringType,
			"prefix_length":  types.StringType,
			"version":        types.Int64Type,
			"null_routed":    types.BoolType,
			"main_ip":        types.BoolType,
			"network_type":   types.StringType,
			"reverse_lookup": types.StringType,
			"ddos":           types.ObjectType{AttrTypes: model.Ddos{}.AttributeTypes()},
		},
		context.TODO(),
		func(ctx context.Context, entity domain.Ip) (*model.Ip, error) {
			ddos, _ := AdaptNullableDomainEntityToResourceObject(
				entity.Ddos,
				model.Ddos{}.AttributeTypes(),
				ctx,
				func(ctx context.Context, ddos domain.Ddos) (*model.Ddos, error) {
					return &model.Ddos{
						DetectionProfile: basetypes.NewStringValue(ddos.DetectionProfile),
						ProtectionType:   basetypes.NewStringValue(ddos.ProtectionType),
					}, nil
				},
			)

			return &model.Ip{
				Ip:            basetypes.NewStringValue(entity.Ip),
				PrefixLength:  basetypes.NewStringValue(entity.PrefixLength),
				Version:       basetypes.NewInt64Value(int64(entity.Version)),
				NullRouted:    basetypes.NewBoolValue(entity.NullRouted),
				MainIp:        basetypes.NewBoolValue(entity.MainIp),
				NetworkType:   basetypes.NewStringValue(string(entity.NetworkType)),
				ReverseLookup: basetypes.NewStringPointerValue(entity.ReverseLookup),
				Ddos:          ddos,
			}, nil
		},
	)

	fmt.Println(listValue)
	// Output: [{"ddos":<null>,"ip":"1.2.3.4","main_ip":true,"network_type":"INTERNAL","null_routed":false,"prefix_length":"prefixLength","reverse_lookup":<null>,"version":2}]
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
