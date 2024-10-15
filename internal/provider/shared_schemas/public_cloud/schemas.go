// Package public_cloud implements schemas used multiple times in public_cloud data sources & resources.
package public_cloud

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/ports"
	customValidator "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/instance/validator"
)

func Contract(
	required bool,
	publicCloudService ports.PublicCloudService,
) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: !required,
		Required: required,
		Attributes: map[string]schema.Attribute{
			"billing_frequency": schema.Int64Attribute{
				Computed:    !required,
				Required:    required,
				Description: "The billing frequency (in months). Valid options are " + publicCloudService.GetBillingFrequencies().Markdown(),
				Validators: []validator.Int64{
					int64validator.OneOf(publicCloudService.GetBillingFrequencies().ToInt64()...),
				},
			},
			"term": schema.Int64Attribute{
				Computed:    !required,
				Required:    required,
				Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + publicCloudService.GetContractTerms().Markdown(),
				Validators: []validator.Int64{
					int64validator.OneOf(publicCloudService.GetContractTerms().ToInt64()...),
				},
			},
			"type": schema.StringAttribute{
				Computed:    !required,
				Required:    required,
				Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
				Validators: []validator.String{
					stringvalidator.OneOf(publicCloudService.GetContractTypes()...),
				},
			},
			"ends_at": schema.StringAttribute{Computed: true},
			"state": schema.StringAttribute{
				Computed: true,
			},
		},
		Validators: []validator.Object{customValidator.ContractTermIsValid()},
	}
}
