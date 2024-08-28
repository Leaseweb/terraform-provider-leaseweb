package dedicated_server

type ControlPanel struct {
	Id   string
	Name string
}

// NewControlPanel creates ControlPanel.
func NewControlPanel(
	id string,
	name string,
) ControlPanel {
	return ControlPanel{
		Id:   id,
		Name: name,
	}
}
