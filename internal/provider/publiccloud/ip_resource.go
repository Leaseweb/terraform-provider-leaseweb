package publiccloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.ResourceWithConfigure   = &ipResource{}
	_ resource.ResourceWithImportState = &ipResource{}
)

type ipResourceModel struct {
	ReverseLookup types.String `tfsdk:"reverse_lookup"`
	InstanceID    types.String `tfsdk:"instance_id"`
	IP            types.String `tfsdk:"ip"`
}

func (i ipResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"reverse_lookup": types.StringType,
		"instance_id":    types.StringType,
		"ip":             types.StringType,
	}
}

func (i ipResourceModel) generateUpdateOpts() publiccloud.UpdateIpOpts {
	return *publiccloud.NewUpdateIpOpts(i.ReverseLookup.ValueString())
}

type ipResource struct {
	name   string
	client publiccloud.PubliccloudAPI
}

func adaptIpDetailsToIPResource(sdkIpDetails publiccloud.IpDetails) ipResourceModel {
	reverseLookup, _ := sdkIpDetails.GetReverseLookupOk()
	return ipResourceModel{
		ReverseLookup: basetypes.NewStringPointerValue(reverseLookup),
		IP:            basetypes.NewStringValue(sdkIpDetails.GetIp()),
	}
}

func adaptIpToIPResource(sdkIp publiccloud.Ip) ipResourceModel {
	reverseLookup, _ := sdkIp.GetReverseLookupOk()
	return ipResourceModel{
		ReverseLookup: basetypes.NewStringPointerValue(reverseLookup),
		IP:            basetypes.NewStringValue(sdkIp.GetIp()),
	}
}

func (i *ipResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf(
				"Expected import identifier with format: instance_id,ip. Got: %q",
				request.ID,
			),
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("instance_id"),
		idParts[0],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("ip"),
		idParts[1],
	)...)
}

func (i *ipResource) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, i.name)
}

func (i *ipResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"reverse_lookup": schema.StringAttribute{
				Required: true,
			},
			"instance_id": schema.StringAttribute{
				Required: true,
			},
			"ip": schema.StringAttribute{
				Required: true,
			},
		},
	}

	utils.AddUnsupportedActionsNotation(
		response,
		[]utils.Action{utils.CreateAction, utils.DeleteAction},
	)
}

func (i *ipResource) Create(
	ctx context.Context,
	_ resource.CreateRequest,
	response *resource.CreateResponse,
) {
	utils.Error(
		ctx,
		&response.Diagnostics,
		fmt.Sprintf(
			"Resource %s can only be imported, not created.",
			i.name,
		),
		nil,
		nil,
	)
}

func (i *ipResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var state ipResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	ip, httpResponse, err := i.client.GetIp(
		ctx,
		state.InstanceID.ValueString(),
		state.IP.ValueString(),
	).Execute()
	if err != nil {
		utils.Error(
			ctx,
			&response.Diagnostics,
			fmt.Sprintf(
				"Reading resource %s for instance_id %q ip %q",
				i.name,
				state.InstanceID.ValueString(),
				state.IP.ValueString(),
			),
			err,
			httpResponse,
		)
		return
	}

	newState := adaptIpDetailsToIPResource(*ip)
	newState.InstanceID = state.InstanceID

	response.Diagnostics.Append(response.State.Set(ctx, newState)...)
}

func (i *ipResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan ipResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := plan.generateUpdateOpts()

	ipDetails, httpResponse, err := i.client.UpdateIp(
		ctx,
		plan.InstanceID.ValueString(),
		plan.IP.ValueString(),
	).UpdateIpOpts(opts).Execute()
	if err != nil {
		utils.Error(
			ctx,
			&response.Diagnostics,
			fmt.Sprintf(
				"Updating resource %s for instance_id %q ip %q",
				i.name,
				plan.InstanceID.ValueString(),
				plan.IP.ValueString(),
			),
			err,
			httpResponse,
		)
		return
	}

	newState := adaptIpDetailsToIPResource(*ipDetails)
	newState.InstanceID = plan.InstanceID

	response.Diagnostics.Append(
		response.State.Set(ctx, newState)...,
	)
}

func (i *ipResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}

func (i *ipResource) Configure(
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

func NewIPResource() resource.Resource {
	return &ipResource{
		name: "public_cloud_ip",
	}
}
