package dedicated_server

type Specs struct {
	Chassis             string
	HardwareRaidCapable bool
	Cpu                 Cpu
	Ram                 Ram
	Hdds                Hdds
	PciCards            PciCards
}

func NewSpecs(chassis string, hardwareRaidCapable bool, cpu Cpu, ram Ram, hdds Hdds, pciCards PciCards) Specs {
	return Specs{
		Chassis:             chassis,
		HardwareRaidCapable: hardwareRaidCapable,
		Cpu:                 cpu,
		Ram:                 ram,
		Hdds:                hdds,
		PciCards:            pciCards,
	}
}
