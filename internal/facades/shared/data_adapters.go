package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AdaptNullableIntToInt64Value converts a NullableInt toTerraform Int64Value.
func AdaptNullableIntToInt64Value(value *int) basetypes.Int64Value {
	if value == nil {
		return basetypes.NewInt64Null()
	}

	return basetypes.NewInt64Value(int64(*value))
}

// AdaptNullableTimeToStringValue converts a NullableTime toTerraform StringValue.
func AdaptNullableTimeToStringValue(value *time.Time) basetypes.StringValue {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(value.String())
}

// AdaptNullableStringToStringValue converts a NullableString toTerraform StringValue.
func AdaptNullableStringToStringValue(value *string) basetypes.StringValue {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(*value)
}

// AdaptNullableDomainEntityToDatasourceModel converts a nullable domain entity to Terraform datasource model.
func AdaptNullableDomainEntityToDatasourceModel[T interface{}, U interface{}](
	entity *T,
	generateModel func(entity T) *U,
) *U {
	if entity == nil {
		return nil
	}

	return generateModel(*entity)
}

// AdaptNullableDomainEntityToResourceObject converts a nullable domain entity to Terraform resource object.
func AdaptNullableDomainEntityToResourceObject[T any, U any](
	entity *T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateResourceObject func(
		ctx context.Context,
		entity T,
	) (*U, error)) (basetypes.ObjectValue, error) {
	if entity == nil {
		return types.ObjectNull(attributeTypes), nil
	}

	resourceObject, err := AdaptDomainEntityToResourceObject(
		*entity,
		attributeTypes,
		ctx,
		generateResourceObject,
	)

	if err != nil {
		return types.ObjectUnknown(attributeTypes), fmt.Errorf(
			"unable to convert domain entity to resource: %w",
			err,
		)
	}

	return resourceObject, nil
}

// AdaptDomainEntityToResourceObject converts a domain entity to Terraform resource object.
func AdaptDomainEntityToResourceObject[T any, U any](
	entity T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateResourceObject func(
		ctx context.Context,
		entity T,
	) (*U, error),
) (basetypes.ObjectValue, error) {
	resourceObject, err := generateResourceObject(ctx, entity)
	if err != nil {
		return types.ObjectUnknown(attributeTypes), fmt.Errorf(
			"unable to convert domain entity to resource: %w",
			err,
		)
	}

	objectValue, diags := types.ObjectValueFrom(
		ctx,
		attributeTypes,
		resourceObject,
	)
	if diags.HasError() {
		for _, diag := range diags {
			return types.ObjectUnknown(attributeTypes), fmt.Errorf(
				"unable to convert domain entity to resource: %q %q",
				diag.Summary(),
				diag.Detail(),
			)
		}

	}

	return objectValue, nil
}

// AdaptEntitiesToListValue converts a domain entities object to a Terraform list value.
func AdaptEntitiesToListValue[T any, U any](
	entities []T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateModel func(
		ctx context.Context,
		entity T,
	) (*U, error),
) (basetypes.ListValue, error) {
	var listValues []U

	for _, value := range entities {
		resourceObject, err := generateModel(ctx, value)
		if err != nil {
			return types.ListUnknown(types.ObjectType{AttrTypes: attributeTypes}), fmt.Errorf(
				"unable to convert domain entity to resource: %w",
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
		for _, diag := range diags {
			return types.ListUnknown(types.ObjectType{AttrTypes: attributeTypes}), fmt.Errorf(
				"unable to convert domain entity to resource: %q %q",
				diag.Summary(),
				diag.Detail(),
			)
		}
	}

	return listObject, nil
}

// AdaptStringPointerValueToNullableString converts a Terraform StringPointerValue to nullable string.
func AdaptStringPointerValueToNullableString(value types.String) *string {
	if value.IsUnknown() {
		return nil
	}

	return value.ValueStringPointer()
}

// AdaptIntArrayToInt64Array converts an array of integers to an array of int64 values.
func AdaptIntArrayToInt64Array(items []int) []int64 {
	var convertedItems []int64

	for _, item := range items {
		convertedItems = append(
			convertedItems,
			int64(item),
		)
	}

	return convertedItems
}

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
