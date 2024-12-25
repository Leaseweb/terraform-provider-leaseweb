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
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
)

type isoResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type contractResourceModel struct {
	BillingFrequency types.Int32  `tfsdk:"billing_frequency"`
	Term             types.Int32  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func (c contractResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_frequency": types.Int32Type,
		"term":              types.Int32Type,
		"type":              types.StringType,
		"ends_at":           types.StringType,
		"state":             types.StringType,
	}
}

func adaptContractToContractResource(contract publiccloud.Contract) contractResourceModel {
	return contractResourceModel{
		BillingFrequency: basetypes.NewInt32Value(int32(contract.GetBillingFrequency())),
		Term:             basetypes.NewInt32Value(int32(contract.GetTerm())),
		Type:             basetypes.NewStringValue(string(contract.GetType())),
		EndsAt:           utils.AdaptNullableTimeToStringValue(contract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(contract.GetState())),
	}
}

type instanceResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Image               types.Object `tfsdk:"image"`
	ISO                 types.Object `tfsdk:"iso"`
	State               types.String `tfsdk:"state"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int32  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	IPs                 types.List   `tfsdk:"ips"`
	Contract            types.Object `tfsdk:"contract"`
	MarketAppID         types.String `tfsdk:"market_app_id"`
}

func adaptInstanceDetailsToInstanceResource(
	instanceDetails publiccloud.InstanceDetails,
	ctx context.Context,
) (*instanceResourceModel, error) {
	instance := instanceResourceModel{
		ID:                  basetypes.NewStringValue(instanceDetails.GetId()),
		Region:              basetypes.NewStringValue(string(instanceDetails.GetRegion())),
		Reference:           basetypes.NewStringPointerValue(instanceDetails.Reference.Get()),
		State:               basetypes.NewStringValue(string(instanceDetails.GetState())),
		Type:                basetypes.NewStringValue(string(instanceDetails.GetType())),
		RootDiskSize:        basetypes.NewInt32Value(instanceDetails.GetRootDiskSize()),
		RootDiskStorageType: basetypes.NewStringValue(string(instanceDetails.GetRootDiskStorageType())),
		MarketAppID:         basetypes.NewStringPointerValue(instanceDetails.MarketAppId.Get()),
	}

	image, err := utils.AdaptSdkModelToResourceObject(
		instanceDetails.Image,
		map[string]attr.Type{
			"id":            types.StringType,
			"instance_id":   types.StringType,
			"name":          types.StringType,
			"custom":        types.BoolType,
			"state":         types.StringType,
			"market_apps":   types.ListType{ElemType: types.StringType},
			"storage_types": types.ListType{ElemType: types.StringType},
			"flavour":       types.StringType,
			"region":        types.StringType,
		},
		ctx,
		func(image publiccloud.Image) imageResourceModel {
			emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})

			return imageResourceModel{
				ID:           basetypes.NewStringValue(image.GetId()),
				Name:         basetypes.NewStringValue(image.GetName()),
				Custom:       basetypes.NewBoolValue(image.GetCustom()),
				Flavour:      basetypes.NewStringValue(string(image.GetFlavour())),
				MarketApps:   emptyList,
				StorageTypes: emptyList,
			}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.Image = image

	ips, err := utils.AdaptSdkModelsToListValue(
		instanceDetails.Ips,
		map[string]attr.Type{
			"reverse_lookup": types.StringType,
			"instance_id":    types.StringType,
			"ip":             types.StringType,
		},
		ctx,
		adaptIpDetailsToIPResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.IPs = ips

	contract, err := utils.AdaptSdkModelToResourceObject(
		instanceDetails.Contract,
		contractResourceModel{}.attributeTypes(),
		ctx,
		adaptContractToContractResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.Contract = contract

	sdkIso, _ := instanceDetails.GetIsoOk()
	iso, err := utils.AdaptNullableSdkModelToResourceObject(
		sdkIso,
		map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
		ctx,
		func(iso publiccloud.Iso) isoResourceModel {
			return isoResourceModel{
				ID:   basetypes.NewStringValue(iso.GetId()),
				Name: basetypes.NewStringValue(iso.GetName()),
			}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.ISO = iso

	return &instance, nil
}

func NewInstanceResource() resource.Resource {
	return &instanceResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "public_cloud_instance",
		},
	}
}

type instanceResource struct {
	utils.ResourceAPI
}

func (i *instanceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan instanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	image := imageResourceModel{}
	imageDiags := plan.Image.As(ctx, &image, basetypes.ObjectAsOptions{})
	if imageDiags != nil {
		resp.Diagnostics.Append(imageDiags...)
		return
	}

	contract := contractResourceModel{}
	contractDiags := plan.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		resp.Diagnostics.Append(contractDiags...)
		return
	}

	opts := publiccloud.NewLaunchInstanceOpts(
		publiccloud.RegionName(plan.Region.ValueString()),
		publiccloud.TypeName(plan.Type.ValueString()),
		image.ID.ValueString(),
		publiccloud.ContractType(contract.Type.ValueString()),
		publiccloud.ContractTerm(contract.Term.ValueInt32()),
		publiccloud.BillingFrequency(contract.BillingFrequency.ValueInt32()),
		publiccloud.StorageType(plan.RootDiskStorageType.ValueString()),
	)
	opts.MarketAppId = utils.AdaptStringPointerValueToNullableString(plan.MarketAppID)
	opts.Reference = utils.AdaptStringPointerValueToNullableString(plan.Reference)
	opts.RootDiskSize = utils.AdaptInt32PointerValueToNullableInt32(plan.RootDiskSize)

	instance, httpResponse, err := i.PubliccloudAPI.LaunchInstance(ctx).
		LaunchInstanceOpts(*opts).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	}

	// Get ISO data from instanceDetails
	instanceDetails, httpResponse, err := i.PubliccloudAPI.GetInstance(
		ctx,
		instance.GetId(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	}

	state, err := adaptInstanceDetailsToInstanceResource(*instanceDetails, ctx)
	if err != nil {
		utils.GeneralError(&resp.Diagnostics, ctx, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (i *instanceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state instanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := i.PubliccloudAPI.TerminateInstance(
		ctx,
		state.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
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

func (i *instanceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state instanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceDetails, httpResponse, err := i.PubliccloudAPI.
		GetInstance(ctx, state.ID.ValueString()).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	}

	newState, err := adaptInstanceDetailsToInstanceResource(
		*instanceDetails,
		ctx,
	)
	if err != nil {
		utils.GeneralError(&resp.Diagnostics, ctx, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (i *instanceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan instanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewUpdateInstanceOpts()
	opts.Reference = utils.AdaptStringPointerValueToNullableString(plan.Reference)
	opts.RootDiskSize = utils.AdaptInt32PointerValueToNullableInt32(plan.RootDiskSize)
	contract := contractResourceModel{}
	diags := plan.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if contract.Type.ValueString() != "" {
		opts.SetContractType(publiccloud.ContractType(contract.Type.ValueString()))
	}
	if contract.Term.ValueInt32() != 0 {
		opts.SetContractTerm(publiccloud.ContractTerm(contract.Term.ValueInt32()))
	}
	if contract.BillingFrequency.ValueInt32() != 0 {
		opts.SetBillingFrequency(publiccloud.BillingFrequency(contract.BillingFrequency.ValueInt32()))
	}
	if plan.Type.ValueString() != "" {
		opts.SetType(publiccloud.TypeName(plan.Type.ValueString()))
	}

	instanceDetails, httpResponse, err := i.PubliccloudAPI.
		UpdateInstance(ctx, plan.ID.ValueString()).
		UpdateInstanceOpts(*opts).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
		return
	}

	state, err := adaptInstanceDetailsToInstanceResource(*instanceDetails, ctx)
	if err != nil {
		utils.GeneralError(&resp.Diagnostics, ctx, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (i *instanceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := utils.NewIntMarkdownList(
		append(
			[]publiccloud.BillingFrequency{0},
			publiccloud.AllowedBillingFrequencyEnumValues...,
		),
	)
	contractTerms := utils.NewIntMarkdownList(publiccloud.AllowedContractTermEnumValues)
	warningError := "**WARNING!** Changing this value once running will cause this instance to be destroyed and a new one to be created."

	resp.Schema = schema.Schema{
		Description: utils.BetaDescription,
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
					utils.StringTypeArrayToMarkdown(publiccloud.AllowedRegionNameEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedRegionNameEnumValues)...),
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
						Description: "Can be either an Operating System or a UUID in case of a Custom Image ID. " + warningError,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"instance_id": schema.StringAttribute{
						Computed: true,
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
			"iso": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ISO ID.",
					},
					"name": schema.StringAttribute{
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
					utils.StringTypeArrayToMarkdown(publiccloud.AllowedTypeNameEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.Expressions{path.MatchRoot("region")}...,
					),
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedTypeNameEnumValues)...),
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
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedStorageTypeEnumValues)...),
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
						"instance_id": schema.StringAttribute{
							Computed: true,
						},
						"reverse_lookup": schema.StringAttribute{
							Computed: true,
						},
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
							stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedContractTypeEnumValues)...),
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
				Description: "Market App ID that must be installed into the instance. " + warningError,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}
