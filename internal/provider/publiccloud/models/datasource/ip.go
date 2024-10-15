package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Ip struct {
	Ip types.String `tfsdk:"ip"`
}
