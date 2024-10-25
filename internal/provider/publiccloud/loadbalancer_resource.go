package publiccloud

import (
	"context"
	"fmt"

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
	_ resource.ResourceWithConfigure   = &loadBalancerResource{}
	_ resource.ResourceWithImportState = &loadBalancerResource{}
	_ resource.ResourceWithModifyPlan  = &loadBalancerResource{}
)

type resourceModelLoadBalancer struct {
	ID        types.String `tfsdk:"id"`
	Region    types.String `tfsdk:"region"`
	Type      types.String `tfsdk:"type"`
	Reference types.String `tfsdk:"reference"`
	Contract  types.Object `tfsdk:"contract"`
}

func (l *resourceModelLoadBalancer) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"region":    types.StringType,
		"type":      types.StringType,
		"reference": types.StringType,
		"contract": types.ObjectType{
			AttrTypes: resourceModelContract{}.AttributeTypes(),
		},
	}
}

func (l *resourceModelLoadBalancer) GetLaunchLoadBalancerOpts(ctx context.Context) (
	*publicCloud.LaunchLoadBalancerOpts,
	error,
) {
	contract := resourceModelContract{}
	contractDiags := l.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		return nil, utils.ReturnError("GetLaunchLoadBalancerOpts", contractDiags)
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
		l.Region.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	sdkTypeName, err := publicCloud.NewTypeNameFromValue(
		l.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	opts := publicCloud.NewLaunchLoadBalancerOpts(
		*sdkRegionName,
		*sdkTypeName,
		*sdkContractType,
		*sdkContractTerm,
		*sdkBillingFrequency,
	)

	opts.Reference = utils.AdaptStringPointerValueToNullableString(l.Reference)

	return opts, nil
}

func (l *resourceModelLoadBalancer) GetUpdateLoadBalancerOpts() (
	*publicCloud.UpdateLoadBalancerOpts,
	error,
) {
	opts := publicCloud.NewUpdateLoadBalancerOpts()
	opts.Reference = utils.AdaptStringPointerValueToNullableString(l.Reference)

	if l.Type.ValueString() != "" {
		instanceType, err := publicCloud.NewTypeNameFromValue(
			l.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateLoadBalancerOpts: %w", err)
		}
		opts.Type = instanceType
	}

	return opts, nil
}

func adaptSdkLoadBalancerDetailsToResourceLoadBalancer(
	sdkLoadBalancerDetails publicCloud.LoadBalancerDetails,
	ctx context.Context,
) (*resourceModelLoadBalancer, error) {
	loadBalancer := resourceModelLoadBalancer{
		ID:        basetypes.NewStringValue(sdkLoadBalancerDetails.GetId()),
		Region:    basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetRegion())),
		Type:      basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetType())),
		Reference: basetypes.NewStringPointerValue(sdkLoadBalancerDetails.Reference.Get()),
	}

	contract, err := utils.AdaptSdkModelToResourceObject(
		sdkLoadBalancerDetails.Contract,
		resourceModelContract{}.AttributeTypes(),
		ctx,
		adaptSdkContractToResourceContract,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	loadBalancer.Contract = contract

	return &loadBalancer, nil
}

type loadBalancerResource struct {
	client client.Client
}

func (l *loadBalancerResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		request,
		response,
	)
}

func (l *loadBalancerResource) Configure(
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

	l.client = coreClient
}

func (l *loadBalancerResource) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = request.ProviderTypeName + "_public_cloud_loadbalancer"
}

func (l *loadBalancerResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	warningError := "**WARNING!** Changing this value once running will cause this loadbalancer to be destroyed and a new one to be created."

	contractTerms := utils.NewIntMarkdownList(publicCloud.AllowedContractTermEnumValues)
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := utils.NewIntMarkdownList(
		append(
			[]publicCloud.BillingFrequency{0},
			publicCloud.AllowedBillingFrequencyEnumValues...,
		),
	)

	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The load balancer unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reference": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "An identifying name you can refer to the load balancer",
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
			"region": schema.StringAttribute{
				Required: true,
				Description: fmt.Sprintf(
					"%s Valid options are %s",
					warningError,
					utils.GenerateMarkdownFromEnumsSlice(publicCloud.AllowedRegionNameEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedRegionNameEnumValues)...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Description: fmt.Sprintf(
					"%s Valid options are %s",
					warningError,
					utils.GenerateMarkdownFromEnumsSlice(publicCloud.AllowedTypeNameEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.Expressions{path.MatchRoot("region")}...,
					),
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publicCloud.AllowedTypeNameEnumValues)...),
				},
			},
		},
	}
}

