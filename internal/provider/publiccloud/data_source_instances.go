package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/dataadapters/shared"
	datasourceModel "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/datasource"
	customValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/validator"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/logging"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/service"
)

var (
	_ datasource.DataSource              = &instancesDataSource{}
	_ datasource.DataSourceWithConfigure = &instancesDataSource{}
)

func NewInstancesDataSource() datasource.DataSource {
	return &instancesDataSource{}
}

type instancesDataSource struct {
	client client.Client
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

	d.client = coreClient
}

func (d *instancesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_public_cloud_instances"
}

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {

	tflog.Info(ctx, "Read public cloud instances")
	instances, err := d.client.PublicCloudRepository.GetAllInstances(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Unable to read instances", err.Error())
		logging.ServiceError(
			ctx,
			err.ErrorResponse,
			&resp.Diagnostics,
			"Unable to read instances",
			err.Error(),
		)

		return
	}

	state := datasourceModel.NewInstances(instances)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *instancesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := service.NewIntMarkdownList(
		append(
			[]publicCloud.BillingFrequency{0},
			publicCloud.AllowedBillingFrequencyEnumValues...,
		),
	)
	contractTerms := service.NewIntMarkdownList(publicCloud.AllowedContractTermEnumValues)

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
										stringvalidator.OneOf(shared.AdaptContractTypesToStringArray(publicCloud.AllowedContractTypeEnumValues)...),
									},
								},
								"ends_at": schema.StringAttribute{Computed: true},
								"state": schema.StringAttribute{
									Computed: true,
								},
							},
							Validators: []validator.Object{customValidator.ContractTermIsValid()},
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
