package publiccloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &loadBalancersDataSource{}
)

type loadBalancerDataSourceModel struct {
	ID        types.String            `tfsdk:"id"`
	IPs       []ipDataSourceModel     `tfsdk:"ips"`
	Reference types.String            `tfsdk:"reference"`
	Contract  contractDataSourceModel `tfsdk:"contract"`
	State     types.String            `tfsdk:"state"`
	Region    types.String            `tfsdk:"region"`
	Type      types.String            `tfsdk:"type"`
}

type loadBalancersDataSourceModel struct {
	LoadBalancers []loadBalancerDataSourceModel `tfsdk:"load_balancers"`
}

func adaptLoadBalancerListItemToLoadBalancerDataSource(loadBalancerListItem publiccloud.LoadBalancerListItem) loadBalancerDataSourceModel {
	var ips []ipDataSourceModel
	for _, ip := range loadBalancerListItem.Ips {
		ips = append(ips, ipDataSourceModel{IP: basetypes.NewStringValue(ip.GetIp())})
	}

	return loadBalancerDataSourceModel{
		ID:        basetypes.NewStringValue(loadBalancerListItem.GetId()),
		IPs:       ips,
		Reference: basetypes.NewStringPointerValue(loadBalancerListItem.Reference.Get()),
		Contract:  adaptContractToContractDataSource(loadBalancerListItem.GetContract()),
		State:     basetypes.NewStringValue(string(loadBalancerListItem.GetState())),
		Region:    basetypes.NewStringValue(string(loadBalancerListItem.GetRegion())),
		Type:      basetypes.NewStringValue(string(loadBalancerListItem.GetType())),
	}
}

func adaptLoadBalancersToLoadBalancersDataSource(sdkLoadBalancers []publiccloud.LoadBalancerListItem) loadBalancersDataSourceModel {
	var loadBalancers loadBalancersDataSourceModel

	for _, sdkLoadBalancerListItem := range sdkLoadBalancers {
		loadBalancer := adaptLoadBalancerListItemToLoadBalancerDataSource(sdkLoadBalancerListItem)
		loadBalancers.LoadBalancers = append(loadBalancers.LoadBalancers, loadBalancer)
	}

	return loadBalancers
}

func getAllLoadBalancers(
	ctx context.Context,
	api publiccloud.PubliccloudAPI,
) ([]publiccloud.LoadBalancerListItem, *http.Response, error) {
	var loadBalancers []publiccloud.LoadBalancerListItem
	var offset *int32

	request := api.GetLoadBalancerList(ctx)

	for {
		result, httpResponse, err := request.Execute()
		if err != nil {
			return nil, httpResponse, fmt.Errorf("getAllLoadBalancers: %w", err)
		}

		loadBalancers = append(loadBalancers, result.GetLoadBalancers()...)

		metadata := result.GetMetadata()

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		request = request.Offset(*offset)
	}

	return loadBalancers, nil, nil
}

type loadBalancersDataSource struct {
	name   string
	client publiccloud.PubliccloudAPI
}

func (l *loadBalancersDataSource) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, l.name)
}

func (l *loadBalancersDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: utils.BetaDescription,
		Attributes: map[string]schema.Attribute{
			"load_balancers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The load balancer unique identifier",
						},
						"ips": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{Computed: true},
								},
							},
						},
						"reference": schema.StringAttribute{
							Computed: true,
						},
						"contract": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"billing_frequency": schema.Int32Attribute{
									Computed:    true,
									Description: "The billing frequency (in months)",
								},
								"term": schema.Int32Attribute{
									Computed:    true,
									Description: "Contract term (in months)",
								},
								"type": schema.StringAttribute{
									Computed: true,
								},
								"ends_at": schema.StringAttribute{Computed: true},
								"state": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (l *loadBalancersDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	loadBalancers, httpResponse, err := getAllLoadBalancers(ctx, l.client)
	if err != nil {
		summary := fmt.Sprintf("Reading data %s", l.name)
		utils.Error(ctx, &response.Diagnostics, summary, err, httpResponse)
		return
	}

	response.Diagnostics.Append(
		response.State.Set(
			ctx,
			adaptLoadBalancersToLoadBalancersDataSource(loadBalancers),
		)...,
	)
}

func (l *loadBalancersDataSource) Configure(
	_ context.Context,
	request datasource.ConfigureRequest,
	response *datasource.ConfigureResponse,
) {
	if request.ProviderData == nil {
		return
	}

	coreClient, ok := request.ProviderData.(client.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected provider.Client, got: %T. Please report this issue to the provider developers.",
				request.ProviderData,
			),
		)

		return
	}

	l.client = coreClient.PubliccloudAPI
}

func NewLoadBalancersDataSource() datasource.DataSource {
	return &loadBalancersDataSource{
		name: "public_cloud_load_balancers",
	}
}
