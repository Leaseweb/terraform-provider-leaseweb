package utils

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AdaptInt32PointerValueToNullableInt32 converts a Terraform
// Int32PointerValue to a nullable int32.
func AdaptInt32PointerValueToNullableInt32(int32Type types.Int32) *int32 {
	if int32Type.IsUnknown() {
		return nil
	}

	value := int32Type.ValueInt32Pointer()
	if value == nil {
		return nil
	}

	return value
}

// AdaptNullableTimeToStringValue converts a nullable Time to a Terraform
// StringValue.
func AdaptNullableTimeToStringValue(value *time.Time) basetypes.StringValue {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(value.String())
}

// AdaptSdkModelToResourceObject converts an sdk model to a Terraform resource object.
func AdaptSdkModelToResourceObject[T any, U any](
	sdkModel T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateResourceObject func(sdkModel T) U,
	diags *diag.Diagnostics,
) basetypes.ObjectValue {
	resourceObject := generateResourceObject(sdkModel)

	objectValue, objectDiags := types.ObjectValueFrom(
		ctx,
		attributeTypes,
		resourceObject,
	)
	if objectDiags.HasError() {
		diags.Append(objectDiags...)
		return types.ObjectUnknown(attributeTypes)
	}

	return objectValue
}

// AdaptNullableSdkModelToResourceObject converts a nullable sdk model to a Terraform resource object.
func AdaptNullableSdkModelToResourceObject[T interface{}, U interface{}](
	sdkModel *T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateResourceObject func(sdkModel T) U,
	diags *diag.Diagnostics,
) basetypes.ObjectValue {
	if sdkModel == nil {
		return basetypes.NewObjectNull(attributeTypes)
	}

	resourceObject := generateResourceObject(*sdkModel)

	objectValue, objectDiags := types.ObjectValueFrom(
		ctx,
		attributeTypes,
		resourceObject,
	)
	if objectDiags.HasError() {
		diags.Append(objectDiags...)
		return types.ObjectUnknown(attributeTypes)
	}

	return objectValue
}

// AdaptSdkModelsToListValue converts a sdk model array to a Terraform
// ListValue.
func AdaptSdkModelsToListValue[T any, U any](
	sdkModels []T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateModel func(sdkModel T) U,
	diags *diag.Diagnostics,
) basetypes.ListValue {
	var listValues []U

	for _, value := range sdkModels {
		listValues = append(listValues, generateModel(value))
	}

	listObject, listDiags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: attributeTypes},
		listValues,
	)

	if listDiags.HasError() {
		diags.Append(listDiags...)
		return types.ListUnknown(types.ObjectType{AttrTypes: attributeTypes})
	}

	return listObject
}

// AdaptStringPointerValueToNullableString converts a Terraform
// StringPointerValue to a nullable string.
func AdaptStringPointerValueToNullableString(value types.String) *string {
	if value.IsUnknown() {
		return nil
	}

	return value.ValueStringPointer()
}

func AdaptStringTypeArrayToStringArray[T ~string](types []T) []string {
	var convertedTypes []string

	for _, contractType := range types {
		convertedTypes = append(convertedTypes, string(contractType))
	}

	return convertedTypes
}

// AdaptBoolPointerValueToNullableBool converts a Terraform BoolPointerValue to a nullable string.
func AdaptBoolPointerValueToNullableBool(value types.Bool) *bool {
	if value.IsUnknown() {
		return nil
	}

	return value.ValueBoolPointer()
}
