package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type Specs struct {
	Chassis             types.String `tfsdk:"chassis"`
	HardwareRaidCapable types.Bool   `tfsdk:"hardware_raid_capable"`
	Cpu                 Cpu          `tfsdk:"cpu"`
	Ram                 Ram          `tfsdk:"ram"`
	Hdds                Hdds         `tfsdk:"hdds"`
	PciCards            PciCards     `tfsdk:"pci_cards"`
}
