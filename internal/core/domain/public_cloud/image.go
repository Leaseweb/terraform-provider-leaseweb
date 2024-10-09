package public_cloud

type Image struct {
	Id      string
	Name    string
	Family  string
	Flavour string
	Custom  bool
}

func NewImage(
	id string,
	name string,
	family string,
	flavour string,
	custom bool,
) Image {
	return Image{
		Id:      id,
		Name:    name,
		Family:  family,
		Flavour: flavour,
		Custom:  custom,
	}
}
