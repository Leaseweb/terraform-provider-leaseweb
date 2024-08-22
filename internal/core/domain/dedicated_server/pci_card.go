package dedicated_server

type PciCard struct {
	Description string
}

func (p PciCard) String() string {
	return p.Description
}

func NewPciCard(description string) PciCard {
	return PciCard{
		Description: description,
	}
}
