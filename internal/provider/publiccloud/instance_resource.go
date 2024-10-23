package publiccloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource                = &instanceResource{}
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
	_ resource.ResourceWithModifyPlan  = &instanceResource{}
)

type reason string

const (
	reasonContractTermCannotBeZero reason = "contract.term cannot be 0 when contract type is MONTHLY"
	reasonContractTermMustBeZero   reason = "contract.term must be 0 when contract type is HOURLY"
	reasonNone                     reason = ""
)

type resourceModelContract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func (c resourceModelContract) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_frequency": types.Int64Type,
		"term":              types.Int64Type,
		"type":              types.StringType,
		"ends_at":           types.StringType,
		"state":             types.StringType,
	}
}

func (c resourceModelContract) IsContractTermValid() (bool, reason) {
	if c.Type.ValueString() == string(publicCloud.CONTRACTTYPE_MONTHLY) && c.Term.ValueInt64() == 0 {
		return false, reasonContractTermCannotBeZero
	}

	if c.Type.ValueString() == string(publicCloud.CONTRACTTYPE_HOURLY) && c.Term.ValueInt64() != 0 {
		return false, reasonContractTermMustBeZero
	}

	return true, reasonNone
}

func adaptSdkContractToResourceContract(sdkContract publicCloud.Contract) resourceModelContract {
	return resourceModelContract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.GetBillingFrequency())),
		Term:             basetypes.NewInt64Value(int64(sdkContract.GetTerm())),
		Type:             basetypes.NewStringValue(string(sdkContract.GetType())),
		EndsAt:           utils.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.GetState())),
	}
}

type resourceModelImage struct {
	ID types.String `tfsdk:"id"`
}

func (i resourceModelImage) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": types.StringType,
	}
}

func adaptSdkImageToResourceImage(sdkImage publicCloud.Image) resourceModelImage {
	return resourceModelImage{
		ID: basetypes.NewStringValue(sdkImage.GetId()),
	}
}

type reasonInstanceCannotBeTerminated string

