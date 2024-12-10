package publiccloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &instanceISOResource{}
	_ resource.ResourceWithImportState = &instanceISOResource{}
)

type invalidISOIDError struct {
	supportedISOIDs []string
}

func (u invalidISOIDError) Error() string {
	return fmt.Sprintf(
		"Attribute iso_id value must be one of: %q",
		u.supportedISOIDs,
	)
}

type instanceISOResourceModel struct {
	DesiredISOID types.String `tfsdk:"desired_iso_id"`
	InstanceID   types.String `tfsdk:"instance_id"`
	ISOID        types.String `tfsdk:"iso_id"`
	Name         types.String `tfsdk:"name"`
}

func adaptIsoToInstanceISOResource(
	desiredISOID *string,
	instanceID string,
	sdkISO *publiccloud.Iso,
) instanceISOResourceModel {
	var isoID *string
	var isoName *string
	if sdkISO != nil {
		isoID = &sdkISO.Id
		isoName = &sdkISO.Name
	}

	return instanceISOResourceModel{
		ISOID:        basetypes.NewStringPointerValue(isoID),
		DesiredISOID: basetypes.NewStringPointerValue(desiredISOID),
		InstanceID:   basetypes.NewStringValue(instanceID),
		Name:         basetypes.NewStringPointerValue(isoName),
	}
}

type instanceISOResource struct {
	name   string
	client publiccloud.PubliccloudAPI
}

func (i *instanceISOResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(
		ctx,
		path.Root("instance_id"),
		request,
		response,
	)

	// Set to unknown as there is no desired state when importing
	response.State.SetAttribute(
		ctx,
		path.Root("desired_iso_id"),
		basetypes.NewStringUnknown(),
	)
}

func (i *instanceISOResource) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, i.name)
}

func (i *instanceISOResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: "Creating this resource attaches the ISO to the instance, deleting it detaches the ISO from the instance.",
		Attributes: map[string]schema.Attribute{
			"desired_iso_id": schema.StringAttribute{
				Optional:    true,
				Description: "The desired ISO ID. Removing this will detach the current ISO from the instance. Changing it will cause the current ISO to be detached and a new one to be attached to the instance.",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "Instance's ID. This value cannot be changed.",
			},
			"iso_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO ID. Detaching/attaching iso_ids is an asynchronous operation. This attribute shows the current iso_id while `desired_iso_id` has the desired iso_id state.",
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (i *instanceISOResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan instanceISOResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	state, httpResponse, err := updateISO(plan, i.client, ctx)
	if err != nil {
		var re invalidISOIDError
		ok := errors.As(err, &re)
		if ok {
			response.Diagnostics.AddAttributeError(
				path.Root("desired_iso_id"),
				"Invalid Attribute Value Match",
				re.Error(),
			)
			return
		}

		utils.Error(
			ctx,
			&response.Diagnostics,
			fmt.Sprintf("Creating resource %s", i.name),
			nil,
			httpResponse,
		)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (i *instanceISOResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var currentState instanceISOResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &currentState)...)
	if response.Diagnostics.HasError() {
		return
	}

	instanceDetails, httpResponse, err := i.client.GetInstance(
		ctx,
		currentState.InstanceID.ValueString(),
	).Execute()
	if err != nil {
		utils.Error(
			ctx,
			&response.Diagnostics,
			fmt.Sprintf(
				"Reading ISO %s for instance %q",
				i.name,
				currentState.InstanceID.ValueString(),
			),
			err,
			httpResponse,
		)
		return
	}

	sdkISO, _ := instanceDetails.GetIsoOk()

	// When ISOID is unknown Read is called from ImportState. There is no current state and desired_iso is the same as iso_id
	if currentState.DesiredISOID.IsUnknown() {
		var desiredISOID *string
		if sdkISO != nil {
			desiredISOID = &sdkISO.Id
		}

		state := adaptIsoToInstanceISOResource(
			desiredISOID,
			instanceDetails.Id,
			sdkISO,
		)
		response.Diagnostics.Append(response.State.Set(ctx, state)...)
		return
	}

	// desired_iso_id can be retrieved from the current state
	state := adaptIsoToInstanceISOResource(
		currentState.DesiredISOID.ValueStringPointer(),
		instanceDetails.Id,
		sdkISO,
	)
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

// Update detaches the current ISO and attaches a new one if a new one is set.
func (i *instanceISOResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var currentState instanceISOResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &currentState)...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan instanceISOResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// instance_id is never allowed to change
	if plan.InstanceID.ValueString() != currentState.InstanceID.ValueString() {
		response.Diagnostics.AddAttributeError(
			path.Root("instance_id"),
			"Invalid Attribute Value Match",
			fmt.Sprintf(
				"Attribute instance_id value cannot be changed: was %q",
				currentState.InstanceID.ValueString(),
			),
		)
		return
	}

	if plan.DesiredISOID.ValueString() != currentState.ISOID.ValueString() {
		state, httpResponse, err := updateISO(plan, i.client, ctx)
		if err != nil {
			var re invalidISOIDError
			ok := errors.As(err, &re)
			if ok {
				response.Diagnostics.AddAttributeError(
					path.Root("desired_iso_id"),
					"Invalid Attribute Value Match",
					re.Error(),
				)
				return
			}

			utils.Error(
				ctx,
				&response.Diagnostics,
				fmt.Sprintf(
					"Attaching/detaching ISO %s for instance %q",
					i.name,
					plan.InstanceID.ValueString(),
				),
				err,
				httpResponse,
			)
			return
		}

		response.Diagnostics.Append(response.State.Set(ctx, state)...)
	}
}

// Delete detaches the current ISO.
func (i *instanceISOResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var currentState instanceISOResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &currentState)...)
	if response.Diagnostics.HasError() {
		return
	}

	currentState.DesiredISOID = basetypes.NewStringPointerValue(nil)
	state, httpResponse, err := updateISO(currentState, i.client, ctx)
	if err != nil {
		utils.Error(
			ctx,
			&response.Diagnostics,
			fmt.Sprintf(
				"Detaching ISO %s from instance %q",
				i.name,
				currentState.InstanceID.ValueString(),
			),
			err,
			httpResponse,
		)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (i *instanceISOResource) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
	if request.ProviderData == nil {
		return
	}

	coreClient, ok := request.ProviderData.(client.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				request.ProviderData,
			),
		)
		return
	}

	i.client = coreClient.PubliccloudAPI
}

