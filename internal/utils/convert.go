package utils

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ConvertNullableSdkIntToInt64Value Convert SDK NullableInt to terraform Int64.
func ConvertNullableSdkIntToInt64Value(value *int32, ok bool) basetypes.Int64Value {
	if value == nil || !ok {
		return basetypes.NewInt64Null()
	}

	return basetypes.NewInt64Value(int64(*value))
}

// ConvertNullableSdkTimeToStringValue Convert SDK NullableTime to terraform String.
func ConvertNullableSdkTimeToStringValue(value *time.Time, ok bool) basetypes.StringValue {
	if value == nil || !ok {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(value.String())
}

// ConvertNullableSdkStringToStringValue Convert SDK NullableString to terraform String.
func ConvertNullableSdkStringToStringValue(value *string, ok bool) basetypes.StringValue {
	if value == nil || !ok {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(*value)
}

// ConvertNullableSdkModelToDatasourceModel Convert nullable SDK model to datasource model.
func ConvertNullableSdkModelToDatasourceModel[T interface{}, U interface{}](
	value *T,
	ok bool,
	generateModel func(sdkEntity T) *U,
) *U {
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}

	return generateModel(*value)
}

// ConvertNullableSdkModelToResourceObject Convert nullable SDK model to resource object.
func ConvertNullableSdkModelToResourceObject[T any, U any](
	sdkEntity *T,
	ok bool,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateTerraformModel func(
		ctx context.Context,
		sdkEntity T,
	) (
		model *U,
		diagnostics diag.Diagnostics,
	)) (basetypes.ObjectValue, diag.Diagnostics) {
	if !ok {
		return types.ObjectNull(attributeTypes), nil
	}
	if sdkEntity == nil {
		return types.ObjectNull(attributeTypes), nil
	}

	object, diags := ConvertSdkModelToResourceObject(
		*sdkEntity,
		attributeTypes,
		ctx,
		generateTerraformModel,
	)

	if diags.HasError() {
		return types.ObjectNull(attributeTypes), diags
	}

	return object, nil
}

// ConvertSdkModelToResourceObject Convert SDK model to resource object.
func ConvertSdkModelToResourceObject[T any, U any](
	sdkEntity T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateTerraformModel func(
		ctx context.Context,
		sdkEntity T,
	) (model *U, diagnostics diag.Diagnostics),
) (basetypes.ObjectValue, diag.Diagnostics) {
	terraformModel, diags := generateTerraformModel(ctx, sdkEntity)
	if diags.HasError() {
		return types.ObjectNull(attributeTypes), diags
	}

	objectValue, diags := types.ObjectValueFrom(
		ctx,
		attributeTypes,
		terraformModel,
	)
	if diags != nil {
		return types.ObjectNull(attributeTypes), diags
	}

	return objectValue, nil
}
