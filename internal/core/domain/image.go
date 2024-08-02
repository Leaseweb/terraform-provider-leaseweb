package domain

type Image struct {
	Id           string
	Name         string
	Version      string
	Family       string
	Flavour      string
	Architecture string
	MarketApps   []string
	StorageTypes []string
}

func NewImage(
	id string,
	name string,
	version string,
	family string,
	flavour string,
	architecture string,
	marketApps []string,
	storageTypes []string,
) Image {
	return Image{
		Id:           id,
		Name:         name,
		Version:      version,
		Family:       family,
		Flavour:      flavour,
		Architecture: architecture,
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
	}
}
