package dedicatedserver

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DedicatedServerModel is implemented incorrectly. Resource & attribute schema should closely match the underlying API https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
type DedicatedServerModel struct {
	ID                           types.String `tfsdk:"id"`
	Reference                    types.String `tfsdk:"reference"`
	ReverseLookup                types.String `tfsdk:"reverse_lookup"`
	DHCPLease                    types.String `tfsdk:"dhcp_lease"`
	PoweredOn                    types.Bool   `tfsdk:"powered_on"`
	PublicNetworkInterfaceOpened types.Bool   `tfsdk:"public_network_interface_opened"`
	PublicIPNullRouted           types.Bool   `tfsdk:"public_ip_null_routed"`
	PublicIP                     types.String `tfsdk:"public_ip"`
	RemoteManagementIP           types.String `tfsdk:"remote_management_ip"`
	InternalMAC                  types.String `tfsdk:"internal_mac"`
	Location                     types.Object `tfsdk:"locationModel"`
}

type locationModel struct {
	Rack  types.String `tfsdk:"rack"`
	Site  types.String `tfsdk:"site"`
	Suite types.String `tfsdk:"suite"`
	Unit  types.String `tfsdk:"unit"`
}

func (l locationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"rack":  types.StringType,
		"site":  types.StringType,
		"suite": types.StringType,
		"unit":  types.StringType,
	}
}
