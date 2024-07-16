package utils

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ConvertNullableIntToInt64Value Convert NullableInt to terraform Int64.
func ConvertNullableIntToInt64Value(value *int) basetypes.Int64Value {
	if value == nil {
		return basetypes.NewInt64Null()
	}

	return basetypes.NewInt64Value(int64(*value))
}

// ConvertNullableTimeToStringValue Convert NullableTime to terraform String.
func ConvertNullableTimeToStringValue(value *time.Time) basetypes.StringValue {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(value.String())
}

// ConvertNullableStringToStringValue Convert NullableString to terraform String.
func ConvertNullableStringToStringValue(value *string) basetypes.StringValue {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(*value)
}

// ConvertNullableDomainEntityToDatasourceModel Convert nullable domain entity to datasource model.
func ConvertNullableDomainEntityToDatasourceModel[T interface{}, U interface{}](
	value *T,
	generateModel func(sdkEntity T) *U,
) *U {
	if value == nil {
		return nil
	}

	return generateModel(*value)
}

// ConvertNullableDomainEntityToResourceObject Convert nullable domain entity to resource object.
func ConvertNullableDomainEntityToResourceObject[T any, U any](
	entity *T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateResourceObject func(
		ctx context.Context,
		entity T,
	) (
		model *U,
		diagnostics diag.Diagnostics,
	)) (basetypes.ObjectValue, diag.Diagnostics) {
	if entity == nil {
		return types.ObjectNull(attributeTypes), nil
	}

	resourceObject, diags := ConvertDomainEntityToResourceObject(
		*entity,
		attributeTypes,
		ctx,
		generateResourceObject,
	)

	if diags.HasError() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	return resourceObject, nil
}

// ConvertDomainEntityToResourceObject Convert domain entity to resource object.
func ConvertDomainEntityToResourceObject[T any, U any](
	entity T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateResourceObject func(
		ctx context.Context,
		entity T,
	) (model *U, diagnostics diag.Diagnostics),
) (basetypes.ObjectValue, diag.Diagnostics) {
	resourceObject, diags := generateResourceObject(ctx, entity)
	if diags.HasError() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objectValue, diags := types.ObjectValueFrom(
		ctx,
		attributeTypes,
		resourceObject,
	)
	if diags.HasError() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	return objectValue, nil
}

// ConvertEntitiesToListValue Convert a slice of entities to a list value.
func ConvertEntitiesToListValue[T any, U any](
	entities []T,
	attributeTypes map[string]attr.Type,
	ctx context.Context,
	generateTerraformModel func(
		ctx context.Context,
		entity T,
	) (model *U, diagnostics diag.Diagnostics),
) (basetypes.ListValue, diag.Diagnostics) {
	var listValues []U

	for _, value := range entities {
		resourceObject, diags := generateTerraformModel(ctx, value)
		if diags.HasError() {
			return types.ListUnknown(types.ObjectType{AttrTypes: attributeTypes}), diags
		}
		listValues = append(listValues, *resourceObject)
	}

	listObject, diags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: attributeTypes},
		listValues,
	)

	if diags.HasError() {
		return types.ListUnknown(types.ObjectType{AttrTypes: attributeTypes}), diags
	}

	return listObject, nil
}
