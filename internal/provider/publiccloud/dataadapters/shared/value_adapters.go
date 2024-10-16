package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

// AdaptInt64PointerValueToNullableInt32 converts a Terraform
// Int64PointerValue to a nullable int32.
func AdaptInt64PointerValueToNullableInt32(int64Type types.Int64) *int32 {
	if int64Type.IsUnknown() {
		return nil
	}

	value := int64Type.ValueInt64Pointer()
	if value == nil {
		return nil
	}

	convertedValue := int32(*value)

	return &convertedValue
}

// AdaptNullableTimeToStringValue converts a nullable Time to a Terraform
// StringValue.
func AdaptNullableTimeToStringValue(value *time.Time) basetypes.StringValue {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(value.String())
}

// AdaptNullableStringToStringValue converts a nullable string to a Terraform
// StringValue.
func AdaptNullableStringToStringValue(value *string) basetypes.StringValue {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(*value)
}

// AdaptSdkModelToResourceObject converts an sdk model to a Terraform resource object.
func AdaptSdkModelToResourceObject[T any, U any](
	sdkModel T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateResourceObject func(
		ctx context.Context,
		sdkModel T,
	) (*U, error),
) (basetypes.ObjectValue, error) {
	resourceObject, err := generateResourceObject(ctx, sdkModel)
	if err != nil {
		return types.ObjectUnknown(attributeTypes), fmt.Errorf(
			"unable to convert sdk sdkModel to resource: %w",
			err,
		)
	}

	objectValue, diags := types.ObjectValueFrom(
		ctx,
		attributeTypes,
		resourceObject,
	)
	if diags.HasError() {
		for _, v := range diags {
			return types.ObjectUnknown(attributeTypes), fmt.Errorf(
				"unable to convert sdk sdkModel to resource: %q %q",
				v.Summary(),
				v.Detail(),
			)
		}

	}

	return objectValue, nil
}

// AdaptSdkModelsToListValue converts a sdk model array to a Terraform
// ListValue.
func AdaptSdkModelsToListValue[T any, U any](
	sdkModels []T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateModel func(
		ctx context.Context,
		sdkModel T,
	) (*U, error),
) (basetypes.ListValue, error) {
	var listValues []U

	for _, value := range sdkModels {
		resourceObject, err := generateModel(ctx, value)
		if err != nil {
			return types.ListUnknown(
					types.ObjectType{AttrTypes: attributeTypes}),
				fmt.Errorf(
					"unable to convert sdk model to resource: %w",
					err,
				)
		}
		listValues = append(listValues, *resourceObject)
	}

	listObject, diags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: attributeTypes},
		listValues,
	)

	if diags.HasError() {
		for _, v := range diags {
			return types.ListUnknown(
					types.ObjectType{AttrTypes: attributeTypes}), fmt.Errorf(
					"unable to convert sdk model to resource: %q %q",
					v.Summary(),
					v.Detail(),
				)
		}
	}

	return listObject, nil
}

// AdaptStringPointerValueToNullableString converts a Terraform
// StringPointerValue to a nullable string.
func AdaptStringPointerValueToNullableString(value types.String) *string {
	if value.IsUnknown() {
		return nil
	}

	return value.ValueStringPointer()
}

// ReturnError returns the first diagnostics error as a golang Error.
func ReturnError(functionName string, diags diag.Diagnostics) error {
	for _, diagError := range diags {
		return fmt.Errorf(
			"%s: %q %q",
			functionName,
			diagError.Summary(),
			diagError.Detail(),
		)
	}

	return nil
}

func AdaptContractTypesToStringArray(types []publicCloud.ContractType) []string {
	var convertedTypes []string

	for _, contractType := range types {
		convertedTypes = append(convertedTypes, string(contractType))
	}

	return convertedTypes
}
