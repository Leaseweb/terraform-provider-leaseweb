// Package public_cloud implements schemas used multiple times in public_cloud data sources & resources.
package public_cloud

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	facade "github.com/leaseweb/terraform-provider-leaseweb/internal/facades/public_cloud"
	customValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/validator"
)

func Resources() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"cpu": schema.SingleNestedAttribute{
				Description: "Number of cores",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{Computed: true},
					"unit":  schema.StringAttribute{Computed: true},
				}},
			"memory": schema.SingleNestedAttribute{
				Description: "Total memory in GiB",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Float64Attribute{Computed: true},
					"unit":  schema.StringAttribute{Computed: true},
				}},
			"public_network_speed": schema.SingleNestedAttribute{
				Description: "Public network speed in Gbps",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{Computed: true},
					"unit":  schema.StringAttribute{Computed: true},
				}},
			"private_network_speed": schema.SingleNestedAttribute{
				Description: "Private network speed in Gbps",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.Int64Attribute{Computed: true},
					"unit":  schema.StringAttribute{Computed: true},
				}},
		},
		Description: "Available resources",
		Computed:    true,
	}
}

func Contract(
	required bool,
	facade facade.PublicCloudFacade,
) schema.SingleNestedAttribute {
	if required {
		return schema.SingleNestedAttribute{
			Computed: !required,
			Required: required,
			Attributes: map[string]schema.Attribute{
				"billing_frequency": schema.Int64Attribute{
					Computed:    !required,
					Required:    required,
					Description: "The billing frequency (in months). Valid options are " + facade.GetBillingFrequencies().Markdown(),
					Validators: []validator.Int64{
						int64validator.OneOf(facade.GetBillingFrequencies().ToInt64()...),
					},
				},
				"term": schema.Int64Attribute{
					Computed:    !required,
					Required:    required,
					Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + facade.GetContractTerms().Markdown(),
					Validators: []validator.Int64{
						int64validator.OneOf(facade.GetContractTerms().ToInt64()...),
					},
				},
				"type": schema.StringAttribute{
					Computed:    !required,
					Required:    required,
					Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
					Validators: []validator.String{
						stringvalidator.OneOf(facade.GetContractTypes()...),
					},
				},
				"state": schema.StringAttribute{
					Computed: true,
				},
				"ends_at": schema.StringAttribute{Computed: true},
				"renewals_at": schema.StringAttribute{
					Computed:    true,
					Description: "Date when the contract will be automatically renewed",
				},
				"created_at": schema.StringAttribute{
					Computed:    true,
					Description: "Date when the contract was created",
				},
			},
			Validators: []validator.Object{customValidator.ContractTermIsValid()},
		}
	}

	return schema.SingleNestedAttribute{
		Computed: !required,
		Required: required,
		Attributes: map[string]schema.Attribute{
			"billing_frequency": schema.Int64Attribute{
				Computed:    !required,
				Required:    required,
				Description: "The billing frequency (in months). Valid options are " + facade.GetBillingFrequencies().Markdown(),
				Validators: []validator.Int64{
					int64validator.OneOf(facade.GetBillingFrequencies().ToInt64()...),
				},
			},
			"term": schema.Int64Attribute{
				Computed:    !required,
				Required:    required,
				Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + facade.GetContractTerms().Markdown(),
				Validators: []validator.Int64{
					int64validator.OneOf(facade.GetContractTerms().ToInt64()...),
				},
			},
			"type": schema.StringAttribute{
				Computed:    !required,
				Required:    required,
				Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
				Validators: []validator.String{
					stringvalidator.OneOf(facade.GetContractTypes()...),
				},
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
		},
		Validators: []validator.Object{customValidator.ContractTermIsValid()},
	}
}

func Network() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"private_network_id": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"subnet": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
