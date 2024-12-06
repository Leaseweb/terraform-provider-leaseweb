package publiccloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &instancesDataSource{}
)

type contractDataSourceModel struct {
	BillingFrequency types.Int32  `tfsdk:"billing_frequency"`
	Term             types.Int32  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func adaptContractToContractDataSource(contract publiccloud.Contract) contractDataSourceModel {
	return contractDataSourceModel{
		BillingFrequency: basetypes.NewInt32Value(int32(contract.GetBillingFrequency())),
		Term:             basetypes.NewInt32Value(int32(contract.GetTerm())),
		Type:             basetypes.NewStringValue(string(contract.GetType())),
		EndsAt:           utils.AdaptNullableTimeToStringValue(contract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(contract.GetState())),
	}
}

type instanceDataSourceModel struct {
	ID                  types.String            `tfsdk:"id"`
	Region              types.String            `tfsdk:"region"`
	Reference           types.String            `tfsdk:"reference"`
	Image               imageModelDataSource    `tfsdk:"image"`
	State               types.String            `tfsdk:"state"`
	Type                types.String            `tfsdk:"type"`
	RootDiskSize        types.Int32             `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String            `tfsdk:"root_disk_storage_type"`
	IPs                 []iPDataSourceModel     `tfsdk:"ips"`
	Contract            contractDataSourceModel `tfsdk:"contract"`
	MarketAppID         types.String            `tfsdk:"market_app_id"`
}

func adaptInstanceToInstanceDataSource(instance publiccloud.Instance) instanceDataSourceModel {
	var ips []iPDataSourceModel
	for _, ip := range instance.Ips {
		ips = append(ips, iPDataSourceModel{IP: basetypes.NewStringValue(ip.GetIp())})
	}

	return instanceDataSourceModel{
		ID:                  basetypes.NewStringValue(instance.GetId()),
		Region:              basetypes.NewStringValue(string(instance.GetRegion())),
		Reference:           basetypes.NewStringPointerValue(instance.Reference.Get()),
		Image:               adaptImageToImageDataSource(instance.GetImage()),
		State:               basetypes.NewStringValue(string(instance.GetState())),
		Type:                basetypes.NewStringValue(string(instance.GetType())),
		RootDiskSize:        basetypes.NewInt32Value(instance.GetRootDiskSize()),
		RootDiskStorageType: basetypes.NewStringValue(string(instance.GetRootDiskStorageType())),
		IPs:                 ips,
		Contract:            adaptContractToContractDataSource(instance.GetContract()),
		MarketAppID:         basetypes.NewStringPointerValue(instance.MarketAppId.Get()),
	}
}

type iPDataSourceModel struct {
	IP types.String `tfsdk:"ip"`
}

type instancesDataSourceModel struct {
	Instances []instanceDataSourceModel `tfsdk:"instances"`
}

func adaptInstancesToInstancesDataSource(sdkInstances []publiccloud.Instance) instancesDataSourceModel {
	var instances instancesDataSourceModel

	for _, sdkInstance := range sdkInstances {
		instance := adaptInstanceToInstanceDataSource(sdkInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func getAllInstances(
	ctx context.Context,
	api publiccloud.PubliccloudAPI,
) ([]publiccloud.Instance, *http.Response, error) {
	var instances []publiccloud.Instance
	var offset *int32

	request := api.GetInstanceList(ctx)

	for {
		result, httpResponse, err := request.Execute()
		if err != nil {
			return nil, httpResponse, fmt.Errorf("getAllInstances: %w", err)
		}
		instances = append(instances, result.Instances...)

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

	return instances, nil, nil
}

func NewInstancesDataSource() datasource.DataSource {
	return &instancesDataSource{
		name: "public_cloud_instances",
	}
}

type instancesDataSource struct {
	name   string
	client publiccloud.PubliccloudAPI
}

func (d *instancesDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	coreClient, ok := req.ProviderData.(client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected provider.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	d.client = coreClient.PubliccloudAPI
}

func (d *instancesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.name)
}

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	instances, httpResponse, err := getAllInstances(ctx, d.client)
	if err != nil {
		summary := fmt.Sprintf("Reading data %s", d.name)
		utils.Error(ctx, &resp.Diagnostics, summary, err, httpResponse)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			adaptInstancesToInstancesDataSource(instances),
		)...,
	)
}

func (d *instancesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := utils.NewIntMarkdownList(
		append(
			[]publiccloud.BillingFrequency{0},
			publiccloud.AllowedBillingFrequencyEnumValues...,
		),
	)
	contractTerms := utils.NewIntMarkdownList(publiccloud.AllowedContractTermEnumValues)

	resp.Schema = schema.Schema{
		Description: utils.BetaDescription,
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The instance unique identifier",
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"reference": schema.StringAttribute{
							Computed:    true,
							Description: "The identifying name set to the instance",
						},
						"image": schema.SingleNestedAttribute{
							Computed:   true,
							Attributes: imageSchemaAttributes(),
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The instance's current state",
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"root_disk_size": schema.Int32Attribute{
							Computed:    true,
							Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances",
						},
						"root_disk_storage_type": schema.StringAttribute{
							Computed:    true,
							Description: "The root disk's storage type",
						},
						"ips": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{Computed: true},
								},
							},
						},
						"contract": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"billing_frequency": schema.Int32Attribute{
									Computed:    true,
									Description: "The billing frequency (in months). Valid options are " + billingFrequencies.Markdown(),
									Validators: []validator.Int32{
										int32validator.OneOf(billingFrequencies.ToInt32()...),
									},
								},
								"term": schema.Int32Attribute{
									Computed:    true,
									Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + contractTerms.Markdown(),
									Validators: []validator.Int32{
										int32validator.OneOf(contractTerms.ToInt32()...),
									},
								},
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
									Validators: []validator.String{
										stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedContractTypeEnumValues)...),
									},
								},
								"ends_at": schema.StringAttribute{Computed: true},
								"state": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"market_app_id": schema.StringAttribute{
							Computed:    true,
							Description: "Market App ID",
						},
					},
				},
			},
		},
	}
}
