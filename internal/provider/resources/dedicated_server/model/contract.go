package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type Contract struct {
	Id             types.String `tfsdk:"id"`
	CustomerId     types.String `tfsdk:"customer_id"`
	DeliveryStatus types.String `tfsdk:"delivery_status"`
	Reference      types.String `tfsdk:"reference"`
	SalesOrgId     types.String `tfsdk:"sales_org_id"`
}
