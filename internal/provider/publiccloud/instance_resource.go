package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
)

type contractResourceModel struct {
	BillingFrequency types.Int32  `tfsdk:"billing_frequency"`
	Term             types.Int32  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func (c contractResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_frequency": types.Int32Type,
		"term":              types.Int32Type,
		"type":              types.StringType,
		"ends_at":           types.StringType,
		"state":             types.StringType,
	}
}

func adaptContractToContractResource(sdkContract publicCloud.Contract) contractResourceModel {
	return contractResourceModel{
		BillingFrequency: basetypes.NewInt32Value(int32(sdkContract.GetBillingFrequency())),
		Term:             basetypes.NewInt32Value(int32(sdkContract.GetTerm())),
		Type:             basetypes.NewStringValue(string(sdkContract.GetType())),
		EndsAt:           utils.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.GetState())),
	}
}

type instanceResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Image               types.Object `tfsdk:"image"`
	State               types.String `tfsdk:"state"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int32  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	IPs                 types.List   `tfsdk:"ips"`
	Contract            types.Object `tfsdk:"contract"`
	MarketAppID         types.String `tfsdk:"market_app_id"`
}

func (i instanceResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"region":    types.StringType,
		"reference": types.StringType,
		"image": types.ObjectType{
			AttrTypes: imageResourceModel{}.AttributeTypes(),
		},
		"state":                  types.StringType,
		"type":                   types.StringType,
		"root_disk_size":         types.Int32Type,
		"root_disk_storage_type": types.StringType,
		"ips": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: iPResourceModel{}.AttributeTypes(),
			},
		},
		"contract": types.ObjectType{
			AttrTypes: contractResourceModel{}.AttributeTypes(),
		},
		"market_app_id": types.StringType,
	}
}

func (i instanceResourceModel) GetLaunchInstanceOpts(ctx context.Context) (
	*publicCloud.LaunchInstanceOpts,
	error,
) {
	sdkRootDiskStorageType, err := publicCloud.NewStorageTypeFromValue(
		i.RootDiskStorageType.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	image := imageResourceModel{}
	imageDiags := i.Image.As(ctx, &image, basetypes.ObjectAsOptions{})
	if imageDiags != nil {
		return nil, utils.ReturnError("GetLaunchInstanceOpts", imageDiags)
	}

	contract := contractResourceModel{}
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
		contract.Term.ValueInt32(),
	)
	if err != nil {
		return nil, err
	}

	sdkBillingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
		contract.BillingFrequency.ValueInt32(),
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

	sdkTypeName, err := publicCloud.NewTypeNameFromValue(
		i.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	opts := publicCloud.NewLaunchInstanceOpts(
		*sdkRegionName,
		*sdkTypeName,
		image.ID.ValueString(),
		*sdkContractType,
		*sdkContractTerm,
		*sdkBillingFrequency,
		*sdkRootDiskStorageType,
	)

	opts.MarketAppId = utils.AdaptStringPointerValueToNullableString(i.MarketAppID)
	opts.Reference = utils.AdaptStringPointerValueToNullableString(i.Reference)
	opts.RootDiskSize = utils.AdaptInt32PointerValueToNullableInt32(i.RootDiskSize)

	return opts, nil
}

func (i instanceResourceModel) GetUpdateInstanceOpts(ctx context.Context) (
	*publicCloud.UpdateInstanceOpts,
	error,
) {
	opts := publicCloud.NewUpdateInstanceOpts()

	opts.Reference = utils.AdaptStringPointerValueToNullableString(i.Reference)
	opts.RootDiskSize = utils.AdaptInt32PointerValueToNullableInt32(i.RootDiskSize)

	contract := contractResourceModel{}
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

	if contract.Term.ValueInt32() != 0 {
		contractTerm, err := publicCloud.NewContractTermFromValue(
			contract.Term.ValueInt32(),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.ContractTerm = contractTerm
	}

	if contract.BillingFrequency.ValueInt32() != 0 {
		billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
			contract.BillingFrequency.ValueInt32(),
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

func adaptInstanceToInstanceResource(
	sdkInstance publicCloud.Instance,
	ctx context.Context,
) (*instanceResourceModel, error) {
	instance := instanceResourceModel{
		ID:                  basetypes.NewStringValue(sdkInstance.GetId()),
		Region:              basetypes.NewStringValue(string(sdkInstance.GetRegion())),
		Reference:           basetypes.NewStringPointerValue(sdkInstance.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstance.GetState())),
		Type:                basetypes.NewStringValue(string(sdkInstance.GetType())),
		RootDiskSize:        basetypes.NewInt32Value(sdkInstance.GetRootDiskSize()),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.GetRootDiskStorageType())),
		MarketAppID:         basetypes.NewStringPointerValue(sdkInstance.MarketAppId.Get()),
	}

	image, err := utils.AdaptSdkModelToResourceObject(
		sdkInstance.Image,
		imageResourceModel{}.AttributeTypes(),
		ctx,
		adaptImageToImageResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.Image = image

	ips, err := utils.AdaptSdkModelsToListValue(
		sdkInstance.Ips,
		iPResourceModel{}.AttributeTypes(),
		ctx,
		adaptIpToIPResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.IPs = ips

	contract, err := utils.AdaptSdkModelToResourceObject(
		sdkInstance.Contract,
		contractResourceModel{}.AttributeTypes(),
		ctx,
		adaptContractToContractResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

func adaptInstanceDetailsToInstanceResource(
	sdkInstanceDetails publicCloud.InstanceDetails,
	ctx context.Context,
) (*instanceResourceModel, error) {
	instance := instanceResourceModel{
		ID:                  basetypes.NewStringValue(sdkInstanceDetails.GetId()),
		Region:              basetypes.NewStringValue(string(sdkInstanceDetails.GetRegion())),
		Reference:           basetypes.NewStringPointerValue(sdkInstanceDetails.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstanceDetails.GetState())),
		Type:                basetypes.NewStringValue(string(sdkInstanceDetails.GetType())),
		RootDiskSize:        basetypes.NewInt32Value(sdkInstanceDetails.GetRootDiskSize()),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstanceDetails.GetRootDiskStorageType())),
		MarketAppID:         basetypes.NewStringPointerValue(sdkInstanceDetails.MarketAppId.Get()),
	}

	image, err := utils.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Image,
		imageResourceModel{}.AttributeTypes(),
		ctx,
		adaptImageToImageResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.Image = image

	ips, err := utils.AdaptSdkModelsToListValue(
		sdkInstanceDetails.Ips,
		iPResourceModel{}.AttributeTypes(),
		ctx,
		adaptIpDetailsToIPResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.IPs = ips

	contract, err := utils.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Contract,
		contractResourceModel{}.AttributeTypes(),
		ctx,
		adaptContractToContractResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

type iPResourceModel struct {
	IP types.String `tfsdk:"ip"`
}

func (i iPResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip": types.StringType,
	}
}

func adaptIpToIPResource(sdkIp publicCloud.Ip) iPResourceModel {
	return iPResourceModel{
		IP: basetypes.NewStringValue(sdkIp.GetIp()),
	}
}

func adaptIpDetailsToIPResource(sdkIpDetails publicCloud.IpDetails) iPResourceModel {
	return iPResourceModel{
		IP: basetypes.NewStringValue(sdkIpDetails.GetIp()),
	}
}

func NewInstanceResource() resource.Resource {
	return &instanceResource{
		name: "public_cloud_instance",
	}
}

type instanceResource struct {
	name   string
	client publicCloud.PublicCloudAPI
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

	i.client = coreClient.PublicCloudAPI
}

func (i *instanceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan instanceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	summary := fmt.Sprintf("Launching resource %s", i.name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Launch Public Cloud instance")

	opts, err := plan.GetLaunchInstanceOpts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, utils.DefaultErrMsg)
		return
	}

	sdkInstance, httpResponse, err := i.client.LaunchInstance(ctx).
		LaunchInstanceOpts(*opts).
		Execute()

	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}

	instance, resourceErr := adaptInstanceToInstanceResource(*sdkInstance, ctx)
	if resourceErr != nil {
		resp.Diagnostics.AddError(summary, utils.DefaultErrMsg)
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
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Terminate Public Cloud instance %q",
		state.ID.ValueString(),
	))
	httpResponse, err := i.client.TerminateInstance(
		ctx,
		state.ID.ValueString(),
	).Execute()

	if err != nil {
		summary := fmt.Sprintf("Terminating resource %s for id %q", i.name, state.ID.ValueString())
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}
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
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, i.name)
}

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state instanceResourceModel

	diags := req.State.Get(ctx, &state)
	summary := fmt.Sprintf(
		"Reading resource %s for id %q",
		i.name,
		state.ID.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(
		ctx,
		fmt.Sprintf("Read Public Cloud instance %q", state.ID.ValueString()),
	)
	sdkInstance, httpResponse, err := i.client.
		GetInstance(ctx, state.ID.ValueString()).
		Execute()
	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}

	tflog.Info(
		ctx,
		fmt.Sprintf(
			"Create publiccloud instance resource for %q",
			state.ID.ValueString(),
		),
	)
	instance, sdkErr := adaptInstanceDetailsToInstanceResource(
		*sdkInstance,
		ctx,
	)
	if sdkErr != nil {
		resp.Diagnostics.AddError(summary, utils.DefaultErrMsg)

		return
	}

	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
}

func (i *instanceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan instanceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	summary := fmt.Sprintf(
		"Updating resource %s for id %q",
		i.name,
		plan.ID.ValueString(),
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(
		ctx,
		fmt.Sprintf("Update Public Cloud instance %q", plan.ID.ValueString()),
	)
	opts, err := plan.GetUpdateInstanceOpts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, utils.DefaultErrMsg)
		return
	}

	sdkInstance, httpResponse, err := i.client.
		UpdateInstance(ctx, plan.ID.ValueString()).
		UpdateInstanceOpts(*opts).
		Execute()
	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}

	diags = resp.State.Set(ctx, sdkInstance)
	resp.Diagnostics.Append(diags...)
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
				Required: true,
				Description: fmt.Sprintf(
					"%s Valid options are %s",
					warningError,
					utils.StringTypeArrayToMarkdown(publicCloud.AllowedRegionNameEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedRegionNameEnumValues)...),
				},
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
					"name": schema.StringAttribute{
						Computed: true,
					},
					"custom": schema.BoolAttribute{
						Computed:    true,
						Description: "Standard or Custom image",
					},
					"state": schema.StringAttribute{
						Computed: true,
					},
					"market_apps": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"storage_types": schema.ListAttribute{
						Computed:    true,
						Description: "The supported storage types for the instance type",
						ElementType: types.StringType,
					},
					"flavour": schema.StringAttribute{
						Computed: true,
					},
					"region": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The instance's current state",
			},
			"type": schema.StringAttribute{
				Required: true,
				Description: fmt.Sprintf(
					"%s Valid options are %s",
					warningError,
					utils.StringTypeArrayToMarkdown(publicCloud.AllowedTypeNameEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.Expressions{path.MatchRoot("region")}...,
					),
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedTypeNameEnumValues)...),
				},
			},
			"root_disk_size": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances. The maximum size is 1000 GB",
				Validators: []validator.Int32{
					int32validator.Between(5, 1000),
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
					"billing_frequency": schema.Int32Attribute{
						Required:    true,
						Description: "The billing frequency (in months). Valid options are " + billingFrequencies.Markdown(),
						Validators: []validator.Int32{
							int32validator.OneOf(billingFrequencies.ToInt32()...),
						},
					},
					"term": schema.Int32Attribute{
						Required:    true,
						Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + contractTerms.Markdown(),
						Validators: []validator.Int32{
							int32validator.OneOf(contractTerms.ToInt32()...),
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
