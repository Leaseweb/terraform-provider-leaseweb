package validator

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = InstanceTypeValidator{}

type InstanceTypeValidator struct {
	isInstanceTypeAvailableForRegion func(
		instanceType string,
		region string,
		ctx context.Context,
	) (bool, []string, error)
	canInstanceTypeBeUsedWithInstance func(
		id string,
		instanceType string,
		ctx context.Context,
	) (bool, []string, error)
	instanceId types.String
	region     types.String
}

func (i InstanceTypeValidator) Description(ctx context.Context) string {
	return "Determines if an instanceType can be used with an instance."
}

func (i InstanceTypeValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i InstanceTypeValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {

	// If the instanceType is unknown, there is nothing to validate.
	if request.ConfigValue.IsUnknown() {
		return
	}

	// Do nothing if instanceType does not change.
	if request.ConfigValue.IsNull() {
		return
	}

	// Instance is being created.
	if i.instanceId.IsNull() {
		i.validateCreatedInstance(request, response, ctx)
		return
	}

	// Instance is being updated.
	i.validateUpdatedInstance(request, response, ctx)
}

func (i InstanceTypeValidator) validateCreatedInstance(
	request validator.StringRequest,
	response *validator.StringResponse,
	ctx context.Context,
) {
	isAvailable, availableInstanceTypes, err := i.isInstanceTypeAvailableForRegion(
		request.ConfigValue.ValueString(),
		i.region.ValueString(),
		ctx,
	)

	if err != nil {
		log.Fatal(err)
	}

	if !isAvailable {
		i.setError(response, request, availableInstanceTypes)
	}
}

func (i InstanceTypeValidator) validateUpdatedInstance(
	request validator.StringRequest,
	response *validator.StringResponse,
	ctx context.Context,
) {
	canInstanceTypeBeUsed, allowedInstanceTypes, err := i.canInstanceTypeBeUsedWithInstance(
		i.instanceId.ValueString(),
		request.ConfigValue.ValueString(),
		ctx,
	)

	if err != nil {
		log.Fatal(err)
	}

	if !canInstanceTypeBeUsed {
		i.setError(response, request, allowedInstanceTypes)
	}
}

func (i InstanceTypeValidator) setError(
	response *validator.StringResponse,
	request validator.StringRequest,
	instanceTypes []string,
) {
	response.Diagnostics.AddAttributeError(
		request.Path,
		"Invalid Instance Type",
		fmt.Sprintf(
			"Attribute type value must be one of: %q, got: %q",
			instanceTypes,
			request.ConfigValue.ValueString(),
		),
	)
}

func NewInstanceTypeValidator(
	isInstanceTypeAvailableForRegion func(
		instanceType string,
		region string,
		ctx context.Context,
	) (bool, []string, error),
	canInstanceTypeBeUsedWithInstance func(
		id string,
		instanceType string,
		ctx context.Context,
	) (bool, []string, error),
	instanceId types.String,
	region types.String,
) InstanceTypeValidator {
	if region.IsUnknown() {
		log.Fatal(errors.New("region must be specified"))
	}

	return InstanceTypeValidator{
		isInstanceTypeAvailableForRegion:  isInstanceTypeAvailableForRegion,
		canInstanceTypeBeUsedWithInstance: canInstanceTypeBeUsedWithInstance,
		instanceId:                        instanceId,
		region:                            region,
	}
}