func NewInstanceIsoResource() resource.Resource {
	return &instanceISOResource{
		name: "public_cloud_instance_iso",
	}
}

// updateISO detaches the current ISO if there is anything attached. If api.desired_api_id is set then a new ISO is attached.
func updateISO(
	iso instanceISOResourceModel,
	api publiccloud.PubliccloudAPI,
	ctx context.Context,
) (*instanceISOResourceModel, *http.Response, error) {
	// If a new ISO is to be attached then check that the ID is valid
	if !iso.DesiredISOID.IsNull() {
		var supportedISOs []publiccloud.Iso
		var offset *int32

		request := api.GetIsoList(ctx)
		for {
			result, httpResponse, err := request.Execute()
			if err != nil {
				return nil, httpResponse, fmt.Errorf("updateISO: %v", err)
			}

			supportedISOs = append(supportedISOs, result.Isos...)

			metadata := result.GetMetadata()
			offset = utils.NewOffset(
				metadata.GetLimit(),
				metadata.GetOffset(),
				metadata.GetTotalCount(),
			)
			if offset == nil {
				break
			}

			request = request.Offset(*offset)
		}

		isValid := false
		for _, supportedISO := range supportedISOs {
			if supportedISO.GetId() == iso.DesiredISOID.ValueString() {
				isValid = true
				break
			}
		}

		if !isValid {
			var supportedISOIDs []string
			for _, supportedISO := range supportedISOs {
				supportedISOIDs = append(supportedISOIDs, supportedISO.GetId())
			}
			return nil, nil, invalidISOIDError{supportedISOIDs: supportedISOIDs}
		}
	}

	instanceDetails, httpResponse, err := api.GetInstance(
		ctx,
		iso.InstanceID.ValueString(),
	).Execute()
	if err != nil {
		return nil, httpResponse, fmt.Errorf("updateISO: %v", err)
	}

	// Detach current ISO if anything is attached
	isoSDK, _ := instanceDetails.GetIsoOk()
	if isoSDK != nil {
		httpResponse, err = api.DetachIso(
			ctx,
			iso.InstanceID.ValueString(),
		).Execute()
		if err != nil {
			return nil, httpResponse, fmt.Errorf("updateISO: %v", err)
		}

		// If a detached ISO is the desired state then exit
		if iso.DesiredISOID.IsNull() {
			instanceDetails, httpResponse, err = api.GetInstance(
				ctx,
				iso.InstanceID.ValueString(),
			).Execute()
			if err != nil {
				return nil, httpResponse, fmt.Errorf("updateISO: %v", err)
			}

			isoSDK, _ = instanceDetails.GetIsoOk()
			updatedISO := adaptIsoToInstanceISOResource(
				nil,
				instanceDetails.Id,
				isoSDK,
			)
			return &updatedISO, nil, nil
		}
	}

	// Attach new ISO
	httpResponse, err = api.AttachIso(
		ctx,
		iso.InstanceID.ValueString(),
	).AttachIsoOpts(*publiccloud.NewAttachIsoOpts(iso.DesiredISOID.ValueString())).
		Execute()
	if err != nil {
		return nil, httpResponse, fmt.Errorf("updateISO: %v", err)
	}

	instanceDetails, httpResponse, err = api.GetInstance(
		ctx,
		iso.InstanceID.ValueString(),
	).Execute()
	if err != nil {
		return nil, httpResponse, fmt.Errorf("updateISO: %v", err)
	}

	isoSDK, _ = instanceDetails.GetIsoOk()
	updatedISO := adaptIsoToInstanceISOResource(
		iso.DesiredISOID.ValueStringPointer(),
		instanceDetails.Id,
		isoSDK,
	)
	return &updatedISO, nil, nil
}
