package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/model"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/resource"
)

var (
	_ datasource.DataSource              = &InstancesDataSource{}
	_ datasource.DataSourceWithConfigure = &InstancesDataSource{}
)

type dataSourceModelContract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func newDataSourceModelContract(sdkContract publicCloud.Contract) dataSourceModelContract {
	return dataSourceModelContract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.BillingFrequency)),
		Term:             basetypes.NewInt64Value(int64(sdkContract.Term)),
		Type:             basetypes.NewStringValue(string(sdkContract.Type)),
		EndsAt:           model.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.State)),
	}
}

type dataSourceModelInstance struct {
	Id                  types.String            `tfsdk:"id"`
	Region              types.String            `tfsdk:"region"`
	Reference           types.String            `tfsdk:"reference"`
	Image               dataSourceModelImage    `tfsdk:"image"`
	State               types.String            `tfsdk:"state"`
	Type                types.String            `tfsdk:"type"`
	RootDiskSize        types.Int64             `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String            `tfsdk:"root_disk_storage_type"`
	Ips                 []dataSourceModelIp     `tfsdk:"ips"`
	Contract            dataSourceModelContract `tfsdk:"contract"`
	MarketAppId         types.String            `tfsdk:"market_app_id"`
}

func newDataSourceModelInstance(sdkInstance publicCloud.Instance) dataSourceModelInstance {
	var ips []dataSourceModelIp
	for _, ip := range sdkInstance.Ips {
		ips = append(ips, newDataSourceModelIp(ip))
	}

	return dataSourceModelInstance{
		Id:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           model.AdaptNullableStringToStringValue(sdkInstance.Reference.Get()),
		Image:               newDataSourceModelImage(sdkInstance.Image),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		Ips:                 ips,
		Contract:            newDataSourceModelContract(sdkInstance.Contract),
		MarketAppId:         model.AdaptNullableStringToStringValue(sdkInstance.MarketAppId.Get()),
	}
}

type dataSourceModelImage struct {
	Id types.String `tfsdk:"id"`
}

func newDataSourceModelImage(sdkImage publicCloud.Image) dataSourceModelImage {
	return dataSourceModelImage{
		Id: basetypes.NewStringValue(sdkImage.Id),
	}
}

type dataSourceModelIp struct {
	Ip types.String `tfsdk:"ip"`
}

func newDataSourceModelIp(sdkIp publicCloud.Ip) dataSourceModelIp {
	return dataSourceModelIp{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}
}

type dataSourceModelInstances struct {
	Instances []dataSourceModelInstance `tfsdk:"instances"`
}

func newDataSourceModelInstances(sdkInstances []publicCloud.Instance) dataSourceModelInstances {
	var instances dataSourceModelInstances

	for _, sdkInstance := range sdkInstances {
		instance := newDataSourceModelInstance(sdkInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

func NewInstancesDataSource() datasource.DataSource {
	return &InstancesDataSource{}
}

type InstancesDataSource struct {
	client client.Client
}

func (d *InstancesDataSource) Configure(
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

	d.client = coreClient
}

func (d *InstancesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_instances"
}

func (d *InstancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	tflog.Info(ctx, "Read public cloud instances")
	instances, err := d.client.PublicCloudRepository.GetAllInstances(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read instances", err.Error())
		logging.LogError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Unable to read instances",
			err.Error(),
		)

		return
	}

	state := newDataSourceModelInstances(instances)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *InstancesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := resource.NewIntMarkdownList(
		append(
			[]publicCloud.BillingFrequency{0},
			publicCloud.AllowedBillingFrequencyEnumValues...,
		),
	)
	contractTerms := resource.NewIntMarkdownList(publicCloud.AllowedContractTermEnumValues)

	resp.Schema = schema.Schema{
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
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "Image ID",
								},
							},
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The instance's current state",
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"root_disk_size": schema.Int64Attribute{
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
								"billing_frequency": schema.Int64Attribute{
									Computed:    true,
									Description: "The billing frequency (in months). Valid options are " + billingFrequencies.Markdown(),
									Validators: []validator.Int64{
										int64validator.OneOf(billingFrequencies.ToInt64()...),
									},
								},
								"term": schema.Int64Attribute{
									Computed:    true,
									Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + contractTerms.Markdown(),
									Validators: []validator.Int64{
										int64validator.OneOf(contractTerms.ToInt64()...),
									},
								},
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
									Validators: []validator.String{
										stringvalidator.OneOf(model.AdaptStringTypeArrayToStringArray(publicCloud.AllowedContractTypeEnumValues)...),
									},
								},
								"ends_at": schema.StringAttribute{Computed: true},
								"state": schema.StringAttribute{
									Computed: true,
								},
							},
							Validators: []validator.Object{contractTermValidator{}},
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
