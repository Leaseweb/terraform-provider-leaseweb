package domain

type Image struct {
	Id           string
	Name         string
	Version      string
	Family       string
	Flavour      string
	MarketApps   []string
	StorageTypes []string
}

func NewImage(
	id string,
	name string,
	version string,
	family string,
	flavour string,
	marketApps []string,
	storageTypes []string,
) Image {
	return Image{
		Id:           id,
		Name:         name,
		Version:      version,
		Family:       family,
		Flavour:      flavour,
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
	}
}
