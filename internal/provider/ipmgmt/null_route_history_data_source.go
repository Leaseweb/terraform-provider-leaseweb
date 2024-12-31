package ipmgmt

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/ipmgmt"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &nullRouteHistoryDataSource{}
)

type nullRouteHistoryDataSourceModel struct {
	ContractID  types.String               `tfsdk:"contract_id"`
	EquipmentID types.String               `tfsdk:"equipment_id"`
	FromDate    types.String               `tfsdk:"from_date"`
	FromIP      types.String               `tfsdk:"from_ip"`
	NulledBy    types.String               `tfsdk:"nulled_by"`
	Sort        []string                   `tfsdk:"sort"`
	TicketID    types.String               `tfsdk:"ticket_id"`
	ToDate      types.String               `tfsdk:"to_date"`
	ToIP        types.String               `tfsdk:"to_ip"`
	UnnulledBy  types.String               `tfsdk:"unnulled_by"`
	NullRoutes  []nullrouteDataSourceModel `tfsdk:"nullroutes"`
}

type nullrouteDataSourceModel struct {
	AssignedContract     *assignedContractDataSourceModel `tfsdk:"assigned_contract"`
	AutomaticUnnullingAt types.String                     `tfsdk:"automatic_unnulling_at"`
	Comment              types.String                     `tfsdk:"comment"`
	EquipmentID          types.String                     `tfsdk:"equipment_id"`
	ID                   types.String                     `tfsdk:"id"`
	IP                   types.String                     `tfsdk:"ip"`
	NulledAt             types.String                     `tfsdk:"nulled_at"`
	NulledBy             types.String                     `tfsdk:"nulled_by"`
	NullLevel            types.Int32                      `tfsdk:"null_level"`
	TicketID             types.String                     `tfsdk:"ticket_id"`
	UnnulledAt           types.String                     `tfsdk:"unnulled_at"`
	UnnulledBy           types.String                     `tfsdk:"unnulled_by"`
}

type nullRouteHistoryDataSource struct {
	utils.DataSourceAPI
}

func (n nullRouteHistoryDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: "Inspect null route history",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by ID of the contract assigned to the IP at the time of null route creation",
			},
			"equipment_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by ID of the contract assigned to the IP at the time of null route creation",
			},
			"from_date": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by ID of the server assigned to the IP at the time of null route creation",
			},
			"from_ip": schema.StringAttribute{
				Optional:    true,
				Description: "Return only IPs greater or equal to the specified address",
			},
			"nulled_by": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by the email address of the user who created the null route",
			},
			"sort": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Sort field names. Prepend the field name with '-' for descending order. E.g. `ip,-nullrouted`. Sortable field names are `ip`, `nullRouted`, `reverseLookup`",
			},
			"ticket_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by the reference stored with the null route",
			},
			"to_date": schema.StringAttribute{
				Optional:    true,
				Description: "Return only null routes active before the specified date and time",
			},
			"to_ip": schema.StringAttribute{
				Optional:    true,
				Description: "Return only IPs lower or equal to the specified address",
			},
			"unnulled_by": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by the email address of the user who removed the null route",
			},

			"nullroutes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"assigned_contract": schema.SingleNestedAttribute{
							Description: "An explanation about the purpose of this instance",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "ID of the contract connected to the IP",
								},
							},
						},
						"automatic_unnulling_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date and time when the null route is to be automatically removed",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "Comment stored with the null route",
						},
						"equipment_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the equipment which was assigned to the IP at the time of null route creation",
						},
						"id": schema.StringAttribute{
							Optional:    true,
							Description: "Null route ID",
						},
						"ip": schema.StringAttribute{
							Computed:    true,
							Description: "IP address",
						},
						"nulled_at": schema.StringAttribute{
							Computed:    true,
							Description: "Null route date",
						},
						"nulled_by": schema.StringAttribute{
							Computed:    true,
							Description: "Email address of the user who created the null route or 'LeaseWeb' if null route was created by LeaseWeb",
						},
						"null_level": schema.Int32Attribute{
							Computed:    true,
							Description: "Null route permission level. If greater than 1 then the null route can only be removed by LeaseWeb",
						},
						"ticket_id": schema.StringAttribute{
							Computed:    true,
							Description: "Reference stored with the null route",
						},
						"unnulled_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date and time when the null route has been removed. If null then the null route is still active",
						},
						"unnulled_by": schema.StringAttribute{
							Computed:    true,
							Description: "Email address of the user who removed the null route or 'LeaseWeb' if null route was removed by LeaseWeb",
						},
					},
				},
			},
		},
	}
}

