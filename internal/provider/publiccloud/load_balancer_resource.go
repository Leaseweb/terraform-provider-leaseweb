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
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &loadBalancerResource{}
	_ resource.ResourceWithImportState = &loadBalancerResource{}
)

type loadBalancerResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Region    types.String `tfsdk:"region"`
	Type      types.String `tfsdk:"type"`
	Reference types.String `tfsdk:"reference"`
	Contract  types.Object `tfsdk:"contract"`
}

func (l *loadBalancerResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"region":    types.StringType,
		"type":      types.StringType,
		"reference": types.StringType,
		"contract": types.ObjectType{
			AttrTypes: contractResourceModel{}.AttributeTypes(),
		},
	}
}

func (l *loadBalancerResourceModel) GetLaunchLoadBalancerOpts(ctx context.Context) (
	*publicCloud.LaunchLoadBalancerOpts,
	error,
) {
	contract := contractResourceModel{}
	contractDiags := l.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		return nil, utils.ReturnError("GetLaunchLoadBalancerOpts", contractDiags)
	}

	sdkContractType, err := publicCloud.NewContractTypeFromValue(contract.Type.ValueString())
	if err != nil {
		return nil, err
	}

	sdkContractTerm, err := publicCloud.NewContractTermFromValue(contract.Term.ValueInt32())
	if err != nil {
		return nil, err
	}

	sdkBillingFrequency, err := publicCloud.NewBillingFrequencyFromValue(contract.BillingFrequency.ValueInt32())
	if err != nil {
		return nil, err
	}

	sdkRegionName, err := publicCloud.NewRegionNameFromValue(l.Region.ValueString())
	if err != nil {
		return nil, err
	}

	sdkTypeName, err := publicCloud.NewTypeNameFromValue(l.Type.ValueString())
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

func (l *loadBalancerResourceModel) GetUpdateLoadBalancerOpts() (
	*publicCloud.UpdateLoadBalancerOpts,
	error,
) {
	opts := publicCloud.NewUpdateLoadBalancerOpts()
	opts.Reference = utils.AdaptStringPointerValueToNullableString(l.Reference)

	if l.Type.ValueString() != "" {
		instanceType, err := publicCloud.NewTypeNameFromValue(l.Type.ValueString())
		if err != nil {
			return nil, fmt.Errorf("GetUpdateLoadBalancerOpts: %w", err)
		}
		opts.Type = instanceType
	}

	return opts, nil
}

func adaptLoadBalancerDetailsToLoadBalancerResource(
	sdkLoadBalancerDetails publicCloud.LoadBalancerDetails,
	ctx context.Context,
) (*loadBalancerResourceModel, error) {
	loadBalancer := loadBalancerResourceModel{
		ID:        basetypes.NewStringValue(sdkLoadBalancerDetails.GetId()),
		Region:    basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetRegion())),
		Type:      basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetType())),
		Reference: basetypes.NewStringPointerValue(sdkLoadBalancerDetails.Reference.Get()),
	}

	contract, err := utils.AdaptSdkModelToResourceObject(
		sdkLoadBalancerDetails.Contract,
		contractResourceModel{}.AttributeTypes(),
		ctx,
		adaptContractToContractResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	loadBalancer.Contract = contract

	return &loadBalancer, nil
}

type loadBalancerResource struct {
	name   string
	client publicCloud.PublicCloudAPI
}

func (l *loadBalancerResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (l *loadBalancerResource) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
	coreClient, ok := utils.GetResourceClient(request, response)
	if !ok {
		return
	}

	l.client = coreClient.PublicCloudAPI
}

func (l *loadBalancerResource) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, l.name)
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
		},
	}
}

func (l *loadBalancerResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan loadBalancerResourceModel

	diags := request.Plan.Get(ctx, &plan)
	summary := fmt.Sprintf("Launching resource %s", l.name)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Launch Public Cloud load balancer")

	opts, err := plan.GetLaunchLoadBalancerOpts(ctx)
	if err != nil {
		response.Diagnostics.AddError(summary, utils.DefaultErrMsg)
		return
	}

	sdkLoadBalancer, httpResponse, err := l.client.LaunchLoadBalancer(ctx).
		LaunchLoadBalancerOpts(*opts).
		Execute()

	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	loadBalancer, resourceErr := adaptLoadBalancerDetailsToLoadBalancerResource(
		*sdkLoadBalancer,
		ctx,
	)
	if resourceErr != nil {
		response.Diagnostics.AddError(summary, utils.DefaultErrMsg)
		return
	}

	diags = response.State.Set(ctx, loadBalancer)
	response.Diagnostics.Append(diags...)
}

func (l *loadBalancerResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var state loadBalancerResourceModel
	diags := request.State.Get(ctx, &state)
	summary := fmt.Sprintf(
		"Reading resource %s for id %q",
		l.name,
		state.ID.ValueString(),
	)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	sdkLoadBalancerDetails, httpResponse, err := l.client.
		GetLoadBalancer(ctx, state.ID.ValueString()).
		Execute()
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	instance, resourceErr := adaptLoadBalancerDetailsToLoadBalancerResource(*sdkLoadBalancerDetails, ctx)
	if resourceErr != nil {
		utils.Error(ctx, &response.Diagnostics, summary, resourceErr, nil)
		return
	}

	diags = response.State.Set(ctx, instance)
	response.Diagnostics.Append(diags...)
}

func (l *loadBalancerResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan loadBalancerResourceModel

	diags := request.Plan.Get(ctx, &plan)
	summary := fmt.Sprintf(
		"Updating resource %s for id %q",
		l.name,
		plan.ID.ValueString(),
	)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	opts, err := plan.GetUpdateLoadBalancerOpts()
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, nil)
		return
	}

	sdkLoadBalancer, httpResponse, err := l.client.
		UpdateLoadBalancer(ctx, plan.ID.ValueString()).
		UpdateLoadBalancerOpts(*opts).
		Execute()
	if err != nil {
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	diags = response.State.Set(ctx, sdkLoadBalancer)
	response.Diagnostics.Append(diags...)
}

func (l *loadBalancerResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state loadBalancerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := l.client.TerminateLoadBalancer(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		summary := fmt.Sprintf("Terminating resource %s for id %q", l.name, state.ID.ValueString())
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}
}

func NewLoadBalancerResource() resource.Resource {
	return &loadBalancerResource{
		name: "public_cloud_load_balancer",
	}
}
