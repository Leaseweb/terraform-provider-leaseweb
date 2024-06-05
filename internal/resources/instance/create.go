package instance

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"strconv"
	"strings"
	"terraform-provider-leaseweb/internal/resources/instance/model"
)

func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.Instance
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	term, err := convertStringToInt32("contract term", plan.Contract.Attributes()["term"].String(), resp)
	if err != nil {
		return
	}

	billingFrequency, err := convertStringToInt32("contract billing frequency", plan.Contract.Attributes()["billing_frequency"].String(), resp)
	if err != nil {
		return
	}

	opts := publicCloud.NewLaunchInstanceOpts(
		plan.Region.ValueString(),
		plan.OperatingSystem.Attributes()["id"].String(),
		strings.Trim(plan.Contract.Attributes()["type"].String(), "\""),
		term,
		billingFrequency,
		plan.RootDiskStorageType.ValueString(),
	)

	request := r.client.SdkClient.PublicCloudAPI.LaunchInstance(r.client.AuthContext()).LaunchInstanceOpts(*opts)
	instance, _, err := r.client.SdkClient.PublicCloudAPI.LaunchInstanceExecute(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating instance",
			"Could not create instance, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Populate(instance, ctx)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertStringToInt32(name string, value string, resp *resource.CreateResponse) (int32, error) {

	convertedValue, err := strconv.Atoi(value)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error setting \"%s\"", name),
			fmt.Sprintf("Could not set \"%s\", unexpected error:  \"%s\"", name, err.Error()),
		)
		return 0, err
	}

	return int32(convertedValue), nil
}
