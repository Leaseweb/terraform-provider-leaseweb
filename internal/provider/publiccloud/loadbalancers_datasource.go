package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &LoadBalancersDataSource{}
)

type datasourceModelLoadBalancer struct {
	ID        types.String            `tfsdk:"id"`
	IPs       []iPDataSourceModel     `tfsdk:"ips"`
	Reference types.String            `tfsdk:"reference"`
	Contract  contractDataSourceModel `tfsdk:"contract"`
	State     types.String            `tfsdk:"state"`
	Region    types.String            `tfsdk:"region"`
	Type      types.String            `tfsdk:"type"`
}

type datasourceModelLoadBalancers struct {
	LoadBalancers []datasourceModelLoadBalancer `tfsdk:"loadbalancers"`
}

func adaptSdkLoadBalancerDetailsToDatasourceLoadBalancer(sdkLoadBalancerDetails publicCloud.LoadBalancerDetails) datasourceModelLoadBalancer {
	var ips []iPDataSourceModel
	for _, ip := range sdkLoadBalancerDetails.Ips {
		ips = append(ips, iPDataSourceModel{IP: basetypes.NewStringValue(ip.GetIp())})
	}

	return datasourceModelLoadBalancer{
		ID:        basetypes.NewStringValue(sdkLoadBalancerDetails.GetId()),
		IPs:       ips,
		Reference: basetypes.NewStringPointerValue(sdkLoadBalancerDetails.Reference.Get()),
		Contract:  adaptContractToContractDataSource(sdkLoadBalancerDetails.GetContract()),
		State:     basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetState())),
		Region:    basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetRegion())),
		Type:      basetypes.NewStringValue(string(sdkLoadBalancerDetails.GetType())),
	}
}

func adaptSdkLoadBalancersToDatasourceLoadBalancers(sdkLoadBalancers []publicCloud.LoadBalancerDetails) datasourceModelLoadBalancers {
	var loadBalancers datasourceModelLoadBalancers

	for _, sdkLoadBalancerDetails := range sdkLoadBalancers {
		loadBalancer := adaptSdkLoadBalancerDetailsToDatasourceLoadBalancer(sdkLoadBalancerDetails)
		loadBalancers.LoadBalancers = append(loadBalancers.LoadBalancers, loadBalancer)
	}

	return loadBalancers
}

func getAllLoadBalancers(
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) ([]publicCloud.LoadBalancerDetails, *utils.SdkError) {
	var loadBalancers []publicCloud.LoadBalancerDetails

	request := api.GetLoadBalancerList(ctx)

	result, response, err := request.Execute()

	if err != nil {
		return nil, utils.NewSdkError("getAllLoadBalancers", err, response)
	}

	metadata := result.GetMetadata()
	pagination := utils.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, utils.NewSdkError("getAllImages", err, response)
		}

		loadBalancers = append(loadBalancers, result.GetLoadBalancers()...)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, utils.NewSdkError("getAllImages", err, response)
		}
	}

	return loadBalancers, nil
}

type LoadBalancersDataSource struct {
	client client.Client
}

func (l *LoadBalancersDataSource) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = request.ProviderTypeName + "_public_cloud_loadbalancers"
}

func (l *LoadBalancersDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"loadbalancers": schema.ListNestedAttribute{
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
							Validators: []validator.Object{contractTermValidator{}},
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

func (l *LoadBalancersDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	tflog.Info(ctx, "Read publiccloud loadBalancers")
	loadBalancers, err := getAllLoadBalancers(ctx, l.client.PublicCloudAPI)

	if err != nil {
		response.Diagnostics.AddError("Unable to read loadBalancers", err.Error())
		utils.LogError(
			ctx,
			err.ErrorResponse,
			&response.Diagnostics,
			"Unable to read loadBalancers",
			err.Error(),
		)

		return
	}

	state := adaptSdkLoadBalancersToDatasourceLoadBalancers(loadBalancers)

	diags := response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (l *LoadBalancersDataSource) Configure(
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
	return &LoadBalancersDataSource{}
}
