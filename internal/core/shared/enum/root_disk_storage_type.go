package enum

type RootDiskStorageType string

func (r RootDiskStorageType) String() string {
	return string(r)
}

func (r RootDiskStorageType) Values() []string {
	var stringValues []string

	for _, rootDiskStorageType := range rootDiskStorageTypes {
		stringValues = append(stringValues, string(rootDiskStorageType))
	}

	return stringValues
}

const (
	RootDiskStorageTypeCentral RootDiskStorageType = "CENTRAL"
	RootDiskStorageTypeLocal   RootDiskStorageType = "LOCAL"
)

var rootDiskStorageTypes = []RootDiskStorageType{
	RootDiskStorageTypeCentral,
	RootDiskStorageTypeLocal,
}

func NewRootDiskStorageType(value string) (RootDiskStorageType, error) {
	return findEnumForString(value, rootDiskStorageTypes, RootDiskStorageTypeLocal)
}