type resourceModelInstance struct {
	ID                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Image               types.Object `tfsdk:"image"`
	State               types.String `tfsdk:"state"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int64  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	IPs                 types.List   `tfsdk:"ips"`
	Contract            types.Object `tfsdk:"contract"`
	MarketAppID         types.String `tfsdk:"market_app_id"`
}

func (i resourceModelInstance) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"region":    types.StringType,
		"reference": types.StringType,
		"image": types.ObjectType{
			AttrTypes: resourceModelImage{}.AttributeTypes(),
		},
		"state":                  types.StringType,
		"type":                   types.StringType,
		"root_disk_size":         types.Int64Type,
		"root_disk_storage_type": types.StringType,
		"ips": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: resourceModelIP{}.AttributeTypes(),
			},
		},
		"contract": types.ObjectType{
			AttrTypes: resourceModelContract{}.AttributeTypes(),
		},
		"market_app_id": types.StringType,
	}
}

func (i resourceModelInstance) GetLaunchInstanceOpts(ctx context.Context) (
	*publicCloud.LaunchInstanceOpts,
	error,
) {
	sdkRootDiskStorageType, err := publicCloud.NewStorageTypeFromValue(
		i.RootDiskStorageType.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	image := resourceModelImage{}
	imageDiags := i.Image.As(ctx, &image, basetypes.ObjectAsOptions{})
	if imageDiags != nil {
		return nil, utils.ReturnError("GetLaunchInstanceOpts", imageDiags)
	}

	contract := resourceModelContract{}
	contractDiags := i.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		return nil, utils.ReturnError("GetLaunchInstanceOpts", contractDiags)
	}

	sdkContractType, err := publicCloud.NewContractTypeFromValue(
		contract.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	sdkContractTerm, err := publicCloud.NewContractTermFromValue(
		int32(contract.Term.ValueInt64()),
	)
	if err != nil {
		return nil, err
	}

	sdkBillingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
		int32(contract.BillingFrequency.ValueInt64()),
	)
	if err != nil {
		return nil, err
	}

	sdkRegionName, err := publicCloud.NewRegionNameFromValue(
		i.Region.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	sdkInstanceType, err := publicCloud.NewTypeNameFromValue(
		i.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	opts := publicCloud.NewLaunchInstanceOpts(
		*sdkRegionName,
		*sdkInstanceType,
		image.ID.ValueString(),
		*sdkContractType,
		*sdkContractTerm,
		*sdkBillingFrequency,
		*sdkRootDiskStorageType,
	)

	opts.MarketAppId = utils.AdaptStringPointerValueToNullableString(i.MarketAppID)
	opts.Reference = utils.AdaptStringPointerValueToNullableString(i.Reference)
	opts.RootDiskSize = utils.AdaptInt64PointerValueToNullableInt32(i.RootDiskSize)

	return opts, nil
}

func (i resourceModelInstance) GetUpdateInstanceOpts(ctx context.Context) (
	*publicCloud.UpdateInstanceOpts,
	error,
) {
	opts := publicCloud.NewUpdateInstanceOpts()

	opts.Reference = utils.AdaptStringPointerValueToNullableString(i.Reference)
	opts.RootDiskSize = utils.AdaptInt64PointerValueToNullableInt32(i.RootDiskSize)

	contract := resourceModelContract{}
	diags := i.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if diags.HasError() {
		return nil, utils.ReturnError("GetUpdateInstanceOpts", diags)
	}

	if contract.Type.ValueString() != "" {
		contractType, err := publicCloud.NewContractTypeFromValue(
			contract.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.ContractType = contractType
	}

	if contract.Term.ValueInt64() != 0 {
		contractTerm, err := publicCloud.NewContractTermFromValue(
			int32(contract.Term.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.ContractTerm = contractTerm
	}

	if contract.BillingFrequency.ValueInt64() != 0 {
		billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
			int32(contract.BillingFrequency.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.BillingFrequency = billingFrequency
	}

	if i.Type.ValueString() != "" {
		instanceType, err := publicCloud.NewTypeNameFromValue(
			i.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.Type = instanceType
	}

	return opts, nil
}

func (i resourceModelInstance) CanBeTerminated(ctx context.Context) *reasonInstanceCannotBeTerminated {
	contract := resourceModelContract{}
	contractDiags := i.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if contractDiags != nil {
		log.Fatal("cannot convert contract objectType to model")
	}

	if i.State.ValueString() == string(publicCloud.STATE_CREATING) || i.State.ValueString() == string(publicCloud.STATE_DESTROYING) || i.State.ValueString() == string(publicCloud.STATE_DESTROYED) {
		reason := reasonInstanceCannotBeTerminated(
			fmt.Sprintf("state is %q", i.State),
		)

		return &reason
	}

	if !contract.EndsAt.IsNull() {
		reason := reasonInstanceCannotBeTerminated(
			fmt.Sprintf("contract.endsAt is %q", contract.EndsAt.ValueString()),
		)

		return &reason
	}

	return nil
}

func adaptSdkInstanceToResourceInstance(
	sdkInstance publicCloud.Instance,
	ctx context.Context,
) (*resourceModelInstance, error) {
	instance := resourceModelInstance{
		ID:                  basetypes.NewStringValue(sdkInstance.GetId()),
		Region:              basetypes.NewStringValue(string(sdkInstance.GetRegion())),
		Reference:           basetypes.NewStringPointerValue(sdkInstance.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstance.GetState())),
		Type:                basetypes.NewStringValue(string(sdkInstance.GetType())),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.GetRootDiskSize())),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.GetRootDiskStorageType())),
		MarketAppID:         basetypes.NewStringPointerValue(sdkInstance.MarketAppId.Get()),
	}

	image, err := utils.AdaptSdkModelToResourceObject(
		sdkInstance.Image,
		resourceModelImage{}.AttributeTypes(),
		ctx,
		adaptSdkImageToResourceImage,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Image = image

	ips, err := utils.AdaptSdkModelsToListValue(
		sdkInstance.Ips,
		resourceModelIP{}.AttributeTypes(),
		ctx,
		adaptSdkIpToResourceIP,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.IPs = ips

	contract, err := utils.AdaptSdkModelToResourceObject(
		sdkInstance.Contract,
		resourceModelContract{}.AttributeTypes(),
		ctx,
		adaptSdkContractToResourceContract,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

func adaptSdkInstanceDetailsToResourceInstance(
	sdkInstanceDetails publicCloud.InstanceDetails,
	ctx context.Context,
) (*resourceModelInstance, error) {
	instance := resourceModelInstance{
		ID:                  basetypes.NewStringValue(sdkInstanceDetails.GetId()),
		Region:              basetypes.NewStringValue(string(sdkInstanceDetails.GetRegion())),
		Reference:           basetypes.NewStringPointerValue(sdkInstanceDetails.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstanceDetails.GetState())),
		Type:                basetypes.NewStringValue(string(sdkInstanceDetails.GetType())),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstanceDetails.GetRootDiskSize())),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstanceDetails.GetRootDiskStorageType())),
		MarketAppID:         basetypes.NewStringPointerValue(sdkInstanceDetails.MarketAppId.Get()),
	}

	image, err := utils.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Image,
		resourceModelImage{}.AttributeTypes(),
		ctx,
		adaptSdkImageToResourceImage,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Image = image

	ips, err := utils.AdaptSdkModelsToListValue(
		sdkInstanceDetails.Ips,
		resourceModelIP{}.AttributeTypes(),
		ctx,
		adaptSdkIpDetailsToResourceIP,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.IPs = ips

	contract, err := utils.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Contract,
		resourceModelContract{}.AttributeTypes(),
		ctx,
		adaptSdkContractToResourceContract,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

type resourceModelIP struct {
	IP types.String `tfsdk:"ip"`
}

func (i resourceModelIP) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip": types.StringType,
	}
}

func adaptSdkIpToResourceIP(
	_ context.Context,
	sdkIp publicCloud.Ip,
) (*resourceModelIP, error) {
	return &resourceModelIP{
		IP: basetypes.NewStringValue(sdkIp.Ip),
	}, nil
}

func adaptSdkIpDetailsToResourceIP(
	_ context.Context,
	sdkIpDetails publicCloud.IpDetails,
) (*resourceModelIP, error) {
	return &resourceModelIP{
		IP: basetypes.NewStringValue(sdkIpDetails.Ip),
	}, nil
}

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

type instanceResource struct {
	client client.Client
}

func (i *instanceResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	coreClient, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	i.client = coreClient
}

func (i *instanceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan resourceModelInstance

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Launch public cloud instance on API")

	opts, err := plan.GetLaunchInstanceOpts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating launch instance opts",
			err.Error(),
		)

		return
	}

	sdkInstance, repositoryErr := launchInstance(
		*opts,
		ctx,
		i.client.PublicCloudAPI,
	)
	if repositoryErr != nil {
		resp.Diagnostics.AddError(
			"Error creating resourceModelInstance",
			repositoryErr.Error(),
		)

		utils.LogError(
			ctx,
			repositoryErr.ErrorResponse,
			&resp.Diagnostics,
			"Error launching public cloud instance",
			repositoryErr.Error(),
		)

		return
	}

	instance, resourceErr := adaptSdkInstanceToResourceInstance(*sdkInstance, ctx)
	if resourceErr != nil {
		resp.Diagnostics.AddError(
			"Error creating public cloud instance resource",
			resourceErr.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (i *instanceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state resourceModelInstance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Terminate public cloud instance %q",
		state.ID.ValueString(),
	))
	err := terminateInstance(state.ID.ValueString(), ctx, i.client.PublicCloudAPI)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error terminating Public Cloud resourceModelInstance",
			fmt.Sprintf(
				"Could not terminate Public Cloud resourceModelInstance, unexpected error: %q",
				err.Error(),
			),
		)

		utils.LogError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Error terminating public cloud instance %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)

		return
	}
}

func getInstance(
	id string,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.InstanceDetails, *utils.SdkError) {
	instance, response, err := api.GetInstance(ctx, id).Execute()

	if err != nil {
		return nil, utils.NewSdkError(
			fmt.Sprintf("getInstance %q", id),
			err,
			response,
		)
	}

	return instance, nil
}

func launchInstance(
	opts publicCloud.LaunchInstanceOpts,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.Instance, *utils.SdkError) {
	instance, response, err := api.LaunchInstance(ctx).LaunchInstanceOpts(opts).Execute()

	if err != nil {
		return nil, utils.NewSdkError("launchInstance", err, response)
	}

	return instance, nil
}

func updateInstance(
	id string,
	opts publicCloud.UpdateInstanceOpts,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) (*publicCloud.InstanceDetails, *utils.SdkError) {
	instance, response, err := api.UpdateInstance(
		ctx,
		id,
	).UpdateInstanceOpts(opts).Execute()
	if err != nil {
		return nil, utils.NewSdkError(
			fmt.Sprintf("updateInstance %q", id),
			err,
			response,
		)
	}

	return instance, nil
}

func terminateInstance(
	id string,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) *utils.SdkError {
	response, err := api.TerminateInstance(ctx, id).Execute()
	if err != nil {
		return utils.NewSdkError(
			fmt.Sprintf("terminateInstance %q", id),
			err,
			response,
		)
	}

	return nil
}

func getAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) ([]string, *utils.SdkError) {
	var instanceTypes []string

	sdkInstanceTypes, response, err := api.GetUpdateInstanceTypeList(ctx, id).
		Execute()
	if err != nil {
		return nil, utils.NewSdkError(
			fmt.Sprintf("getAvailableInstanceTypesForUpdate %q", id),
			err,
			response,
		)
	}

	for _, sdkInstanceType := range sdkInstanceTypes.InstanceTypes {
		instanceTypes = append(instanceTypes, string(sdkInstanceType.Name))
	}

	return instanceTypes, nil
}

func getRegions(
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) ([]string, *utils.SdkError) {
	var regions []string

	request := api.GetRegionList(ctx)

	result, response, err := request.Execute()

	if err != nil {
		return nil, utils.NewSdkError("getRegions", err, response)
	}

	metadata := result.GetMetadata()
	pagination := utils.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, utils.NewSdkError("getRegions", err, response)
		}

		for _, sdkRegion := range result.Regions {
			regions = append(regions, string(sdkRegion.Name))
		}

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, utils.NewSdkError("GetAllInstances", err, response)
		}
	}

	return regions, nil
}

func getInstanceTypesForRegion(
	region string,
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) ([]string, *utils.SdkError) {
	var instanceTypes []string

	request := api.GetInstanceTypeList(ctx).Region(publicCloud.RegionName(region))

	result, response, err := request.Execute()

	if err != nil {
		return nil, utils.NewSdkError(
			"GetInstanceTypesForRegion",
			err,
			response,
		)
	}

	metadata := result.GetMetadata()
	pagination := utils.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, utils.NewSdkError(
				"GetInstanceTypesForRegion",
				err,
				response,
			)
		}

		for _, sdkInstanceType := range result.InstanceTypes {
			instanceTypes = append(instanceTypes, string(sdkInstanceType.Name))
		}

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, utils.NewSdkError("GetAllInstances", err, response)
		}
	}

	return instanceTypes, nil
}

func (i *instanceResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {

	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		req,
		resp,
	)
}

func (i *instanceResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_instance"
}

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state resourceModelInstance
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Read public cloud instance %q",
		state.ID.ValueString(),
	))
	sdkInstance, err := getInstance(
		state.ID.ValueString(),
		ctx,
		i.client.PublicCloudAPI,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading resourceModelInstance", err.Error())

		utils.LogError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf("Unable to read resourceModelInstance %q", state.ID.ValueString()),
			err.Error(),
		)

		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Create public cloud instance resource for %q",
		state.ID.ValueString(),
	))
	instance, resourceErr := adaptSdkInstanceDetailsToResourceInstance(
		*sdkInstance,
		ctx,
	)
	if resourceErr != nil {
		resp.Diagnostics.AddError(
			"Error creating public cloud instance resource",
			resourceErr.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (i *instanceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan resourceModelInstance

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Update public cloud instance %q",
		plan.ID.ValueString(),
	))
	opts, err := plan.GetUpdateInstanceOpts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating UpdateInstanceOpts",
			err.Error(),
		)
		return
	}

	sdkInstance, repositoryErr := updateInstance(
		plan.ID.ValueString(),
		*opts,
		ctx,
		i.client.PublicCloudAPI,
	)
	if repositoryErr != nil {
		resp.Diagnostics.AddError(
			"Error updating instance",
			repositoryErr.Error(),
		)

		utils.LogError(
			ctx,
			repositoryErr.ErrorResponse,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Unable to update public cloud instance %q",
				plan.ID.ValueString(),
			),
			repositoryErr.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, sdkInstance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (i *instanceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := utils.NewIntMarkdownList(
		append(
			[]publicCloud.BillingFrequency{0},
			publicCloud.AllowedBillingFrequencyEnumValues...,
		),
	)
	contractTerms := utils.NewIntMarkdownList(publicCloud.AllowedContractTermEnumValues)
	warningError := "**WARNING!** Changing this value once running will cause this instance to be destroyed and a new one to be created."

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The instance unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Our current regions can be found in the [developer documentation](https://developer.leaseweb.com/api-docs/publiccloud_v1.html#tag/Instances/operation/launchInstance)" + warningError,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reference": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The identifying name set to the instance",
			},
			"image": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Can be either an Operating System or a UUID in case of a Custom Image ID." + warningError,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The instance's current state",
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.Expressions{path.MatchRoot("region")}...,
					),
				},
			},
			"root_disk_size": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances. The maximum size is 1000 GB",
				Validators: []validator.Int64{
					int64validator.Between(5, 1000),
				},
			},
			"root_disk_storage_type": schema.StringAttribute{
				Required:    true,
				Description: "The root disk's storage type. Can be *LOCAL* or *CENTRAL*. " + warningError,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedStorageTypeEnumValues)...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{Computed: true},
					},
				},
			},
			"contract": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"billing_frequency": schema.Int64Attribute{
						Required:    true,
						Description: "The billing frequency (in months). Valid options are " + billingFrequencies.Markdown(),
						Validators: []validator.Int64{
							int64validator.OneOf(billingFrequencies.ToInt64()...),
						},
					},
					"term": schema.Int64Attribute{
						Required:    true,
						Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + contractTerms.Markdown(),
						Validators: []validator.Int64{
							int64validator.OneOf(contractTerms.ToInt64()...),
						},
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
						Validators: []validator.String{
							stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedContractTypeEnumValues)...),
						},
					},
					"ends_at": schema.StringAttribute{Computed: true},
					"state": schema.StringAttribute{
						Computed: true,
					},
				},
				Validators: []validator.Object{contractTermValidator{}},
			},
			"market_app_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Market App ID that must be installed into the instance." + warningError,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

// ModifyPlan calls validators that require access to the handler.
// This needs to be done here as client.Client isn't properly initialized when
// the schema is called.
func (i *instanceResource) ModifyPlan(
	ctx context.Context,
	request resource.ModifyPlanRequest,
	response *resource.ModifyPlanResponse,
) {
	planInstance := resourceModelInstance{}
	request.Plan.Get(ctx, &planInstance)

	planImage := resourceModelImage{}
	planInstance.Image.As(ctx, &planImage, basetypes.ObjectAsOptions{})

	stateInstance := resourceModelInstance{}
	request.State.Get(ctx, &stateInstance)

	stateImage := resourceModelImage{}
	stateInstance.Image.As(ctx, &stateImage, basetypes.ObjectAsOptions{})

	// Before deletion, determine if the instance is allowed to be deleted
	if request.Plan.Raw.IsNull() {
		i.validateInstance(stateInstance, response, ctx)
		if response.Diagnostics.HasError() {
			return
		}
	}

	regions, err := getRegions(ctx, i.client.PublicCloudAPI)
	if err != nil {
		response.Diagnostics.AddError("Cannot get regions", err.Error())
		return
	}

	// The Region has
	//to be validated first or getAvailableInstanceTypes will throw an error on creation,
	//as the region could be invalid.
	i.validateRegion(planInstance.Region, response, regions, ctx)
	if response.Diagnostics.HasError() {
		return
	}

	availableInstanceTypes := i.getAvailableInstanceTypes(
		response,
		stateInstance.ID,
		planInstance.Region.ValueString(),
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	i.validateInstanceType(
		planInstance.Type,
		stateInstance.Type,
		response,
		availableInstanceTypes,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}
}

// When creating a new resourceModelInstance,
// any instanceType available in the region is good.
// On update, the criteria is more limited.
func (i *instanceResource) getAvailableInstanceTypes(
	response *resource.ModifyPlanResponse,
	id basetypes.StringValue,
	region string,
	ctx context.Context,
) []string {
	// resourceModelInstance is being created.
	if id.IsNull() {
		availableInstanceTypes, err := getInstanceTypesForRegion(
			region,
			ctx,
			i.client.PublicCloudAPI,
		)
		if err != nil {
			response.Diagnostics.AddError(
				"Cannot get available instanceTypes for region",
				err.Error(),
			)
			return nil
		}

		return availableInstanceTypes
	}

	availableInstanceTypes, err := getAvailableInstanceTypesForUpdate(
		id.ValueString(),
		ctx,
		i.client.PublicCloudAPI,
	)
	if err != nil {
		response.Diagnostics.AddError(
			"Cannot get available instanceTypes for update",
			err.Error(),
		)
		return nil
	}

	return availableInstanceTypes
}

func (i *instanceResource) validateRegion(
	plannedValue types.String,
	response *resource.ModifyPlanResponse,
	regions []string,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: plannedValue}
	regionResponse := validator.StringResponse{}

	regionValidator := regionValidator{
		regions: regions,
	}
	regionValidator.ValidateString(ctx, request, &regionResponse)
	if regionResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(regionResponse.Diagnostics.Errors()...)
	}
}

func (i *instanceResource) validateInstanceType(
	instanceType types.String,
	currentInstanceType types.String,
	response *resource.ModifyPlanResponse,
	availableInstanceTypes []string,
	ctx context.Context,
) {
	request := validator.StringRequest{ConfigValue: instanceType}
	instanceResponse := validator.StringResponse{}

	instanceTypeValidator := newInstanceTypeValidator(
		currentInstanceType,
		availableInstanceTypes,
	)

	instanceTypeValidator.ValidateString(ctx, request, &instanceResponse)
	if instanceResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(instanceResponse.Diagnostics.Errors()...)
	}
}

// Checks if instance can be deleted.
func (i *instanceResource) validateInstance(
	instance resourceModelInstance,
	response *resource.ModifyPlanResponse,
	ctx context.Context,
) {
	instanceObject, diags := basetypes.NewObjectValueFrom(
		ctx,
		resourceModelInstance{}.AttributeTypes(),
		instance,
	)
	if diags.HasError() {
		response.Diagnostics.Append(diags.Errors()...)
		return
	}

	instanceRequest := validator.ObjectRequest{ConfigValue: instanceObject}
	instanceResponse := validator.ObjectResponse{}
	validateInstanceTermination := instanceTerminationValidator{}
	validateInstanceTermination.ValidateObject(
		ctx,
		instanceRequest,
		&instanceResponse,
	)

	if instanceResponse.Diagnostics.HasError() {
		response.Diagnostics.Append(instanceResponse.Diagnostics.Errors()...)
	}
}
