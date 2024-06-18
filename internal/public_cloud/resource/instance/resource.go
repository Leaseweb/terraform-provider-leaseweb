package instance

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/client"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/validator"
)

var (
	_ resource.Resource                = &instanceResource{}
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
	_ resource.ResourceWithModifyPlan  = &instanceResource{}
)

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

type instanceResource struct {
	client *client.Client
}

func (i *instanceResource) ModifyPlan(
	ctx context.Context,
	request resource.ModifyPlanRequest,
	response *resource.ModifyPlanResponse,
) {
	planInstance := model.Instance{}
	request.Plan.Get(ctx, &planInstance)

	stateInstance := model.Instance{}
	request.State.Get(ctx, &stateInstance)

	typeValidator := validator.NewTypeValidator(
		stateInstance.Id,
		stateInstance.Type,
		planInstance.Type,
	)

	hasTypeChanged := typeValidator.HashTypeChanged()

	if !hasTypeChanged {
		return
	}

	allowedInstanceTypesRequest := i.client.SdkClient.PublicCloudAPI.
		GetUpdateInstanceTypeList(
			i.client.AuthContext(ctx),
			stateInstance.Id.ValueString(),
		)
	allowedInstanceTypes, sdkResponse, err := i.client.SdkClient.PublicCloudAPI.
		GetUpdateInstanceTypeListExecute(allowedInstanceTypesRequest)

	if err != nil {
		if sdkResponse != nil {
			buf := new(strings.Builder)
			_, sdkResponseError := io.Copy(buf, sdkResponse.Body)

			if sdkResponseError == nil {
				tflog.Debug(
					ctx,
					err.Error(),
					map[string]interface{}{"response": buf.String()},
				)
			}
		}

		response.Diagnostics.AddError(
			fmt.Sprintf(
				"Error getting updateInstanceType list for %q",
				stateInstance.Id.ValueString(),
			),
			err.Error(),
		)
		return
	}

	if typeValidator.IsTypeValid(allowedInstanceTypes.GetInstanceTypes()) {
		return
	}

	response.Diagnostics.AddAttributeError(
		path.Root("type"),
		"Invalid Instance Type",
		fmt.Sprintf(
			"Allowed types are %v",
			convertAllowedInstancesTypesToString(allowedInstanceTypes.GetInstanceTypes()),
		),
	)
}

func convertAllowedInstancesTypesToString(updateInstanceTypes []publicCloud.UpdateInstanceType) (instanceTypes []string) {
	for _, instanceType := range updateInstanceTypes {
		instanceTypes = append(instanceTypes, instanceType.GetName())
	}

	return
}
