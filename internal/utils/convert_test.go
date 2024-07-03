package utils

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func mockGenerator(sdkEntity mockSdkEntity) *string {
	value := "tralala"
	return &value
}

type mockSdkEntity struct {
}

type mockModel struct {
	Value string `tfsdk:"value"`
}

func TestConvertNullableSdkIntToInt64Value(t *testing.T) {
	value := int32(1234)

	type args struct {
		value *int32
		ok    bool
	}
	tests := []struct {
		name string
		args args
		want basetypes.Int64Value
	}{
		{
			name: "value is not set",
			args: args{ok: false, value: &value},
			want: basetypes.NewInt64Null(),
		},
		{
			name: "value has been set to nil",
			args: args{ok: true, value: nil},
			want: basetypes.NewInt64Null(),
		},
		{
			name: "value has been set",
			args: args{ok: true, value: &value},
			want: basetypes.NewInt64Value(1234),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				ConvertNullableSdkIntToInt64Value(tt.args.value, tt.args.ok),
				"ConvertNullableSdkIntToInt64Value(%v, %v)",
				tt.args.value,
				tt.args.ok,
			)
		})
	}
}

func TestConvertNullableSdkTimeToStringValue(t *testing.T) {
	value, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")

	type args struct {
		value *time.Time
		ok    bool
	}
	tests := []struct {
		name string
		args args
		want basetypes.StringValue
	}{
		{
			name: "value is not set",
			args: args{ok: false, value: &value},
			want: basetypes.NewStringNull(),
		},
		{
			name: "time is not set",
			args: args{nil, true},
			want: basetypes.NewStringNull(),
		},
		{
			name: "time is set",
			args: args{&value, true},
			want: basetypes.NewStringValue("2019-09-08 00:00:00 +0000 UTC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ConvertNullableSdkTimeToStringValue(
				tt.args.value,
				tt.args.ok,
			), "ConvertNullableSdkTimeToStringValue(%v)", tt.args.value)
		})
	}
}

func TestConvertNullableSdkModelToDatasourceModel(t *testing.T) {
	sdkEntity := mockSdkEntity{}

	t.Run("value is not set", func(t *testing.T) {
		got := ConvertNullableSdkModelToDatasourceModel(&sdkEntity, false, mockGenerator)
		assert.Nil(t, got)
	})
	t.Run("value is nil", func(t *testing.T) {
		got := ConvertNullableSdkModelToDatasourceModel(nil, true, mockGenerator)
		assert.Nil(t, got)
	})
	t.Run("value is set", func(t *testing.T) {
		got := ConvertNullableSdkModelToDatasourceModel(&sdkEntity, true, mockGenerator)
		assert.Equal(t, "tralala", *got)
	})
}

func TestConvertNullableSdkModelToResourceObject(t *testing.T) {
	sdkEntity := mockSdkEntity{}

	t.Run("value is not set", func(t *testing.T) {
		got, gotDiags := ConvertNullableSdkModelToResourceObject(
			&sdkEntity,
			false,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				sdkEntity mockSdkEntity,
			) (model *mockModel, diagnostics diag.Diagnostics) {

				return &mockModel{Value: "tralala"}, nil
			},
		)
		assert.Equal(t, types.ObjectNull(map[string]attr.Type{}), got)
		assert.Nil(t, gotDiags)
	})

	t.Run("value is nil", func(t *testing.T) {
		got, gotDiags := ConvertNullableSdkModelToResourceObject(
			nil,
			true,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				sdkEntity mockSdkEntity,
			) (model *mockModel, diagnostics diag.Diagnostics) {

				return &mockModel{Value: "tralala"}, nil
			},
		)
		assert.Equal(t, types.ObjectNull(map[string]attr.Type{}), got)
		assert.Nil(t, gotDiags)
	})

	t.Run("value is set", func(t *testing.T) {
		got, gotDiags := ConvertNullableSdkModelToResourceObject(
			&sdkEntity,
			true,
			map[string]attr.Type{"value": types.StringType},
			context.TODO(),
			func(
				ctx context.Context,
				sdkEntity mockSdkEntity,
			) (model *mockModel, diagnostics diag.Diagnostics) {

				return &mockModel{Value: "tralala"}, nil
			},
		)

		assert.Nil(t, gotDiags)
		assert.Equal(t, "\"tralala\"", got.Attributes()["value"].String())
	})

	t.Run("generateTerraformModel returns an error", func(t *testing.T) {
		got, diags := ConvertNullableSdkModelToResourceObject(
			&sdkEntity,
			true,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				sdkEntity mockSdkEntity,
			) (model *mockModel, diagnostics diag.Diagnostics) {
				diagnostics.AddError("tralala", "")

				return nil, diagnostics
			},
		)

		assert.Equal(t, types.ObjectNull(map[string]attr.Type{}), got)
		assert.Equal(t, 1, diags.ErrorsCount())
		assert.Equal(t, "tralala", diags[0].Summary())
	})
}

func TestConvertSdkModelToResourceObject(t *testing.T) {
	sdkEntity := mockSdkEntity{}

	t.Run("generateTerraformModel returns an error", func(t *testing.T) {
		got, diags := ConvertSdkModelToResourceObject(
			sdkEntity,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				sdkEntity mockSdkEntity,
			) (model *mockModel, diagnostics diag.Diagnostics) {
				diagnostics.AddError("tralala", "")

				return nil, diagnostics
			},
		)

		assert.Equal(t, types.ObjectNull(map[string]attr.Type{}), got)
		assert.Equal(t, 1, diags.ErrorsCount())
		assert.Equal(t, "tralala", diags[0].Summary())
	})

	t.Run("attributeTypes are incorrect", func(t *testing.T) {
		got, diags := ConvertSdkModelToResourceObject(
			sdkEntity,
			map[string]attr.Type{},
			context.TODO(),
			func(
				ctx context.Context,
				sdkEntity mockSdkEntity,
			) (model *mockModel, diagnostics diag.Diagnostics) {

				return &mockModel{}, nil
			},
		)

		assert.Equal(t, types.ObjectNull(map[string]attr.Type{}), got)
		assert.Equal(t, 1, diags.ErrorsCount())
		assert.Equal(t, "Value Conversion Error", diags[0].Summary())
	})

	t.Run("sdkModel is processed properly", func(t *testing.T) {
		got, diags := ConvertSdkModelToResourceObject(
			sdkEntity,
			map[string]attr.Type{"value": types.StringType},
			context.TODO(),
			func(
				ctx context.Context,
				sdkEntity mockSdkEntity,
			) (model *mockModel, diagnostics diag.Diagnostics) {

				return &mockModel{Value: "tralala"}, nil
			},
		)

		assert.Nil(t, diags)
		assert.Equal(t, "\"tralala\"", got.Attributes()["value"].String())
	})
}

func TestConvertNullableSdkStringToInt64Value(t *testing.T) {
	value := "tralala"

	type args struct {
		value *string
		ok    bool
	}
	tests := []struct {
		name string
		args args
		want basetypes.StringValue
	}{
		{
			name: "value is not set",
			args: args{ok: false, value: &value},
			want: basetypes.NewStringNull(),
		},
		{
			name: "value has been set to nil",
			args: args{ok: true, value: nil},
			want: basetypes.NewStringNull(),
		},
		{
			name: "value has been set",
			args: args{ok: true, value: &value},
			want: basetypes.NewStringValue("tralala"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				ConvertNullableSdkStringToStringValue(tt.args.value, tt.args.ok),
				"ConvertNullableSdkIntToInt64Value(%v, %v)",
				tt.args.value,
				tt.args.ok,
			)
		})
	}
}