func (l *loadBalancerResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan resourceModelLoadBalancer

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Launch publiccloud loadBalancer")

	opts, err := plan.GetLaunchLoadBalancerOpts(ctx)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating launch loadBalancer opts",
			err.Error(),
		)

		return
	}

	sdkLoadBalancer, apiResponse, err := l.client.PublicCloudAPI.LaunchLoadBalancer(ctx).
		LaunchLoadBalancerOpts(*opts).
		Execute()

	if err != nil {
		sdkErr := utils.NewSdkError("", err, apiResponse)
		response.Diagnostics.AddError(
			"Error launching loadBalancer",
			sdkErr.Error(),
		)

		utils.LogError(
			ctx,
			sdkErr.ErrorResponse,
			&response.Diagnostics,
			"Error launching publiccloud loadBalancer",
			sdkErr.Error(),
		)

		return
	}

	loadBalancer, resourceErr := adaptSdkLoadBalancerDetailsToResourceLoadBalancer(
		*sdkLoadBalancer,
		ctx,
	)
	if resourceErr != nil {
		response.Diagnostics.AddError(
			"Error creating publiccloud loadBalancer resource",
			resourceErr.Error(),
		)

		return
	}

	diags = response.State.Set(ctx, loadBalancer)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (l *loadBalancerResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var state resourceModelLoadBalancer
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(
		ctx,
		fmt.Sprintf(
			"Read publiccloud loadBalancer %q",
			state.ID.ValueString(),
		),
	)
	sdkLoadBalancerDetails, apiResponse, err := l.client.PublicCloudAPI.
		GetLoadBalancer(ctx, state.ID.ValueString()).
		Execute()
	if err != nil {
		sdkError := utils.NewSdkError("", err, apiResponse)
		response.Diagnostics.AddError(
			"Error reading publiccloud loadBalancer",
			sdkError.Error(),
		)

		utils.LogError(
			ctx,
			sdkError.ErrorResponse,
			&response.Diagnostics,
			fmt.Sprintf(
				"Unable to read publiccloud loadBalancer %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)

		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Create publiccloud loadBalancer resource for %q",
		state.ID.ValueString(),
	))
	instance, resourceErr := adaptSdkLoadBalancerDetailsToResourceLoadBalancer(
		*sdkLoadBalancerDetails,
		ctx,
	)
	if resourceErr != nil {
		response.Diagnostics.AddError(
			"Error creating publiccloud loadBalancer resource",
			resourceErr.Error(),
		)

		return
	}

	diags = response.State.Set(ctx, instance)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (l *loadBalancerResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan resourceModelLoadBalancer

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Update publiccloud loadBalancer %q",
		plan.ID.ValueString(),
	))
	opts, err := plan.GetUpdateLoadBalancerOpts()
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating UpdateInstanceOpts",
			err.Error(),
		)
		return
	}

	sdkLoadBalancer, apiResponse, err := l.client.PublicCloudAPI.
		UpdateLoadBalancer(ctx, plan.ID.ValueString()).
		UpdateLoadBalancerOpts(*opts).
		Execute()
	if err != nil {
		sdkErr := utils.NewSdkError("", err, apiResponse)

		response.Diagnostics.AddError(
			"Error updating publiccloud loadBalancer",
			sdkErr.Error(),
		)

		utils.LogError(
			ctx,
			sdkErr.ErrorResponse,
			&response.Diagnostics,
			fmt.Sprintf(
				"Unable to update publiccloud loadBalancer %q",
				plan.ID.ValueString(),
			),
			sdkErr.Error(),
		)

		return
	}

	diags = response.State.Set(ctx, sdkLoadBalancer)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (l *loadBalancerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state resourceModelLoadBalancer
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf(
		"Terminate publiccloud loadBalancer %q",
		state.ID.ValueString(),
	))
	err := terminateInstance(state.ID.ValueString(), ctx, l.client.PublicCloudAPI)

	if err != nil {
		response.Diagnostics.AddError(
			"Error terminating publiccloud loadBalancer",
			fmt.Sprintf(
				"Could not terminate publiccloud loadBalancer, unexpected error: %q",
				err.Error(),
			),
		)

		utils.LogError(
			ctx,
			err.ErrorResponse,
			&response.Diagnostics,
			fmt.Sprintf(
				"Error terminating publiccloud loadBalancer %q",
				state.ID.ValueString(),
			),
			err.Error(),
		)

		return
	}
}

func (l *loadBalancerResource) ModifyPlan(
	_ context.Context,
	_ resource.ModifyPlanRequest,
	_ *resource.ModifyPlanResponse,
) {
}

func NewLoadBalancerResource() resource.Resource {
	return &loadBalancerResource{}
}
