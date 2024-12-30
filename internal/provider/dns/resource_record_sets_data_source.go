package dns

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &resourceRecordSet{}
)

type resourceRecordSetsDataSourceModel struct {
	DomainName         types.String                       `tfsdk:"domain_name"`
	InfoMessage        types.String                       `tfsdk:"info_message"`
	ResourceRecordSets []resourceRecordSetDataSourceModel `tfsdk:"resource_record_sets"`
}

type resourceRecordSetDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	RecordType types.String `tfsdk:"type"`
	Content    []string     `tfsdk:"content"`
	TTL        types.Int32  `tfsdk:"ttl"`
}

type resourceRecordSet struct {
	utils.DataSourceAPI
}

func (r *resourceRecordSet) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: "List resource record sets",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "Domain Name",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"info_message": schema.StringAttribute{
				Computed:    true,
				Description: "Optional additional information",
			},
			"resource_record_sets": schema.ListNestedAttribute{
				Description: "Array of resource record sets",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the resource record set",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the resource record set",
						},
						"content": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Array of resource record set Content entries",
						},
						"ttl": schema.Int32Attribute{
							Computed:    true,
							Description: "Time to live of the resource record set",
						},
					},
				},
			},
		},
	}
}

func (r *resourceRecordSet) Read(
	ctx context.Context,
	request datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	var config resourceRecordSetsDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)

	result, httpResponse, err := r.DNSAPI.GetResourceRecordSetList(
		ctx,
		config.DomainName.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	var resourceRecordSets []resourceRecordSetDataSourceModel
	for _, resourceRecordSetDetails := range result.GetResourceRecordSets() {
		resourceRecordSets = append(
			resourceRecordSets,
			resourceRecordSetDataSourceModel{
				Name:       basetypes.NewStringValue(resourceRecordSetDetails.GetName()),
				RecordType: basetypes.NewStringValue(string(resourceRecordSetDetails.GetType())),
				Content:    resourceRecordSetDetails.GetContent(),
				TTL:        basetypes.NewInt32Value(int32(resourceRecordSetDetails.GetTtl())),
			},
		)
	}

	response.Diagnostics.Append(
		response.State.Set(
			ctx,
			resourceRecordSetsDataSourceModel{
				DomainName:         config.DomainName,
				InfoMessage:        basetypes.NewStringValue(result.GetInfoMessage()),
				ResourceRecordSets: resourceRecordSets,
			},
		)...,
	)
}

func NewResourceRecordSetsDataSource() datasource.DataSource {
	return &resourceRecordSet{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "dns_resource_record_sets",
		},
	}
}
