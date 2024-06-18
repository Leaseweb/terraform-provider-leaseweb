package validator

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type TypeValidator struct {
	stateInstanceId   types.String
	stateInstanceType types.String
	planInstanceType  types.String
}

func NewTypeValidator(
	stateInstanceId types.String,
	stateInstanceType types.String,
	planInstanceType types.String,
) TypeValidator {
	return TypeValidator{
		stateInstanceId:   stateInstanceId,
		stateInstanceType: stateInstanceType,
		planInstanceType:  planInstanceType,
	}
}
func (v TypeValidator) HashTypeChanged() bool {
	// There is nothing to check when creating
	if v.stateInstanceId.ValueString() == "" {
		return false
	}

	// There is nothing to check when importing
	if v.planInstanceType.ValueString() == "" {
		return false
	}

	// Nothing to validate if nothing changes
	if v.planInstanceType.ValueString() == v.stateInstanceType.ValueString() {
		return false
	}

	return true
}

func (v TypeValidator) IsTypeValid(
	allowedInstanceTypes []publicCloud.UpdateInstanceType,
) bool {
	for _, allowedInstanceType := range allowedInstanceTypes {
		if allowedInstanceType.GetName() == v.planInstanceType.ValueString() {
			return true
		}
	}

	return false
}
