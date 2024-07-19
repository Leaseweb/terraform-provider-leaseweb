package shared

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

type mockDomainEntity struct {
}

type mockModel struct {
	Value string `tfsdk:"value"`
}

func TestConvertNullableIntToInt64Value(t *testing.T) {
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
				ConvertNullableIntToInt64Value(tt.args.value),
				"ConvertNullableIntToInt64Value(%v)",
				tt.args.value,
			)
		})
	}
}

func TestConvertNullableTimeToStringValue(t *testing.T) {
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
			assert.Equalf(t, tt.want, ConvertNullableTimeToStringValue(
				tt.args.value,
			), "ConvertNullableTimeToStringValue(%v)", tt.args.value)
		})
	}
}

func TestConvertNullableDomainEntityToDatasourceModel(t *testing.T) {
	entity := mockDomainEntity{}
	mockGenerator := func(domainEntity mockDomainEntity) *string {
		value := "tralala"
		return &value
	}

	t.Run("value is nil", func(t *testing.T) {
		got := ConvertNullableDomainEntityToDatasourceModel(nil, mockGenerator)
		assert.Nil(t, got)
	})
	t.Run("value is set", func(t *testing.T) {
		got := ConvertNullableDomainEntityToDatasourceModel(&entity, mockGenerator)
		assert.Equal(t, "tralala", *got)
	})
}

func TestConvertNullableDomainEntityToResourceObject(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run("value is nil", func(t *testing.T) {
		got, gotDiags := ConvertNullableDomainEntityToResourceObject(
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
		got, gotDiags := ConvertNullableDomainEntityToResourceObject(
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
		got, err := ConvertNullableDomainEntityToResourceObject(
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

func TestConvertDomainEntityToResourceObject(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run("generateTerraformModel returns an error", func(t *testing.T) {
		got, err := ConvertDomainEntityToResourceObject(
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
		got, err := ConvertDomainEntityToResourceObject(
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
		got, diags := ConvertDomainEntityToResourceObject(
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

func TestConvertNullableStringToStringValue(t *testing.T) {
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
				ConvertNullableStringToStringValue(tt.args.value),
				"ConvertNullableStringToStringValue(%v)",
				tt.args.value,
			)
		})
	}
}

func TestConvertDomainSliceToListValue(t *testing.T) {
	entity := mockDomainEntity{}

	t.Run(
		"slice can successfully be converted into a ListValue",
		func(t *testing.T) {
			got, diags := ConvertEntitiesToListValue(
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
			_, err := ConvertEntitiesToListValue(
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
			_, err := ConvertEntitiesToListValue(
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

func TestConvertValueStringPointerToString(t *testing.T) {
	t.Run("returns nil when value is unknown", func(t *testing.T) {
		value := basetypes.NewStringUnknown()
		assert.Nil(t, ConvertValueStringPointerToString(value))
	})

	t.Run("returns pointer when value is set", func(t *testing.T) {
		target := "tralala"
		value := basetypes.NewStringPointerValue(&target)

		assert.Equal(t, target, *ConvertValueStringPointerToString(value))
	})

	t.Run("returns nil when value is not set", func(t *testing.T) {
		value := basetypes.NewStringPointerValue(nil)

		assert.Nil(t, ConvertValueStringPointerToString(value))
	})
}
