package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type Hdd struct {
	Id              types.String  `tfsdk:"id"`
	Amount          types.Int32   `tfsdk:"amount"`
	Size            types.Float32 `tfsdk:"size"`
	Type            types.String  `tfsdk:"type"`
	Unit            types.String  `tfsdk:"unit"`
	PerformanceType types.String  `tfsdk:"performance_type"`
}
