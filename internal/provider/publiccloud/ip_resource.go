package publiccloud

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
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

type ipResource struct {
	utils.ResourceAPI
}

func adaptIpDetailsToIPResource(ipDetails publiccloud.IpDetails) ipResourceModel {
	reverseLookup, _ := ipDetails.GetReverseLookupOk()
	return ipResourceModel{
		ReverseLookup: basetypes.NewStringPointerValue(reverseLookup),
		IP:            basetypes.NewStringValue(ipDetails.GetIp()),
	}
}

func (i *ipResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		utils.UnexpectedImportIdentifierError(
			&response.Diagnostics,
			"instance_id,ip",
			request.ID,
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
	_ context.Context,
	_ resource.CreateRequest,
	response *resource.CreateResponse,
) {
	utils.ImportOnlyError(&response.Diagnostics)
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

	ip, httpResponse, err := i.PubliccloudAPI.GetInstanceIP(
		ctx,
		state.InstanceID.ValueString(),
		state.IP.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
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

	ipDetails, httpResponse, err := i.PubliccloudAPI.UpdateInstanceIP(
		ctx,
		plan.InstanceID.ValueString(),
		plan.IP.ValueString(),
	).UpdateIPOpts(*publiccloud.NewUpdateIPOpts(plan.ReverseLookup.ValueString())).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptIpDetailsToIPResource(*ipDetails)
	state.InstanceID = plan.InstanceID

	response.Diagnostics.Append(
		response.State.Set(ctx, state)...,
	)
}

func (i *ipResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}

func NewIPResource() resource.Resource {
	return &ipResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "public_cloud_ip",
		},
	}
}