func (n nullRouteHistoryDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var config nullRouteHistoryDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	var nullRoutes []ipmgmt.NullRoutedIP
	var offset *int32
	var state nullRouteHistoryDataSourceModel

	nullRouteRequest := n.IPmgmtAPI.GetNullRouteHistoryList(ctx)
	if !config.ContractID.IsNull() {
		nullRouteRequest = nullRouteRequest.ContractId(config.ContractID.ValueString())
		state.ContractID = config.ContractID
	}
	if !config.EquipmentID.IsNull() {
		nullRouteRequest = nullRouteRequest.EquipmentId(config.EquipmentID.ValueString())
		state.EquipmentID = config.EquipmentID
	}
	if !config.FromDate.IsNull() {
		nullRouteRequest = nullRouteRequest.FromDate(config.FromDate.ValueString())
		state.FromDate = config.FromDate
	}
	if !config.FromIP.IsNull() {
		nullRouteRequest = nullRouteRequest.FromIp(config.FromIP.ValueString())
		state.FromIP = config.FromIP
	}
	if !config.NulledBy.IsNull() {
		nullRouteRequest = nullRouteRequest.NulledBy(config.NulledBy.ValueString())
		state.NulledBy = config.NulledBy
	}
	if len(config.Sort) > 0 {
		nullRouteRequest = nullRouteRequest.Sort(strings.Join(config.Sort[:], ","))
		state.Sort = config.Sort
	}
	if !config.TicketID.IsNull() {
		nullRouteRequest = nullRouteRequest.TicketId(config.TicketID.ValueString())
		state.TicketID = config.TicketID
	}
	if !config.ToDate.IsNull() {
		nullRouteRequest = nullRouteRequest.ToDate(config.ToDate.ValueString())
		state.ToDate = config.ToDate
	}
	if !config.ToIP.IsNull() {
		nullRouteRequest = nullRouteRequest.ToIp(config.ToIP.ValueString())
		state.ToIP = config.ToIP
	}
	if !config.UnnulledBy.IsNull() {
		nullRouteRequest = nullRouteRequest.UnnulledBy(config.UnnulledBy.ValueString())
		state.UnnulledBy = config.UnnulledBy
	}

	for {
		result, httpResponse, err := nullRouteRequest.Execute()
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
			utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
			return
		}

		nullRoutes = append(nullRoutes, result.GetNullroutes()...)

		metadata := result.GetMetadata()

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		nullRouteRequest = nullRouteRequest.Offset(*offset)
	}

	for _, sdkNullRoute := range nullRoutes {
		var assignedContract *assignedContractDataSourceModel
		sdkAssignedContract, _ := sdkNullRoute.GetAssignedContractOk()
		if sdkAssignedContract != nil {
			assignedContract = &assignedContractDataSourceModel{
				ID: basetypes.NewStringValue(sdkAssignedContract.GetId()),
			}
		}

		automaticUnnullingAt, _ := sdkNullRoute.GetAutomatedUnnullingAtOk()
		comment, _ := sdkNullRoute.GetCommentOk()
		ticketID, _ := sdkNullRoute.GetTicketIdOk()
		unnulledAt, _ := sdkNullRoute.GetUnnulledAtOk()
		unnulledBy, _ := sdkNullRoute.GetUnnulledByOk()

		nullRoute := nullrouteDataSourceModel{
			AssignedContract:     assignedContract,
			AutomaticUnnullingAt: utils.AdaptNullableTimeToStringValue(automaticUnnullingAt),
			Comment:              basetypes.NewStringPointerValue(comment),
			EquipmentID:          basetypes.NewStringValue(sdkNullRoute.GetEquipmentId()),
			ID:                   basetypes.NewStringValue(sdkNullRoute.GetId()),
			IP:                   basetypes.NewStringValue(sdkNullRoute.GetIp()),
			NulledAt:             basetypes.NewStringValue(sdkNullRoute.GetNulledAt().String()),
			NulledBy:             basetypes.NewStringValue(sdkNullRoute.GetNulledBy()),
			NullLevel:            basetypes.NewInt32Value(sdkNullRoute.GetNullLevel()),
			TicketID:             basetypes.NewStringPointerValue(ticketID),
			UnnulledAt:           utils.AdaptNullableTimeToStringValue(unnulledAt),
			UnnulledBy:           basetypes.NewStringPointerValue(unnulledBy),
		}
		state.NullRoutes = append(state.NullRoutes, nullRoute)
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func NewNullRouteHistoryDataSource() datasource.DataSource {
	return &nullRouteHistoryDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "ipmgmt_null_route_history",
		},
	}
}
