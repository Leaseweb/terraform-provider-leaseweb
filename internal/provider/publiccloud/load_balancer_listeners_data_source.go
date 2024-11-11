package publiccloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &loadBalancerListenersDataSource{}
)

type loadBalancerListenersDataSourceModel struct {
	LoadBalancerID types.String                          `tfsdk:"load_balancer_id"`
	Listeners      []loadBalancerListenerDataSourceModel `tfsdk:"listeners"`
}

func (l loadBalancerListenersDataSourceModel) generateRequest(
	ctx context.Context,
	api publicCloud.PublicCloudAPI,
) publicCloud.ApiGetLoadBalancerListenerListRequest {
	return api.GetLoadBalancerListenerList(ctx, l.LoadBalancerID.ValueString())
}

func adaptLoadBalancerListenersToLoadBalancerListenersDataSource(sdkLoadBalancerListeners []publicCloud.LoadBalancerListener) loadBalancerListenersDataSourceModel {
	var listeners loadBalancerListenersDataSourceModel

	for _, sdkLoadBalancerListener := range sdkLoadBalancerListeners {
		listener := loadBalancerListenerDataSourceModel{
			ID: basetypes.NewStringValue(sdkLoadBalancerListener.GetId()),
		}
		listeners.Listeners = append(listeners.Listeners, listener)
	}

	return listeners
}

type loadBalancerListenerDataSourceModel struct {
	ID types.String `tfsdk:"id"`
}

func getAllLoadBalancerListeners(request publicCloud.ApiGetLoadBalancerListenerListRequest) (
	[]publicCloud.LoadBalancerListener,
	*http.Response,
	error,
) {
	var listeners []publicCloud.LoadBalancerListener
	var offset *int32

	for {
		result, httpResponse, err := request.Execute()
		if err != nil {
			return nil, httpResponse, fmt.Errorf(
				"getAllLoadBalancerListeners: %w",
				err,
			)
		}

		listeners = append(listeners, result.GetListeners()...)

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

	return listeners, nil, nil
}

type loadBalancerListenersDataSource struct {
	name   string
	client publicCloud.PublicCloudAPI
}

func (l *loadBalancerListenersDataSource) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf("%s_%s", request.ProviderTypeName, l.name)
}

func (l *loadBalancerListenersDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer ID",
			},
			"listeners": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The listener unique identifier",
						},
					},
				},
			},
		},
	}
}

func (l *loadBalancerListenersDataSource) Read(
	ctx context.Context,
	request datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	var config loadBalancerListenersDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	summary := fmt.Sprintf(
		"Reading data %s for load_balancer_id %q",
		l.name,
		config.LoadBalancerID,
	)

	tflog.Info(ctx, summary)
	listeners, httpResponse, err := getAllLoadBalancerListeners(config.generateRequest(ctx, l.client))

	if err != nil {
		utils.Error(
			ctx,
			&response.Diagnostics,
			summary,
			err,
			httpResponse,
		)

		return
	}

	state := adaptLoadBalancerListenersToLoadBalancerListenersDataSource(listeners)
	state.LoadBalancerID = basetypes.NewStringValue(config.LoadBalancerID.ValueString())

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (l *loadBalancerListenersDataSource) Configure(
	_ context.Context,
	request datasource.ConfigureRequest,
	response *datasource.ConfigureResponse,
) {
	coreClient, ok := utils.GetDataSourceClient(request, response)
	if !ok {
		return
	}

	l.client = coreClient.PublicCloudAPI
}

func NewLoadBalancerListenersDataSource() datasource.DataSource {
	return &loadBalancerListenersDataSource{
		name: "public_cloud_load_balancer_listeners",
	}
}
