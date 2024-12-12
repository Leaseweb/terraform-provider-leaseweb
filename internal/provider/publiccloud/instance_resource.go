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
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
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

func (i isoResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

func adaptIsoToISOResource(iso publiccloud.Iso) isoResourceModel {
	return isoResourceModel{
		ID:   basetypes.NewStringValue(iso.GetId()),
		Name: basetypes.NewStringValue(iso.GetName()),
	}
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

func (i instanceResourceModel) getLaunchOpts(ctx context.Context) (
	*publiccloud.LaunchInstanceOpts,
	error,
) {
	rootDiskStorageType, err := publiccloud.NewStorageTypeFromValue(
		i.RootDiskStorageType.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	image := imageResourceModel{}
	imageDiags := i.Image.As(ctx, &image, basetypes.ObjectAsOptions{})
	if imageDiags != nil {
		return nil, utils.ReturnError("getLaunchOpts", imageDiags)
	}

	contract := contractResourceModel{}
	contractDiags := i.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		return nil, utils.ReturnError("getLaunchOpts", contractDiags)
	}

	contractType, err := publiccloud.NewContractTypeFromValue(
		contract.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	contractTerm, err := publiccloud.NewContractTermFromValue(
		contract.Term.ValueInt32(),
	)
	if err != nil {
		return nil, err
	}

	billingFrequency, err := publiccloud.NewBillingFrequencyFromValue(
		contract.BillingFrequency.ValueInt32(),
	)
	if err != nil {
		return nil, err
	}

	regionName, err := publiccloud.NewRegionNameFromValue(
		i.Region.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	typeName, err := publiccloud.NewTypeNameFromValue(
		i.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	opts := publiccloud.NewLaunchInstanceOpts(
		*regionName,
		*typeName,
		image.ID.ValueString(),
		*contractType,
		*contractTerm,
		*billingFrequency,
		*rootDiskStorageType,
	)

	opts.MarketAppId = utils.AdaptStringPointerValueToNullableString(i.MarketAppID)
	opts.Reference = utils.AdaptStringPointerValueToNullableString(i.Reference)
	opts.RootDiskSize = utils.AdaptInt32PointerValueToNullableInt32(i.RootDiskSize)

	return opts, nil
}

func (i instanceResourceModel) getUpdateOpts(ctx context.Context) (
	*publiccloud.UpdateInstanceOpts,
	error,
) {
	opts := publiccloud.NewUpdateInstanceOpts()

	opts.Reference = utils.AdaptStringPointerValueToNullableString(i.Reference)
	opts.RootDiskSize = utils.AdaptInt32PointerValueToNullableInt32(i.RootDiskSize)

	contract := contractResourceModel{}
	diags := i.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if diags.HasError() {
		return nil, utils.ReturnError("getUpdateOpts", diags)
	}

	if contract.Type.ValueString() != "" {
		contractType, err := publiccloud.NewContractTypeFromValue(
			contract.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf("getUpdateOpts: %w", err)
		}
		opts.ContractType = contractType
	}

	if contract.Term.ValueInt32() != 0 {
		contractTerm, err := publiccloud.NewContractTermFromValue(
			contract.Term.ValueInt32(),
		)
		if err != nil {
			return nil, fmt.Errorf("getUpdateOpts: %w", err)
		}
		opts.ContractTerm = contractTerm
	}

	if contract.BillingFrequency.ValueInt32() != 0 {
		billingFrequency, err := publiccloud.NewBillingFrequencyFromValue(
			contract.BillingFrequency.ValueInt32(),
		)
		if err != nil {
			return nil, fmt.Errorf("getUpdateOpts: %w", err)
		}
		opts.BillingFrequency = billingFrequency
	}

	if i.Type.ValueString() != "" {
		instanceType, err := publiccloud.NewTypeNameFromValue(
			i.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf("getUpdateOpts: %w", err)
		}
		opts.Type = instanceType
	}

	return opts, nil
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
		imageResourceModel{}.attributeTypes(),
		ctx,
		adaptImageToImageResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.Image = image

	ips, err := utils.AdaptSdkModelsToListValue(
		instanceDetails.Ips,
		ipResourceModel{}.attributeTypes(),
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
		isoResourceModel{}.attributeTypes(),
		ctx,
		adaptIsoToISOResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptInstanceToInstanceResource: %w", err)
	}
	instance.ISO = iso

	return &instance, nil
}

func NewInstanceResource() resource.Resource {
	return &instanceResource{
		name: "public_cloud_instance",
	}
}

type instanceResource struct {
	name   string
	client publiccloud.PubliccloudAPI
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

	i.client = coreClient.PubliccloudAPI
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

	summary := fmt.Sprintf("Launching resource %s", i.name)

	opts, err := plan.getLaunchOpts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, utils.DefaultErrMsg)
		return
	}

	instance, httpResponse, err := i.client.LaunchInstance(ctx).
		LaunchInstanceOpts(*opts).
		Execute()
	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}

	// Get ISO data from instanceDetails
	instanceDetails, httpResponse, err := i.client.GetInstance(
		ctx,
		instance.GetId(),
	).Execute()
	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}

	state, err := adaptInstanceDetailsToInstanceResource(*instanceDetails, ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, utils.DefaultErrMsg)
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

	httpResponse, err := i.client.TerminateInstance(
		ctx,
		state.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.Error(
			ctx,
			&resp.Diagnostics,
			fmt.Sprintf(
				"Terminating resource %s for id %q",
				i.name,
				state.ID.ValueString(),
			),
			err,
			httpResponse,
		)
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
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	summary := fmt.Sprintf(
		"Reading resource %s for id %q",
		i.name,
		state.ID.ValueString(),
	)

	instanceDetails, httpResponse, err := i.client.
		GetInstance(ctx, state.ID.ValueString()).
		Execute()
	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}

	newState, err := adaptInstanceDetailsToInstanceResource(
		*instanceDetails,
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError(summary, utils.DefaultErrMsg)
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

	summary := fmt.Sprintf(
		"Updating resource %s for id %q",
		i.name,
		plan.ID.ValueString(),
	)

	opts, err := plan.getUpdateOpts(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, utils.DefaultErrMsg)
		return
	}

	instanceDetails, httpResponse, err := i.client.
		UpdateInstance(ctx, plan.ID.ValueString()).
		UpdateInstanceOpts(*opts).
		Execute()
	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}

	state, err := adaptInstanceDetailsToInstanceResource(*instanceDetails, ctx)
	if err != nil {
		utils.Error(ctx, &resp.Diagnostics, summary, err, nil)
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
						Description: "Can be either an Operating System or a UUID in case of a Custom Image ID." + warningError,
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
				Description: "Market App ID that must be installed into the instance." + warningError,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}
