package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InstanceType struct {
	Name         types.String `tfsdk:"name"`
	Resources    Resources    `tfsdk:"resources"`
	Prices       Prices       `tfsdk:"prices"`
	StorageTypes []string     `tfsdk:"storage_types"`
}
