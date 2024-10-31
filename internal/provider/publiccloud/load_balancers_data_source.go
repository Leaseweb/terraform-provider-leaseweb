package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &loadBalancersDataSource{}
)

type loadBalancerDataSourceModel struct {
	ID        types.String            `tfsdk:"id"`
	IPs       []iPDataSourceModel     `tfsdk:"ips"`
	Reference types.String            `tfsdk:"reference"`
	Contract  contractDataSourceModel `tfsdk:"contract"`
	State     types.String            `tfsdk:"state"`
	Region    types.String            `tfsdk:"region"`
	Type      types.String            `tfsdk:"type"`
}

type loadBalancersDataSourceModel struct {
	LoadBalancers []loadBalancerDataSourceModel `tfsdk:"load_balancers"`
}

func adaptLoadBalancerDetailsToLoadBalancerDataSource(sdkLoadBalancerDetails publicCloud.LoadBalancerDetails) loadBalancerDataSourceModel {
	var ips []iPDataSourceModel
	for _, ip := range sdkLoadBalancerDetails.Ips {
		ips = append(ips, iPDataSourceModel{IP: basetypes.NewStringValue(ip.GetIp())})
	}

	return loadBalancerDataSourceModel{
		ID:        basetypes.NewStringValue(sdkLoadBalancerDetails.GetId()),
		IPs:       ips,
		Reference: basetypes.NewStringPointerValue(sdkLoadBalancerDetails.Reference.Get()),
		Contract:  adaptContractToContractDataSource(sdkLoadBalancerDetails.GetContract()),
		State:     basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetState())),
		Region:    basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetRegion())),
		Type:      basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetType())),
	}
}

func adaptLoadBalancersToLoadBalancersDataSource(sdkLoadBalancers []publicCloud.LoadBalancerDetails) loadBalancersDataSourceModel {
	var loadBalancers loadBalancersDataSourceModel

	for _, sdkLoadBalancerDetails := range sdkLoadBalancers {
		loadBalancer := adaptLoadBalancerDetailsToLoadBalancerDataSource(sdkLoadBalancerDetails)
		loadBalancers.LoadBalancers = append(loadBalancers.LoadBalancers, loadBalancer)
	}

	return loadBalancers
}

func getAllLoadBalancers(
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) ([]publicCloud.LoadBalancerDetails, *utils.SdkError) {
	var loadBalancers []publicCloud.LoadBalancerDetails
	var offset *int32

	request := api.GetLoadBalancerList(ctx)

	result, response, err := request.Execute()

	if err != nil {
		return nil, utils.NewSdkError("getAllLoadBalancers", err, response)
	}

	metadata := result.GetMetadata()
	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, utils.NewSdkError("getAllLoadBalancers", err, response)
		}

		loadBalancers = append(loadBalancers, result.GetLoadBalancers()...)

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		request.Offset(*offset)
	}

	return loadBalancers, nil
}

type loadBalancersDataSource struct {
	client client.Client
}

func (l *loadBalancersDataSource) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = request.ProviderTypeName + "_public_cloud_load_balancers"
}

func (l *loadBalancersDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
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
								"billing_frequency": schema.Int64Attribute{
									Computed:    true,
									Description: "The billing frequency (in months)",
								},
								"term": schema.Int64Attribute{
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
	tflog.Info(ctx, "Read Public Cloud load balancers")
	loadBalancers, err := getAllLoadBalancers(ctx, l.client.PublicCloudAPI)

	if err != nil {
		response.Diagnostics.AddError("Unable to read Public Cloud load balancers", err.Error())
		utils.LogError(
			ctx,
			err.ErrorResponse,
			&response.Diagnostics,
			"Unable to read Public Cloud load balancers",
			err.Error(),
		)

		return
	}

	state := adaptLoadBalancersToLoadBalancersDataSource(loadBalancers)

	diags := response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
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

	l.client = coreClient
}

func NewLoadBalancersDataSource() datasource.DataSource {
	return &loadBalancersDataSource{}
}
