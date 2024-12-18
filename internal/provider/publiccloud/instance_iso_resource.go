package publiccloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &instanceISOResource{}
	_ resource.ResourceWithImportState = &instanceISOResource{}
)

type invalidIDError struct {
	supportedIDs []string
}

func (u invalidIDError) Error() string {
	return fmt.Sprintf(
		"Attribute id value must be one of: %q",
		u.supportedIDs,
	)
}

type instanceISOResourceModel struct {
	DesiredID  types.String `tfsdk:"desired_id"`
	InstanceID types.String `tfsdk:"instance_id"`
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
}

func adaptIsoToInstanceISOResource(
	desiredID *string,
	instanceID string,
	iso *publiccloud.Iso,
) instanceISOResourceModel {
	var id *string
	var name *string
	if iso != nil {
		id = &iso.Id
		name = &iso.Name
	}

	return instanceISOResourceModel{
		ID:         basetypes.NewStringPointerValue(id),
		DesiredID:  basetypes.NewStringPointerValue(desiredID),
		InstanceID: basetypes.NewStringValue(instanceID),
		Name:       basetypes.NewStringPointerValue(name),
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
		path.Root("desired_id"),
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
			"desired_id": schema.StringAttribute{
				Optional:    true,
				Description: "The desired ISO ID. Removing this will detach the current ISO from the instance. Changing it will cause the current ISO to be detached and a new one to be attached to the instance.",
			},
			"instance_id": schema.StringAttribute{
				Required:      true,
				Description:   "Instance's ID. **WARNING!** Changing instance_id will cause the current ISO to be detached.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO ID. Detaching/attaching ids is an asynchronous operation. This attribute shows the current id while `desired_id` has the desired id state.",
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
		var re invalidIDError
		ok := errors.As(err, &re)
		if ok {
			response.Diagnostics.AddAttributeError(
				path.Root("desired_id"),
				"Invalid Attribute Value Match",
				re.Error(),
			)
			return
		}
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
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
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	iso, _ := instanceDetails.GetIsoOk()

	// When ID is unknown Read is called from ImportState. There is no current state and desired_id is the same as id
	if currentState.DesiredID.IsUnknown() {
		var desiredID *string
		if iso != nil {
			desiredID = &iso.Id
		}

		state := adaptIsoToInstanceISOResource(
			desiredID,
			instanceDetails.Id,
			iso,
		)
		response.Diagnostics.Append(response.State.Set(ctx, state)...)
		return
	}

	// desired_id can be retrieved from the current state
	state := adaptIsoToInstanceISOResource(
		currentState.DesiredID.ValueStringPointer(),
		instanceDetails.Id,
		iso,
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

	if plan.DesiredID.ValueString() != currentState.ID.ValueString() {
		state, httpResponse, err := updateISO(plan, i.client, ctx)
		if err != nil {
			var re invalidIDError
			ok := errors.As(err, &re)
			if ok {
				response.Diagnostics.AddAttributeError(
					path.Root("desired_id"),
					"Invalid Attribute Value Match",
					re.Error(),
				)
				return
			}
			utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
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

	currentState.DesiredID = basetypes.NewStringPointerValue(nil)
	state, httpResponse, err := updateISO(currentState, i.client, ctx)
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
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
		utils.ConfigError(&response.Diagnostics, request.ProviderData)
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
	if !iso.DesiredID.IsNull() {
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
			if supportedISO.GetId() == iso.DesiredID.ValueString() {
				isValid = true
				break
			}
		}

		if !isValid {
			var supportedIDs []string
			for _, supportedISO := range supportedISOs {
				supportedIDs = append(supportedIDs, supportedISO.GetId())
			}
			return nil, nil, invalidIDError{supportedIDs: supportedIDs}
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
		httpResponse, err = api.DetachInstanceISO(
			ctx,
			iso.InstanceID.ValueString(),
		).Execute()
		if err != nil {
			return nil, httpResponse, fmt.Errorf("updateISO: %v", err)
		}

		// If a detached ISO is the desired state then exit
		if iso.DesiredID.IsNull() {
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
	httpResponse, err = api.AttachInstanceISO(
		ctx,
		iso.InstanceID.ValueString(),
	).AttachIsoOpts(*publiccloud.NewAttachIsoOpts(iso.DesiredID.ValueString())).
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
		iso.DesiredID.ValueStringPointer(),
		instanceDetails.Id,
		isoSDK,
	)
	return &updatedISO, nil, nil
}
