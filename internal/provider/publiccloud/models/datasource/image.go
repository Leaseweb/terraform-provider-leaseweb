package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Image struct {
	Id types.String `tfsdk:"id"`
}
