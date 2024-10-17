package validator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	serviceErrors "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service/errors"
)

var _ validator.String = InstanceTypeValidator{}

type InstanceTypeValidator struct {
	isInstanceTypeAvailableForRegion func(
		instanceType string,
		region string,
		ctx context.Context,
	) (bool, []string, *serviceErrors.ServiceError)
	instanceId                      types.String
	region                          types.String
	currentInstanceType             types.String
	availableInstanceTypesForUpdate []string
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
	// Nothing to validate here.
	if request.ConfigValue.IsUnknown() || request.ConfigValue.IsNull() {
		return
	}

	// Instance is being created.
	if i.instanceId.IsNull() {
		i.validateCreatedInstance(request, response, ctx)
		return
	}

	// Instance is being updated.
	i.validateUpdatedInstance(request, response)
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
) {
	if !slices.Contains(i.availableInstanceTypesForUpdate, request.ConfigValue.ValueString()) {
		i.setError(response, request, i.availableInstanceTypesForUpdate)
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
	) (bool, []string, *serviceErrors.ServiceError),
	instanceId types.String,
	region types.String,
	currentInstanceType types.String,
	availableInstanceTypesForUpdate []string,
) InstanceTypeValidator {
	if region.IsUnknown() {
		log.Fatal(errors.New("region must be specified"))
	}

	// Include the current instance type as it isn't returned the by api.
	availableInstanceTypesForUpdate = append(availableInstanceTypesForUpdate, currentInstanceType.ValueString())

	return InstanceTypeValidator{
		isInstanceTypeAvailableForRegion: isInstanceTypeAvailableForRegion,
		instanceId:                       instanceId,
		region:                           region,
		currentInstanceType:              currentInstanceType,
		availableInstanceTypesForUpdate:  availableInstanceTypesForUpdate,
	}
}
