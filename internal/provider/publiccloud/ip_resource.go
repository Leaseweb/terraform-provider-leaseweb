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
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
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

func (i ipResourceModel) generateUpdateOpts() publiccloud.UpdateIPOpts {
	return *publiccloud.NewUpdateIPOpts(i.ReverseLookup.ValueString())
}

type ipResource struct {
	name   string
	client publiccloud.PubliccloudAPI
}

func adaptIpDetailsToIPResource(ipDetails publiccloud.IpDetails) ipResourceModel {
	reverseLookup, _ := ipDetails.GetReverseLookupOk()
	return ipResourceModel{
		ReverseLookup: basetypes.NewStringPointerValue(reverseLookup),
		IP:            basetypes.NewStringValue(ipDetails.GetIp()),
	}
}

func adaptIpToIPResource(ip publiccloud.Ip) ipResourceModel {
	reverseLookup, _ := ip.GetReverseLookupOk()
	return ipResourceModel{
		ReverseLookup: basetypes.NewStringPointerValue(reverseLookup),
		IP:            basetypes.NewStringValue(ip.GetIp()),
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

	ip, httpResponse, err := i.client.GetInstanceIP(
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

	opts := plan.generateUpdateOpts()

	ipDetails, httpResponse, err := i.client.UpdateInstanceIP(
		ctx,
		plan.InstanceID.ValueString(),
		plan.IP.ValueString(),
	).UpdateIPOpts(opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
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
		utils.ConfigError(&response.Diagnostics, request.ProviderData)
		return
	}

	i.client = coreClient.PubliccloudAPI
}

func NewIPResource() resource.Resource {
	return &ipResource{
		name: "public_cloud_ip",
	}
}
